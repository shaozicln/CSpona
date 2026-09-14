package api

import (
	"BlogBack/utils"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type aboutMeBody struct {
	Content string `json:"content"`
}

// GetAboutMe 公开读取首页「关于我」Markdown
func GetAboutMe(c *gin.Context) {
	path := utils.GetAboutMePath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusOK, gin.H{
				"status":  200,
				"data":    gin.H{"content": ""},
				"message": "about-me 文件不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "读取失败: " + err.Error(),
		})
		return
	}
	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   gin.H{"content": string(data)},
	})
}

// PutAboutMe 登录后保存「关于我」Markdown（管理页）
func PutAboutMe(c *gin.Context) {
	var body aboutMeBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "参数错误",
		})
		return
	}

	path := utils.GetAboutMePath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "创建目录失败: " + err.Error(),
		})
		return
	}
	if err := os.WriteFile(path, []byte(body.Content), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "保存失败: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "保存成功",
		"data":    gin.H{"content": body.Content},
	})
}
