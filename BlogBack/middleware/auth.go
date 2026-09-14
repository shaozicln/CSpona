package middleware

import (
	"BlogBack/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CtxUserID   = "authUserId"
	CtxUsername = "authUsername"
	CtxRoleQx   = "authRoleQx"
)

// AuthRequired 要求已登录（HttpOnly Cookie 中的 JWT）
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := readClaims(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": utils.ErrUnauthorized,
				"msg":  "请先登录",
				"data": nil,
			})
			return
		}
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxUsername, claims.Username)
		c.Set(CtxRoleQx, claims.RoleQx)
		c.Next()
	}
}

// AuthOptional 有 Cookie 则解析，没有也不拦截
func AuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		if claims, err := readClaims(c); err == nil {
			c.Set(CtxUserID, claims.UserID)
			c.Set(CtxUsername, claims.Username)
			c.Set(CtxRoleQx, claims.RoleQx)
		}
		c.Next()
	}
}

func readClaims(c *gin.Context) (*utils.SessionClaims, error) {
	tokenStr, err := c.Cookie(utils.SessionCookieName)
	if err != nil || tokenStr == "" {
		// 兼容：也接受 Authorization: Bearer
		auth := c.GetHeader("Authorization")
		if len(auth) > 7 && (auth[:7] == "Bearer " || auth[:7] == "bearer ") {
			tokenStr = auth[7:]
		} else {
			return nil, err
		}
	}
	return utils.ParseSessionToken(tokenStr)
}

func GetUserID(c *gin.Context) (uint, bool) {
	v, ok := c.Get(CtxUserID)
	if !ok {
		return 0, false
	}
	id, ok := v.(uint)
	return id, ok
}

func GetRoleQx(c *gin.Context) string {
	v, _ := c.Get(CtxRoleQx)
	s, _ := v.(string)
	return s
}
