package middleware

import (
	"fmt"
	"gin-use-demo/response"
	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				response.Fail(c, fmt.Sprint(err), nil)
				return
			}
		}()
		c.Next()
	}
}
