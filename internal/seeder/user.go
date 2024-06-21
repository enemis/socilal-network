package seeder

import (
	"fmt"
	"github.com/essentialkaos/translit/v2"
	"github.com/go-faker/faker/v4"
	"math/rand"
	"reflect"
	"social-network-otus/internal/user"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type threadSafeSlice[T comparable] struct {
	Slice []T
	mutex *sync.Mutex
}

func (s *threadSafeSlice[T]) Append(val T) {
	s.mutex.Lock()
	s.Slice = append(s.Slice, val)
	s.mutex.Unlock()
}

func (s *Seeder) UserSeed(count uint) {
	wg := &sync.WaitGroup{}
	userIds := threadSafeSlice[string]{Slice: make([]string, 0, count), mutex: &sync.Mutex{}}
	jobs := make(chan int)

	for i := 0; i < 2; i++ {
		// more than 2 gouroutins bring to tcp connection error, looks like it's docker issues.
		// https://github.com/lib/pq/issues/835
		go userWorker(s, i, jobs, wg, &userIds)
		wg.Add(1)
	}

	for j := 0; j < int(count); j++ {
		jobs <- j
	}

	close(jobs)
	wg.Wait()
	s.makeFriendshipWorker(&userIds)
}

func userWorker(s *Seeder, workerId int, jobs <-chan int, wg *sync.WaitGroup, userIds *threadSafeSlice[string]) {
	for j := range jobs {
		userModel := user.User{}

		person := faker.GetPerson()
		v := reflect.ValueOf(person)
		if j%2 == 0 {
			result, err := person.RussianFirstNameMale(v)
			if err != nil {
				logrus.Panic(err)
			}
			userModel.Name = result.(string)
			result, err = person.RussianLastNameMale(v)
			if err != nil {
				logrus.Panic(err)
			}
			userModel.Surname = result.(string)
		} else {
			result, err := person.RussianFirstNameFemale(v)
			if err != nil {
				logrus.Panic(err)
			}
			userModel.Name = result.(string)
			result, err = person.RussianLastNameFemale(v)
			if err != nil {
				logrus.Panic(err)
			}
			userModel.Surname = result.(string)
		}

		userModel.Birthday, _ = time.Parse(time.DateOnly, faker.Date())
		userModel.Biography = faker.Paragraph()
		address := faker.GetRealAddress()
		userModel.Email = fmt.Sprintf("%s.%s%d@gmail.com", translit.ISO9B(userModel.Surname), translit.ISO9B(userModel.Name), j)
		userModel.City = address.City
		userModel.Password = "12345"
		userId, err := s.userService.CreateUser(&userModel)
		if err != nil {
			fmt.Println(err.OriginalError())
		}

		if userId != nil {
			s.PostSeed(*userId, uint(rand.Int31n(100)+10))
			userIds.Append(userId.String())
		}
		fmt.Printf("worker %d done iteration #%d\n%s %s (%s)\n", workerId, j, userModel.Surname, userModel.Name, userModel.Email)
	}
	wg.Done()
}

func (s *Seeder) makeFriendshipWorker(userIds *threadSafeSlice[string]) {
	jobs := make(chan string)
	wg := &sync.WaitGroup{}
	for i := 0; i < 2; i++ {
		go s.makeFriendship(jobs, wg, userIds)
		wg.Add(1)
	}
	for _, userId := range userIds.Slice {
		jobs <- userId
	}
	close(jobs)
	wg.Wait()
}

func (s *Seeder) makeFriendship(jobs <-chan string, wg *sync.WaitGroup, userIds *threadSafeSlice[string]) {
	for userId := range jobs {
		friendCount := min(int32(rand.Int31n(int32(len(userIds.Slice)-1))+10), int32(len(userIds.Slice)-1))
		friendCount = min(friendCount, 1000)
		friendCountFreeze := friendCount
		friendsIds := make(map[string]string)
		for friendCount > 0 {
			friendId := userIds.Slice[rand.Int31n(int32(len(userIds.Slice)))]
			if friendsIds[friendId] != "" {
				continue
			}

			friendsIds[friendId] = friendId
			friendCount--
		}
		query := "Insert Into friends (user_id, friend_id, created_at) VALUES "
		elementIndex := 0
		for _, friendId := range friendsIds {
			query += fmt.Sprintf("('%s', '%s', '%s')", userId, friendId, faker.Date())
			if elementIndex != int(friendCountFreeze-1) {
				query += fmt.Sprintf(",")
			}
			elementIndex++
		}
		_, err := s.db.Master().Exec(query)
		if err != nil {
			logrus.Debugln(err)
		}
	}

	wg.Done()
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
