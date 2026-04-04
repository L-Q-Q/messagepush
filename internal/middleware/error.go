package middleware

import (
	"log"
	"net/http"

	"message-push-system/pkg/response"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				response.Error(c, http.StatusInternalServerError, 500, "internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}
