package seeder

import (
	"social-network-otus/internal/auth"
	"social-network-otus/internal/database"
	"social-network-otus/internal/post"
	"social-network-otus/internal/user"
)

type Seeder struct {
	authService *auth.AuthService
	userService *user.Service
	postService *post.PostService
	db          *database.DatabaseStack
}

func NewSeeder(authService *auth.AuthService, userService *user.Service, postService *post.PostService, database *database.DatabaseStack) *Seeder {
	return &Seeder{authService: authService, userService: userService, postService: postService, db: database}
}

func RunImport(s *Seeder) {
	s.Seed()
}

func (s *Seeder) Seed() {
	s.UserSeed(10000)
}
