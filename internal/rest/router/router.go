package router

import (
	"github.com/gin-gonic/gin"

	"social-network-otus/internal/auth"
	"social-network-otus/internal/rest/handler"
	"social-network-otus/internal/rest/middleware"
	"social-network-otus/internal/rest/response"
	"social-network-otus/internal/user"
)

type Router struct {
	handler     *handler.RestHandler
	engine      *gin.Engine
	authService *auth.AuthService
	userService *user.Service
	response    *response.ResponseFactory
}

func NewRouter(engine *gin.Engine, handler *handler.RestHandler, authService *auth.AuthService, userService *user.Service, response *response.ResponseFactory) *Router {
	router := Router{handler: handler, engine: engine, authService: authService, userService: userService, response: response}
	router.initRoutes()

	return &router
}

func (r *Router) initRoutes() {
	unauthorised := r.engine.Group("/")
	authorized := r.engine.Group("/", middleware.AuthRequired(r.authService, r.userService, r.response))
	r.initAuthRoutes(unauthorised)
	r.initUserRoutes(authorized)
	r.initFriendRoutes(authorized)
}
