package router

import (
	"github.com/gin-gonic/gin"
)

func (r *Router) initAuthRoutes(authorizedGroup *gin.RouterGroup) {
	auth := authorizedGroup.Group("/auth")
	{
		auth.POST("/sign-in", r.handler.Auth.SignIn)
		auth.POST("/sign-up", r.handler.Auth.SignUp)
	}

}
