package handler

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewRestHandler,
		fx.Annotate(NewAuthHandler, fx.As(new(AuthHandler))),
		fx.Annotate(NewUserHandler, fx.As(new(UserHandler))),
		fx.Annotate(NewFriendHandler, fx.As(new(FriendHandler))),
		fx.Annotate(NewPostHandler, fx.As(new(PostHandler))),
		fx.Annotate(NewFeedHandler, fx.As(new(FeedHandler))),
	),
)
