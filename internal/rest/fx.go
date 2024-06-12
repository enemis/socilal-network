package rest

import (
	"go.uber.org/fx"

	"social-network-otus/internal/rest/handler"
	"social-network-otus/internal/rest/response"
	"social-network-otus/internal/rest/router"
)

var Module = fx.Options(
	handler.Module,
	fx.Provide(
		response.NewResponseFactory,
		NewRestServer,
		router.NewRouter,
	),
	fx.Invoke(
		InitHooks,
	),
)
