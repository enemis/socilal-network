package router

import "github.com/gin-gonic/gin"

func (r *Router) initWSRoutes(authorizedGroup *gin.RouterGroup) {
	authorizedGroup.Any("/ws", r.handler.WS.Connect)
}
