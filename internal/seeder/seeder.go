package seeder

import (
	"social-network-otus/internal/auth"
)

type Seeder struct {
	authService *auth.AuthService
}

func NewSeeder(authService *auth.AuthService) *Seeder {
	return &Seeder{authService: authService}
}

func (s *Seeder) Seed() {
	s.UserSeed(10000)
}
