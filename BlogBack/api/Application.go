package api

import (
	"BlogBack/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Application struct {
	Id           uint   `gorm:"primary key; autoIncrement" column:"id"`
	Username     string `gorm:"type:varchar(255)" column:"username"`
	Email        string `gorm:"type:varchar(255)" column:"email"`
	Name         string `gorm:"type:varchar(255)" column:"name"`
	Web          string `gorm:"type:varchar(255)" column:"web"`
	Introduction string `gorm:"type:varchar(255)" column:"introduction"`
	Img          string `gorm:"type:varchar(255)" column:"img"`
	Avatar       string `gorm:"type:varchar(255)" column:"avatar"`
	Background   string `gorm:"type:varchar(255)" column:"background"`
	Description  string `gorm:"type:varchar(255)" column:"description"`
}

func (Application) TableName() string {
	return "applications"
}

func GetApplication(c *gin.Context) {
	id := c.Query("id")
	username := c.Query("username")
	web := c.Query("web")
	var applications []Application
	if username != "" {
		db.Where("name LIKE ?", "%"+username+"%").Find(&applications)
	} else if web != "" {
		db.Where("web = ?", web).Find(&applications)
	} else if id != "" {
		db.Where("Id = ?", id).Find(&applications)
	} else {
		db.Find(&applications)
	}
	c.JSON(http.StatusOK, gin.H{"data": applications})
}

func PostApplication(c *gin.Context) {
	var application Application

	// 获取文件
	img, _ := c.FormFile("img")
	background, _ := c.FormFile("background")

	baseDir := utils.GetImageBaseDir()
	_ = os.MkdirAll(baseDir, os.ModePerm)

	// 处理图片文件
	imgname := filepath.Base(img.Filename)
	savePath := filepath.Join(baseDir, imgname)
	savePath = strings.ReplaceAll(savePath, "\\", "/")
	_ = c.SaveUploadedFile(img, savePath)

	// 处理背景文件
	backgroundname := filepath.Base(background.Filename)
	savePath2 := filepath.Join(baseDir, backgroundname)
	savePath2 = strings.ReplaceAll(savePath2, "\\", "/")
	_ = c.SaveUploadedFile(background, savePath2)

	// 填充数据
	application.Name = c.PostForm("name")
	application.Username = c.PostForm("username")
	application.Email = c.PostForm("email")
	application.Web = c.PostForm("web")
	application.Introduction = c.PostForm("introduction")
	application.Img = imgname
	application.Avatar = c.PostForm("avatar")
	application.Background = backgroundname
	application.Description = c.PostForm("description")

	db.Create(&application)
	c.JSON(200, gin.H{"message": "created successfully", "data": application})
}

func DeleteApplication(c *gin.Context) {
	var applications []Application
	id := c.Param("id")
	id1, _ := strconv.ParseUint(id, 10, 64)
	// 检查是否存在该ID的记录
	db.Where("Id=?", id1).First(&applications)
	// 如果记录存在，则删除记录
	err = db.Where("Id=?", id1).Delete(&applications).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "删除成功"})
}
