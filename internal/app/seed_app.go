package app

import (
	"go.uber.org/fx"
	"social-network-otus/internal/session"

	"social-network-otus/internal/auth"
	"social-network-otus/internal/config"
	"social-network-otus/internal/database"
	"social-network-otus/internal/friend"
	"social-network-otus/internal/logger"
	"social-network-otus/internal/post"
	"social-network-otus/internal/seeder"
	"social-network-otus/internal/token"
	"social-network-otus/internal/user"
)

type SeedApp struct {
	config    *config.Config
	seeder    *seeder.Seeder
	container *fx.App
}

func NewSeeder() (*SeedApp, error) {
	fxContainer := fx.New(
		config.Module,
		logger.Module,
		database.Module,
		auth.Module,
		seeder.Module,
		token.Module,
		user.Module,
		friend.Module,
		post.Module,
		session.Module,
	)

	app := SeedApp{container: fxContainer}

	return &app, nil
}
