package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Category struct {
	Id          uint
	Name        string
	Description string
	Img         string
}

func (Category) TableName() string {
	return "categories"
}

func GetCategory(c *gin.Context) {
	name := c.Query("name")
	id := c.Query("id")
	var categories []Category
	if name != "" {
		db.Where("name LIKE ?", "%"+name+"%").Find(&categories)
	} else if id != "" {
		db.Where("id = ?", id).Find(&categories)
	} else {
		db.Find(&categories)
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func PostCatehgory(c *gin.Context) {
	var category Category
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

	category.Img = filename            //只存储文件名字
	category.Name = c.PostForm("name") //数据库其他字段和数据
	category.Description = c.PostForm("description")
	db.Create(&category)
	c.JSON(200, gin.H{"message": "created successfully", "data": category})
}

func PutCategory(c *gin.Context) {
	id := c.Param("id")

	// 绑定请求体到Article结构体
	var category Category
	if err := c.BindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	id1, _ := strconv.ParseUint(id, 10, 64)
	category.Id = uint(id1)
	// 执行更新操作
	if result := db.Model(&Category{}).Where("id = ?", id).Updates(Category{
		Id:          category.Id,
		Name:        category.Name,
		Description: category.Description,
	}); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": category})
	}
}

func DeleteCategory(c *gin.Context) {
	var category Category
	db.Delete(&category, c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"data": category})
}
