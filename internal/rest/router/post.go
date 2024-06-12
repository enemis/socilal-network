package router

import "github.com/gin-gonic/gin"

func (r *Router) initPostRoutes(authorizedGroup *gin.RouterGroup) {
	user := authorizedGroup.Group("/posts")
	{
		user.GET("/:id", r.handler.Post.GetPost)
		user.POST("/", r.handler.Post.CreatePost)
		user.PUT("/", r.handler.Post.UpdatePost)
		user.DELETE("/:id", r.handler.Post.DeletePost)
	}
}
