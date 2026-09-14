package api

import (
	"BlogBack/utils"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Advice struct {
	Id         uint   `gorm:"primary key; autoIncrement" column:"id"`
	Username   string `gorm:"type:varchar(255)" column:"username"`
	Email      string `gorm:"type:varchar(255)" column:"email"`
	Type       string `gorm:"type:varchar(255)" column:"type"`
	Content    string `gorm:"type:text" column:"content"`
	MusicTitle string `gorm:"type:varchar(255);column:music_title"`
	MusicKind  string `gorm:"type:varchar(16);column:music_kind"`
	MusicRef   string `gorm:"type:varchar(512);column:music_ref"`
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

func PostAdvice(c *gin.Context) {
	var advice Advice

	ct := c.ContentType()
	if strings.HasPrefix(ct, "multipart/") {
		advice.Username = c.PostForm("username")
		advice.Email = c.PostForm("email")
		advice.Type = c.PostForm("type")
		advice.Content = c.PostForm("content")
		advice.MusicTitle = strings.TrimSpace(c.PostForm("music_title"))
		advice.MusicKind = strings.TrimSpace(c.PostForm("music_kind"))
		advice.MusicRef = strings.TrimSpace(c.PostForm("music_ref"))

		if file, err := c.FormFile("music_file"); err == nil && file != nil {
			if file.Size > 15<<20 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "音频不能超过 15MB"})
				return
			}
			filename := filepath.Base(file.Filename)
			ext := strings.ToLower(filepath.Ext(filename))
			switch ext {
			case ".mp3", ".m4a", ".flac", ".ogg", ".wav", ".aac":
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持 mp3/m4a/flac/ogg/wav/aac"})
				return
			}
			baseDir := utils.GetMusicDir()
			uniqueFilename := time.Now().Format("20060102150405") + "_" + filename
			savePath := filepath.Join(baseDir, uniqueFilename)
			if err := os.MkdirAll(baseDir, 0755); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "创建目录失败: " + err.Error()})
				return
			}
			if err := c.SaveUploadedFile(file, savePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "保存失败: " + err.Error()})
				return
			}
			advice.MusicKind = "file"
			advice.MusicRef = "music/" + uniqueFilename
		}
	} else {
		if err := c.ShouldBindJSON(&advice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if advice.Type == "背景音乐の推荐" {
		if advice.MusicTitle == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请填写歌名"})
			return
		}
		if advice.MusicKind != "netease" && advice.MusicKind != "file" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请选择网易云链接或上传文件"})
			return
		}
		if advice.MusicRef == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请填写网易云链接或上传音乐文件"})
			return
		}
		if advice.Content == "" {
			advice.Content = "推荐：" + advice.MusicTitle
		}
	}

	db.Create(&advice)
	c.JSON(200, gin.H{"data": advice})
}
