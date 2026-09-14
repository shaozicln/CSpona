package middleware

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Cors 跨域中间件（Cookie 模式不能用 AllowOrigins: *）
func Cors() gin.HandlerFunc {
	allowList := map[string]bool{
		"http://localhost:5173":  true,
		"http://127.0.0.1:5173":  true,
		"http://localhost:4173":  true,
		"http://127.0.0.1:4173":  true,
		"https://cspona.top":     true,
		"https://www.cspona.top": true,
		"http://cspona.top":      true,
	}

	return cors.New(
		cors.Config{
			AllowOriginFunc: func(origin string) bool {
				if origin == "" {
					return true
				}
				if allowList[origin] {
					return true
				}
				if strings.HasPrefix(origin, "http://localhost:") || strings.HasPrefix(origin, "http://127.0.0.1:") {
					return true
				}
				return false
			},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Authorization", "Content-Type", "Accept", "Origin", "X-Requested-With", "X-Visitor-Id"},
			ExposeHeaders:    []string{"Content-Length", "Content-Type"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		},
	)
}
