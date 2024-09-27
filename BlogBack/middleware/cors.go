package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors(config cors.Config) gin.HandlerFunc {
	return cors.New(config)
}
