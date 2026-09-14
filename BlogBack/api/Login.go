package api

import (
	"BlogBack/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Login(c *gin.Context) {
	var cinuser User
	if err := c.ShouldBindJSON(&cinuser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "无效的请求数据", "data": nil})
		return
	}
	var user User
	result := db.Where("username = ?", cinuser.Username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "喵喵喵？注册了吗就来登录？", "message": "喵喵喵？注册了吗就来登录？"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "查询用户失败"})
		}
		return
	}
	if user.Password != ScryptPw(cinuser.Password) {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "密码错误", "message": "密码错误"})
		return
	}

	token, err := utils.SignSessionToken(user.Id, user.Username, user.RoleQx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "创建会话失败"})
		return
	}
	setSessionCookie(c, token)

	articleCount := int64(0)
	db.Model(&Article{}).Where("user_id = ?", user.Id).Count(&articleCount)
	user.ArticleCount = int(articleCount)

	// 不再把 password 回给浏览器
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"msg":     "登录成功",
		"message": "登录成功",
		"id":      user.Id,
		"qx":      user.RoleQx,
		"username": user.Username,
		"email":   user.Email,
		"avatar":  user.Avatar,
		"avatarUrl": utils.ResolveImageURL(user.Avatar, "boli.jpg"),
		"articleCount": user.ArticleCount,
	})
}
