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
	// 获取文件
	file, err := c.FormFile("img")
	if err != nil {
		c.JSON(500, gin.H{"error": "无法获取文件"})
		return
	}

	// 定义基本目录
	baseDir := utils.GetImageBaseDir()
	// 获取上传的文件名
	filename := filepath.Base(file.Filename)
	// 生成唯一文件名
	uniqueFilename := time.Now().Format("20060102150405") + "_" + filename
	// 连接字段，形成存储路径
	savePath := filepath.Join(baseDir, uniqueFilename)
	savePath = strings.ReplaceAll(savePath, "\\", "/") // 兼容不同操作系统

	// 创建保存路径所在目录
	_ = os.MkdirAll(baseDir, os.ModePerm)
	// 保存文件
	_ = c.SaveUploadedFile(file, savePath)

	// 返回图片 URL
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
	//获取文件，错误处理
	file, err := c.FormFile("img")
	if err != nil {
		c.JSON(500, gin.H{"error": "无法获取文件"})
		return
	}
	baseDir := utils.GetImageBaseDir()         //定义基本目录
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
    
    // 先从数据库获取原始文章信息
    var article Article
    if err := db.First(&article, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
        return
    }

    // 处理文件上传（背景图）
    file, err := c.FormFile("img")
    if err == nil { // 如果有上传背景图文件
        baseDir := utils.GetImageBaseDir()
        filename := filepath.Base(file.Filename)
        savePath := filepath.Join(baseDir, filename)
        savePath = strings.ReplaceAll(savePath, "\\", "/")

        // 创建目录并保存文件
        _ = os.MkdirAll(baseDir, os.ModePerm)
        if err := c.SaveUploadedFile(file, savePath); err == nil {
            article.Img = filename // 更新背景图文件名
        }
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
