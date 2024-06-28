package handler

import "go.uber.org/fx"

type RestHandlerParams struct {
	fx.In

	Auth   AuthHandler
	User   UserHandler
	Friend FriendHandler
	Post   PostHandler
	Feed   FeedHandler
	WS     WSHandler
}

type RestHandler struct {
	Feed   FeedHandler
	Auth   AuthHandler
	User   UserHandler
	Friend FriendHandler
	Post   PostHandler
	WS     WSHandler
}

func NewRestHandler(params RestHandlerParams) *RestHandler {
	return &RestHandler{
		Auth:   params.Auth,
		User:   params.User,
		Friend: params.Friend,
		Post:   params.Post,
		Feed:   params.Feed,
		WS:     params.WS,
	}
}
