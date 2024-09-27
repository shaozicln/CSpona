package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	Id           uint      `gorm:"primary key"`
	Email        string    `gorm:"type:varchar(255)" column:"email"`
	Username     string    `gorm:"type:varchar(255)" column:"username"`
	Password     string    `gorm:"type:varchar(255)" column:"password"`
	ArticleCount int       `gorm:"type:int" column:"article_count"`
	CreatedAt    time.Time `gorm:"type:timestamp" column:"created_at"`
	RoleQx       string    `gorm:"type:varchar(20)" column:"role_qx"`
	Avatar       string    `gorm:"type:varchar(255)" column:"avatar"`
	Article      []Article `gorm:"foreign key"`
	Comments     []Comment
	Messages     []Message
}

func GetUserWithUsername(c *gin.Context) {
	username := c.Query("username")
	var users []User
	if username != "" {
		db.Select("id", "username", "email", "avatar", "role_qx", "article_count", "created_at").Where("username LIKE ?", "%"+username+"%").Find(&users)
	} else {
		db.Select("id", "username", "email", "avatar", "role_qx", "article_count", "created_at").Find(&users)
	}
	for i, user := range users {
		articleCount := int64(0)
		db.Model(&Article{}).Where("user_id = ?", user.Id).Count(&articleCount)
		users[i].ArticleCount = int(articleCount)
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func PostUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err == nil {
		db.Create(&user)
		c.JSON(200, gin.H{"data": user})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
func PutUser(c *gin.Context) {
	username := c.Param("username")
	// 绑定请求体到Article结构体
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	//id1, _ := strconv.ParseUint(id, 10, 64)
	//user.Id = uint(id1)
	// 执行更新操作
	if result := db.Model(&User{}).Where("username = ?", username).Updates(User{
		Id:       user.Id,
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
		RoleQx:   user.RoleQx,
		Avatar:   user.Avatar,
	}); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": user})
	}
}
func DeleteUser(c *gin.Context) {
	var user User
	db.Delete(&user, c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func GetUserWithId(c *gin.Context) {
	id1 := c.Param("id")
	id, err := strconv.ParseUint(id1, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "无效的用户ID"})
		return
	}
	var user User
	articleCount := int64(0)
	db.Where("id=?", id).First(&user)
	db.Model(&Article{}).Where("user_id = ?", user.Id).Count(&articleCount)
	user.ArticleCount = int(articleCount)
	c.JSON(http.StatusOK, gin.H{"data": user})
}
