package router

import "github.com/gin-gonic/gin"

func (r *Router) initFriendRoutes(authorizedGroup *gin.RouterGroup) {
	user := authorizedGroup.Group("/friends")
	{
		user.PUT("/", r.handler.Friend.AddFriend)
		user.DELETE("/", r.handler.Friend.RemoveFriend)
	}
}
