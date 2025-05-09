package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Advice struct {
	Id       uint   `gorm:"primary key; autoIncrement" column:"id"`
	Username string `gorm:"type:varchar(255)" column:"username"`
	Email    string `gorm:"type:varchar(255)" column:"email"`
	Type     string `gorm:"type:varchar(255)" column:"type"`
	Content  string `gorm:"type:varchar(255)" column:"content"`
}

func (Advice) TableName() string {
	return "advices"
}

func GetAdvice(c *gin.Context) {
	username := c.Query("username")
	typeA := c.Query("type")
	var advices []Advice
	if username != "" {
		db.Where("Username = ?", username).Find(&advices)
	} else if typeA != "" {
		db.Where("Type = ?", typeA).Find(&advices)
	} else {
		db.Find(&advices)
	}
	c.JSON(http.StatusOK, gin.H{"data": advices})
}

//func GetAdvice(c *gin.Context) {
//	var advices []Advice
//	result := db.Find(&advices)
//	if result.Error != nil {
//		fmt.Printf("查询出错: %v\n", result.Error)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
//		return
//	}
//	fmt.Printf("查询到的建议数量: %d\n", len(advices))
//	c.JSON(http.StatusOK, gin.H{"data": advices})
//}

func PostAdvice(c *gin.Context) {
	var advice Advice
	if err := c.ShouldBindJSON(&advice); err == nil {
		db.Create(&advice)
		c.JSON(200, gin.H{"data": advice})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
