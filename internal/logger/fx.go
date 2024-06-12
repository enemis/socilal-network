package logger

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			NewLogger,
			fx.As(new(LoggerInterface)),
		),
	),
)
