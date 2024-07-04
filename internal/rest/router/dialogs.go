package router

import (
	"github.com/gin-gonic/gin"
)

func (r *Router) initDialogsRoutes(authorizedGroup *gin.RouterGroup) {
	dialogs := authorizedGroup.Group("/dialog")
	{
		dialogs.POST("/send", r.handler.Dialogs.SendMessage)
		dialogs.GET("/:user_id/list", r.handler.Dialogs.GetDialog)
	}
}
