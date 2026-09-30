package api

import (
	"BlogBack/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Application struct {
	Id           uint   `gorm:"primaryKey;autoIncrement" json:"Id"`
	Username     string `gorm:"type:varchar(255)" json:"Username"`
	Email        string `gorm:"type:varchar(255)" json:"Email"`
	Name         string `gorm:"type:varchar(255)" json:"Name"`
	Web          string `gorm:"type:varchar(255)" json:"Web"`
	Introduction string `gorm:"type:varchar(255)" json:"Introduction"`
	Img          string `gorm:"type:varchar(255)" json:"Img"`
	Avatar       string `gorm:"type:varchar(255)" json:"Avatar"`
	Background   string `gorm:"type:varchar(255)" json:"Background"`
	Description  string `gorm:"type:varchar(255)" json:"Description"`
}

func (Application) TableName() string {
	return "applications"
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
		db.Where("id = ?", id).Find(&applications)
	} else {
		db.Order("id desc").Find(&applications)
	}
	c.JSON(http.StatusOK, gin.H{"data": applications})
}

func saveApplicationImage(c *gin.Context, field, baseDir string) (string, error) {
	file, err := c.FormFile(field)
	if err != nil {
		return "", fmt.Errorf("缺少图片字段 %s", field)
	}
	if file.Size > 2<<20 {
		return "", fmt.Errorf("%s 图片不能超过 2MB", field)
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp":
	default:
		return "", fmt.Errorf("%s 仅支持常见图片格式", field)
	}
	base := filepath.Base(file.Filename)
	base = strings.ReplaceAll(base, " ", "_")
	unique := time.Now().Format("20060102150405") + "_" + strconv.FormatInt(time.Now().UnixNano()%1e6, 10) + "_" + base
	savePath := filepath.Join(baseDir, unique)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		return "", fmt.Errorf("保存 %s 失败: %v", field, err)
	}
	return unique, nil
}

func PostApplication(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	web := strings.TrimSpace(c.PostForm("web"))
	intro := strings.TrimSpace(c.PostForm("introduction"))
	if name == "" || name == "undefined" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写网站名称"})
		return
	}
	if web == "" || web == "undefined" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写网站地址"})
		return
	}

	baseDir := utils.GetImageBaseDir()
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建图片目录失败: " + err.Error()})
		return
	}

	imgName, err := saveApplicationImage(c, "img", baseDir)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bgName, err := saveApplicationImage(c, "background", baseDir)
	if err != nil {
		// 封面已存，尽量清理
		_ = os.Remove(filepath.Join(baseDir, imgName))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clean := func(s string) string {
		s = strings.TrimSpace(s)
		if s == "undefined" || s == "null" {
			return ""
		}
		return s
	}

	application := Application{
		Name:         name,
		Username:     clean(c.PostForm("username")),
		Email:        clean(c.PostForm("email")),
		Web:          web,
		Introduction: intro,
		Img:          imgName,
		Avatar:       clean(c.PostForm("avatar")),
		Background:   bgName,
		Description:  clean(c.PostForm("description")),
	}

	if err := db.Create(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存申请失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "created successfully", "data": application})
}

func DeleteApplication(c *gin.Context) {
	id := c.Param("id")
	id1, err := strconv.ParseUint(id, 10, 64)
	if err != nil || id1 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效 ID"})
		return
	}
	if err := db.Where("id = ?", id1).Delete(&Application{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "删除成功"})
}
