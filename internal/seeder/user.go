package seeder

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"otus-social-network/internal/dto"

	"github.com/sirupsen/logrus"
)

func (s *Seeder) UserSeed(count uint) {
	wg := &sync.WaitGroup{}

	jobs := make(chan int)

	for i := 0; i < 2; i++ {
		// more than 2 gouroutins bring to tcp connection error, looks like it's docker issues.
		// https://github.com/lib/pq/issues/835
		go userWorker(s, i, jobs, wg)
		wg.Add(1)
	}

	for j := 0; j < int(count); j++ {
		jobs <- j
	}

	close(jobs)
	wg.Wait()
}

func userWorker(s *Seeder, workerId int, jobs <-chan int, wg *sync.WaitGroup) {
	for j := range jobs {
		user := dto.SignUpInput{}

		person := faker.GetPerson()
		v := reflect.ValueOf(person)
		if j%2 == 0 {
			result, err := person.RussianFirstNameMale(v)
			if err != nil {
				logrus.Panic(err)
			}
			user.Name = result.(string)
			result, err = person.RussianLastNameMale(v)
			if err != nil {
				logrus.Panic(err)
			}
			user.Surname = result.(string)
		} else {
			result, err := person.RussianFirstNameFemale(v)
			if err != nil {
				logrus.Panic(err)
			}
			user.Name = result.(string)
			result, err = person.RussianLastNameFemale(v)
			if err != nil {
				logrus.Panic(err)
			}
			user.Surname = result.(string)
		}

		user.Birthday, _ = time.Parse(time.DateOnly, faker.Date())
		user.Biography = faker.Paragraph()
		address := faker.GetRealAddress()
		user.Email = fmt.Sprintf("%s.%s%d@gmail.com", translit.ISO9B(user.Surname), translit.ISO9B(user.Name), j)
		user.City = address.City
		user.Password = "12345"
		_, err := s.authService.CreateUser(&user)
		if err != nil {
			fmt.Println(err.OriginalError())
		}
		fmt.Printf("worker %d done iteration #%d\n%s %s (%s)\n", workerId, j, user.Surname, user.Name, user.Email)
	}
	wg.Done()
}
