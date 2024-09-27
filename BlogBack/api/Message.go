package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Message struct {
	Id        uint
	UserId    int
	Content   string
	CreatedAt time.Time
	User      User
}

func GetMessage(c *gin.Context) {
	var messages []Message
	db.Find(&messages)
	db.Preload("User").Find(&messages)
	var messageCount int64
	db.Model(&Message{}).Count(&messageCount)
	c.JSON(http.StatusOK, gin.H{
		"data":         messages,
		"messageCount": uint(messageCount),
	})
}

func PostMessage(c *gin.Context) {
	var message Message
	if err := c.ShouldBindJSON(&message); err == nil {
		db.Create(&message)
		c.JSON(200, gin.H{"data": message})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func DeleteMessage(c *gin.Context) {
	var messages []Message
	id1 := c.Param("id")
	id, err := strconv.ParseUint(id1, 10, 64) // 解析ID，并处理可能的错误
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的ID格式"})
		return
	}

	// 检查是否存在该ID的记录
	db.Where("id=?", id).First(&messages)
	// 如果记录存在，则删除记录
	err = db.Where("id=?", id).Delete(&messages).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "删除成功"})
}
