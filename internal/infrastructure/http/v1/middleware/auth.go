package middleware

import (
	"context"
	"prueba-go/pkg/types"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		actor := c.GetHeader("X-User-Id")
		if actor != "" {
			c.Set(string(types.ActorKey), actor)

			ctx := context.WithValue(c.Request.Context(), types.ActorKey, actor)
			c.Request = c.Request.WithContext(ctx)
		}
		c.Next()
	}
}
