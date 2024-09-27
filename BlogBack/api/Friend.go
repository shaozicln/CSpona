package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Friend struct {
	Id           uint   `gorm:"type:int" column:"id"`
	Name         string `gorm:"type:varchar(255)" column:"name"`
	Web          string `gorm:"type:varchar(255)" column:"web"`
	Introduction string `gorm:"type:varchar(255)" column:"introduction"`
	Img          string `gorm:"type:varchar(255)" column:"img"`
}

func GetFriend(c *gin.Context) {
	id := c.Query("id")
	name := c.Query("name")
	web := c.Query("web")
	var friends []Friend
	if name != "" {
		db.Where("Name LIKE ?", "%"+name+"%").Find(&friends)
	} else if web != "" {
		db.Where("Web = ?", web).Find(&friends)
	} else if id != "" {
		db.Where("Id = ?", id).Find(&friends)
	} else {
		db.Find(&friends)
	}
	c.JSON(http.StatusOK, gin.H{"data": friends})
}

func PostFriend(c *gin.Context) {
	var friend Friend
	if err := c.ShouldBindJSON(&friend); err == nil {
		db.Create(&friend)
		c.JSON(200, gin.H{"data": friend})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func PutFriend(c *gin.Context) {
	id := c.Param("id")
	// 绑定请求体到Article结构体
	var friend Friend
	id1, _ := strconv.ParseUint(id, 10, 64)
	friend.Id = uint(id1)
	if err := c.BindJSON(&friend); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	// 执行更新操作
	if result := db.Model(&Friend{}).Where("Id = ?", id1).Updates(Friend{
		Id:           friend.Id,
		Name:         friend.Name,
		Web:          friend.Web,
		Introduction: friend.Introduction,
		Img:          friend.Img,
	}); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": friend})
	}
}

func DeleteFriend(c *gin.Context) {
	var friends []Friend
	id := c.Param("id")
	id1, _ := strconv.ParseUint(id, 10, 64)
	// 检查是否存在该ID的记录
	db.Where("Id=?", id1).First(&friends)
	// 如果记录存在，则删除记录
	err = db.Where("Id=?", id1).Delete(&friends).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "删除成功"})
}
