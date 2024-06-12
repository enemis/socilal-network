package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"social-network-otus/internal/auth"
	"social-network-otus/internal/rest/response"
	"social-network-otus/internal/user"
)

const (
	AuthorizationHeadder = "Authorization"
	UserContext          = "User"
)

func AuthRequired(authService *auth.AuthService, userService *user.Service, response *response.ResponseFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(AuthorizationHeadder)

		if header == "" {
			response.Unauthorised(c)
			return
		}

		headerParts := strings.Split(header, " ")

		if len(headerParts) != 2 {
			response.Unauthorised(c)
			return
		}

		userId, err := authService.ParseToken(headerParts[1])

		if err != nil {
			response.Unauthorised(c)
			return
		}

		userModel, err := userService.GetUserById(userId.String())

		if err != nil {
			response.Unauthorised(c)
			return
		}
		c.Set(UserContext, userModel)
		c.Next()
	}
}
