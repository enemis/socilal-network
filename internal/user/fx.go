package user

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewUserService,
		fx.Annotate(
			NewUserRepository,
			fx.As(new(UserRepository)),
		),
	),
)
