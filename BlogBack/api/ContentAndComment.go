package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Comment struct {
	Id        uint      `gorm:"primary key; autoIncrement" column:"id"`
	ArticleId int       `gorm:"column:article_id"`
	UserId    int       `gorm:"column:user_id"`
	Idea      string    `gorm:"column:idea"`
	CreatedAt time.Time `gorm:"type:timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp"`
	User      User
}

func (Comment) TableName() string {
	return "comments"
}

func GetContent(c *gin.Context) {
	id := c.Param("id")
	var article Article
	if err := db.Preload("User").Preload("Comments").Preload("Comments.User").Find(&article, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
		return
	}
	db.Where("id = ?", id).First(&article)
	article.ViewCount++
	db.Save(&article)
	var articleCount int64
	db.Model(&Article{}).Where("user_id = ?", article.UserId).Count(&articleCount)
	article.User.ArticleCount = int(articleCount)
	var commentCount int64
	db.Model(&Comment{}).Where("article_id = ?", id).Count(&commentCount)
	article.CommentCount = uint(commentCount)
	var category Category
	db.Model(&Category{}).Where("id = ?", article.CategoryId).Find(&category)
	article.Category.Name = category.Name
	c.JSON(http.StatusOK, gin.H{"data": article})
}

func PostContent(c *gin.Context) {
	id1 := c.Param("id")
	id2, err := strconv.ParseUint(id1, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "无效的文章ID"})
		return
	}
	var comment Comment
	if err := c.ShouldBindJSON(&comment); err == nil {
		comment.ArticleId = int(id2)
		db.Create(&comment)
		db.Where("article_id = ?", id2).Preload("User").Find(&comment)
		c.JSON(200, gin.H{"data": comment})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	fmt.Println(comment.ArticleId)
}

func DeleteContent(c *gin.Context) {
	var comments []Comment
	id1 := c.Param("id")
	id, err := strconv.ParseUint(id1, 10, 64) // 解析ID，并处理可能的错误
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的ID格式"})
		return
	}

	// 检查是否存在该ID的记录
	db.Where("id=?", id).First(&comments)
	// 如果记录存在，则删除记录
	err = db.Where("id=?", id).Delete(&comments).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "删除成功"})
}

func GetComment(c *gin.Context) {
	articleIdStr := c.Param("articleId")
	articleId, err := strconv.Atoi(articleIdStr)
	if err != nil || articleId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}
	idea := c.Query("idea")
	username := c.Query("username")
	var comments []Comment
	if idea != "" {
		db.Where("article_id = ?  AND  idea LIKE ?", articleId, "%"+idea+"%").Preload("User").Find(&comments)
	} else if username != "" {
		var user User
		db.Where("username = ?", username).Find(&user)
		db.Where("article_id = ? AND user_id = ?", articleId, user.Id).Preload("User").Find(&comments)
	} else {
		var user User
		db.Where("username = ?", username).Find(&user)
		db.Where("article_id = ?", articleId).Preload("User").Find(&comments)
	}

	c.JSON(http.StatusOK, gin.H{"data": comments})
}
