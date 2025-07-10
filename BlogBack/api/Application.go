package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"BlogBack/utils"
)

type Application struct {
	Id           uint   `gorm:"primary key; autoIncrement" column:"id"`
	Username     string `gorm:"type:varchar(255)" column:"username"`
	Email        string `gorm:"type:varchar(255)" column:"email"`
	Name         string `gorm:"type:varchar(255)" column:"name"`
	Web          string `gorm:"type:varchar(255)" column:"web"`
	Introduction string `gorm:"type:varchar(255)" column:"introduction"`
	Img          string `gorm:"type:varchar(255)" column:"img"`
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
	//获取文件，错误处理
	file, err := c.FormFile("img")
	if err != nil {
		c.JSON(500, gin.H{"error": "无法获取文件"})
		return
	}
	baseDir := utils.GetImageBaseDir()          //定义基本目录
	filename := filepath.Base(file.Filename)           //获取上传的文件名
	savePath := filepath.Join(baseDir, filename)       //连接字段，形成存储路径
	savePath = strings.ReplaceAll(savePath, "\\", "/") // 将路径用正斜杠保存，兼容不同操作系统，同时方便前端读取

	_ = os.MkdirAll(baseDir, os.ModePerm) //创建保存路径所在目录，目录存在则忽略，同时忽略了错误处理
	_ = c.SaveUploadedFile(file, savePath)

	application.Img = filename            //只存储文件名字
	application.Name = c.PostForm("name") //数据库其他字段和数据
	application.Username = c.PostForm("username")
	application.Email = c.PostForm("email")
	application.Web = c.PostForm("web")
	application.Introduction = c.PostForm("introduction")

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
