package friend

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewFriendService,
		fx.Annotate(
			NewFriendRepository,
			fx.As(new(FriendRepository)),
		),
	),
)
