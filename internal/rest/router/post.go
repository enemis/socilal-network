package router

import "github.com/gin-gonic/gin"

func (r *Router) initPostRoutes(authorizedGroup *gin.RouterGroup) {
	post := authorizedGroup.Group("/posts")
	{
		post.GET("/:id", r.handler.Post.GetPost)
		post.POST("", r.handler.Post.CreatePost)
		post.PUT("", r.handler.Post.UpdatePost)
		post.DELETE("/:id", r.handler.Post.DeletePost)
	}
}
