package router

import (
	"github.com/gin-gonic/gin"

	"social-network-otus/internal/rest/middleware"
)

func (r *Router) initAuthRoutes(authorizedGroup *gin.RouterGroup) {
	authorizedGroup.Use(middleware.AuthRequired(r.authService, r.userService, r.response))
	{
		auth := authorizedGroup.Group("/auth")
		{
			auth.POST("/sign-in", r.handler.Auth.SignIn)
			auth.POST("/sign-up", r.handler.Auth.SignUp)
		}
	}
}
