package router

import (
	"github.com/gin-gonic/gin"
)

func (r *Router) initAuthRoutes(unauthorizedGroup *gin.RouterGroup) {
	auth := unauthorizedGroup.Group("/auth")
	{
		auth.POST("/sign-in", r.handler.Auth.SignIn)
		auth.POST("/sign-up", r.handler.Auth.SignUp)
	}

}
