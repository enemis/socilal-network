package seeder

import (
	"github.com/google/uuid"
)

func (s *Seeder) PostSeed(user uuid.UUID, postCount uint) {
	faker.Paragraph()
}
