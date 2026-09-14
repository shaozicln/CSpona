package api

import (
	"BlogBack/middleware"
	"BlogBack/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setSessionCookie(c *gin.Context, token string) {
	secure := utils.GetEnv() == utils.ProdEnv
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		utils.SessionCookieName,
		token,
		int(utils.SessionMaxAge.Seconds()),
		"/",
		"",
		secure,
		true, // HttpOnly：JS 读不到
	)
}

func clearSessionCookie(c *gin.Context) {
	secure := utils.GetEnv() == utils.ProdEnv
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(utils.SessionCookieName, "", -1, "/", "", secure, true)
}

// Logout 清除会话 Cookie
func Logout(c *gin.Context) {
	clearSessionCookie(c)
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "已退出登录", "data": nil})
}

// AuthMe 当前登录用户（靠 Cookie）；未登录返回 200 + data:null
func AuthMe(c *gin.Context) {
	uid, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "未登录", "data": nil})
		return
	}
	var user User
	if err := db.Select("id", "username", "email", "avatar", "role_qx", "article_count", "created_at", "description").
		First(&user, uid).Error; err != nil {
		clearSessionCookie(c)
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "用户不存在或会话失效", "data": nil})
		return
	}
	var articleCount int64
	db.Model(&Article{}).Where("user_id = ?", user.Id).Count(&articleCount)
	user.ArticleCount = int(articleCount)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": gin.H{
			"id":           user.Id,
			"username":     user.Username,
			"email":        user.Email,
			"avatar":       user.Avatar,
			"avatarUrl":    utils.ResolveImageURL(user.Avatar, "boli.jpg"),
			"qx":           user.RoleQx,
			"articleCount": user.ArticleCount,
			"description":  user.Description,
		},
	})
}
