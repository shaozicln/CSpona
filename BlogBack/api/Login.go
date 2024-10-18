package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func Login(c *gin.Context) {
	var cinuser User
	if err := c.ShouldBindJSON(&cinuser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}
	var user User
	result := db.Where("username = ?", cinuser.Username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 用户不存在
			c.JSON(200, gin.H{"message": "喵喵喵？注册了吗就来登录？"})
		} else {
			// 查询出错
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
		}
		return
	}
	if user.Password == ScryptPw(cinuser.Password) {
		// 登录成功
		articleCount := int64(0)
		db.Model(&Article{}).Where("user_id = ?", user.Id).Count(&articleCount)
		user.ArticleCount = int(articleCount)
		c.JSON(http.StatusOK, gin.H{"message": "登录成功", "id": user.Id, "qx": user.RoleQx, "username": user.Username, "email": user.Email, "password": user.Password, "avatar": user.Avatar, "articleCount": user.ArticleCount})
	} else {
		// 密码错误
		c.JSON(200, gin.H{"message": "密码错误"})
	}
}
