package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	if err := c.ShouldBindJSON(&category); err == nil {
		db.Create(&category)
		c.JSON(200, gin.H{"data": category})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
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
