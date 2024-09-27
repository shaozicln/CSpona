package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	if err := c.ShouldBindJSON(&application); err == nil {
		db.Create(&application)
		c.JSON(200, gin.H{"data": application})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
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
