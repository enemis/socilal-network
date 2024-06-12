package token

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			NewPasswordGenerator,
			fx.As(new(PasswordGenerator)),
		),
	),
)
