package router

import "github.com/gin-gonic/gin"

func (r *Router) initUserRoutes(authorizedGroup *gin.RouterGroup) {
	user := authorizedGroup.Group("/profile")
	{
		user.GET("/search", r.handler.User.FindUsers)
		user.POST("/update-password", r.handler.User.UpdatePassword)
		user.GET(":id", r.handler.User.UserPage)
	}
}
