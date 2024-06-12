package router

import "github.com/gin-gonic/gin"

func (r *Router) initUserRoutes(authorizedGroup *gin.RouterGroup) {
	user := authorizedGroup.Group("/users")
	{
		user.GET("/search", r.handler.User.FindUsers)
		user.GET(":id", r.handler.User.UserPage)
	}
}
