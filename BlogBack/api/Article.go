package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Article struct {
	Id           uint      `gorm:"primary key"`
	Title        string    `gorm:"type:varchar(255)"`
	Content      string    `gorm:"type:text"`
	CategoryId   uint      `gorm:"type:int"`
	UserId       uint      `gorm:"type:int"`
	ViewCount    uint      `gorm:"type:int"`
	CommentCount uint      `gorm:"type:int"`
	CreatedAt    time.Time `gorm:"type:timestamp"`
	UpdatedAt    time.Time `gorm:"type:timestamp"`
	Img          string    `gorm:"type:varchar(255)"`
	User         User
	Comments     []Comment
	Category     Category
}

func (Article) TableName() string {
	return "articles"
}

func SearchArticle(c *gin.Context) {
	title := c.Query("title")
	id := c.Query("id")
	content := c.Query("content")
	var articles []Article
	if title != "" {
		db.Where("title LIKE ?", "%"+title+"%").Find(&articles)
	} else if id != "" {
		db.Where("id = ?", id).Find(&articles)
	} else if content != "" {
		db.Where("content LIKE ?", "%"+content+"%").Find(&articles)
	} else {
		db.Find(&articles)
	}
	c.JSON(http.StatusOK, gin.H{"data": articles})
}

func PostArticle(c *gin.Context) {
	var article Article
	//获取文件，错误处理
	file, err := c.FormFile("img")
	if err != nil {
		c.JSON(500, gin.H{"error": "无法获取文件"})
		return
	}
	baseDir := "../BlogFont/public/Pictures"           //定义基本目录
	filename := filepath.Base(file.Filename)           //获取上传的文件名
	savePath := filepath.Join(baseDir, filename)       //连接字段，形成存储路径
	savePath = strings.ReplaceAll(savePath, "\\", "/") // 将路径用正斜杠保存，兼容不同操作系统，同时方便前端读取

	_ = os.MkdirAll(baseDir, os.ModePerm) //创建保存路径所在目录，目录存在则忽略，同时忽略了错误处理
	_ = c.SaveUploadedFile(file, savePath)

	article.Img = filename              //只存储文件名字
	article.Title = c.PostForm("title") //数据库其他字段和数据
	article.Content = c.PostForm("content")
	categoryId, _ := strconv.Atoi(c.PostForm("category_id"))
	article.CategoryId = uint(categoryId)

	userId, _ := strconv.Atoi(c.PostForm("user_id"))
	article.UserId = uint(userId)

	db.Create(&article)
	c.JSON(200, gin.H{"message": "created successfully", "data": article})
}

func PutArticle(c *gin.Context) {
	id := c.Param("id")

	// 绑定请求体到Article结构体
	var article Article
	if err := c.BindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	id1, _ := strconv.ParseUint(id, 10, 64)
	article.Id = uint(id1)
	// 执行更新操作
	if result := db.Model(&Article{}).Where("id = ?", id).Updates(Article{
		Id:         article.Id,
		Title:      article.Title,
		Content:    article.Content,
		CategoryId: article.CategoryId,
		UserId:     article.UserId,
	}); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": article})
	}
}

func DeleteArticle(c *gin.Context) {
	var article Article
	db.Delete(&article, c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"data": article})
}
