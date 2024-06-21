package seeder

import (
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/google/uuid"
	"math/rand"
	"social-network-otus/internal/post"
	"social-network-otus/internal/utils"
	"time"
)

func (s *Seeder) PostSeed(userId uuid.UUID, postCount uint) {
	for i := uint(0); i < postCount; i++ {
		poststatus := post.Draft
		var deleted *time.Time

		if rand.Int31n(10) > 4 {
			poststatus = post.Published
		}
		if rand.Int31n(10) > 8 {
			deleted = utils.Ptr(time.Unix(faker.UnixTime(), 0))
		}

		s.postService.CreatePost(&post.Post{
			UserId:    userId,
			Title:     faker.Sentence(options.WithRandomStringLength(20)),
			Post:      faker.Paragraph(),
			Status:    poststatus,
			CreatedAt: time.Unix(faker.UnixTime(), 0),
			UpdatedAt: time.Unix(faker.UnixTime(), 0),
			DeletedAt: deleted,
		})
	}
}
