package handler

import "go.uber.org/fx"

type RestHandlerParams struct {
	fx.In

	Auth   AuthHandler
	User   UserHandler
	Friend FriendHandler
	Post   PostHandler
}

type RestHandler struct {
	Auth   AuthHandler
	User   UserHandler
	Friend FriendHandler
	Post   PostHandler
}

func NewRestHandler(params RestHandlerParams) *RestHandler {
	return &RestHandler{
		Auth:   params.Auth,
		User:   params.User,
		Friend: params.Friend,
		Post:   params.Post,
	}
}
