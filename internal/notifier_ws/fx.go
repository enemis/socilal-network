package notifier_ws

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewNotifier,
	),
)
