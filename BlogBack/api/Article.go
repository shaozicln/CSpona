package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"BlogBack/utils"
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

// 图片上传接口
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("img")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法获取文件: " + err.Error()})
		return
	}
	if file.Size > 8<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图片不能超过 8MB"})
		return
	}

	baseDir := utils.GetImageBaseDir()
	filename := filepath.Base(file.Filename)
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp":
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持 jpg/png/gif/webp/bmp"})
		return
	}
	uniqueFilename := time.Now().Format("20060102150405") + "_" + filename
	savePath := filepath.Join(baseDir, uniqueFilename)

	if err := os.MkdirAll(baseDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建目录失败: " + err.Error()})
		return
	}
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存失败: " + err.Error()})
		return
	}

	imageUrl := "/Pictures/" + uniqueFilename
	c.JSON(200, gin.H{"imageUrl": imageUrl})
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
	file, err := c.FormFile("img")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法获取封面图: " + err.Error()})
		return
	}
	baseDir := utils.GetImageBaseDir()
	filename := filepath.Base(file.Filename)
	uniqueFilename := time.Now().Format("20060102150405") + "_" + filename
	savePath := filepath.Join(baseDir, uniqueFilename)

	if err := os.MkdirAll(baseDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建目录失败: " + err.Error()})
		return
	}
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "封面保存失败: " + err.Error()})
		return
	}

	article.Img = uniqueFilename
	article.Title = c.PostForm("title")
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
    
    // 先从数据库获取原始文章信息
    var article Article
    if err := db.First(&article, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
        return
    }

    // 处理文件上传（背景图）
    file, err := c.FormFile("img")
    if err == nil {
        baseDir := utils.GetImageBaseDir()
        filename := filepath.Base(file.Filename)
        uniqueFilename := time.Now().Format("20060102150405") + "_" + filename
        savePath := filepath.Join(baseDir, uniqueFilename)

        if mkErr := os.MkdirAll(baseDir, 0755); mkErr != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "创建目录失败: " + mkErr.Error()})
            return
        }
        if saveErr := c.SaveUploadedFile(file, savePath); saveErr != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "封面保存失败: " + saveErr.Error()})
            return
        }
        article.Img = uniqueFilename
    }

    // 处理其他字段，若前端传递的参数为空则使用原始数据
    title := c.PostForm("title")
    if title != "" {
        article.Title = title
    }

    content := c.PostForm("content")
    if content != "" {
        article.Content = content
    }

    categoryId := c.PostForm("category_id")
    if categoryId != "" {
        cid, _ := strconv.Atoi(categoryId)
        article.CategoryId = uint(cid)
    }

    userId := c.PostForm("user_id")
    if userId != "" {
        uid, _ := strconv.Atoi(userId)
        article.UserId = uint(uid)
    }

    // 执行更新操作
    if result := db.Save(&article); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": article})
}

func DeleteArticle(c *gin.Context) {
	var article Article
	db.Delete(&article, c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"data": article})
}
