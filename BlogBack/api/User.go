package api

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/scrypt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"BlogBack/utils"
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

func (User) TableName() string {
	return "users"
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
	//获取文件，错误处理
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(500, gin.H{"error": "无法获取文件"})
		return
	}
	baseDir := utils.GetImageBaseDir()       //定义基本目录
	filename := filepath.Base(file.Filename)           //获取上传的文件名
	savePath := filepath.Join(baseDir, filename)       //连接字段，形成存储路径
	savePath = strings.ReplaceAll(savePath, "\\", "/") // 将路径用正斜杠保存，兼容不同操作系统，同时方便前端读取

	_ = os.MkdirAll(baseDir, os.ModePerm) //创建保存路径所在目录，目录存在则忽略，同时忽略了错误处理
	_ = c.SaveUploadedFile(file, savePath)

	user.Avatar = filename                 //只存储文件名字
	user.Username = c.PostForm("username") //数据库其他字段和数据
	user.Email = c.PostForm("email")
	user.Password = ScryptPw(c.PostForm("password")) //密码加密储存
	user.RoleQx = c.PostForm("role_qx")

	db.Create(&user)
	c.JSON(200, gin.H{"message": "created successfully", "data": user})
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

// 密码加密
func ScryptPw(password string) string {
	const KeyLen = 10
	salt := make([]byte, 8)
	salt = []byte{12, 56, 45, 89, 52, 123, 45, 96}

	HashPw, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, KeyLen)
	if err != nil {
		log.Fatal(err)
	}
	fpw := base64.StdEncoding.EncodeToString(HashPw)
	return fpw
}
