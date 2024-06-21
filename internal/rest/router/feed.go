package router

import "github.com/gin-gonic/gin"

func (r *Router) initFeedRoutes(authorizedGroup *gin.RouterGroup) {
	user := authorizedGroup.Group("/posts/feed")
	{
		user.GET("/", r.handler.Feed.GetFeed)
	}
}
