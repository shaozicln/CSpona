package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	if err := c.ShouldBindJSON(&article); err == nil {
		db.Create(&article)
		c.JSON(200, gin.H{"data": article})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
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
