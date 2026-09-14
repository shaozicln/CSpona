package api

import (
	"BlogBack/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
)

type MusicTrack struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type MusicSettings struct {
	Mode       string       `json:"mode"`       // netease(歌单) | song(单曲) | files
	PlaylistId string       `json:"playlistId"` // 网易云歌单或单曲 ID
	Tracks     []MusicTrack `json:"tracks"`     // 本地曲目，url 为 music/xxx.mp3
}

func defaultMusicSettings() MusicSettings {
	return MusicSettings{
		Mode:       "netease",
		PlaylistId: "",
		Tracks:     []MusicTrack{},
	}
}

func loadMusicSettings() (MusicSettings, error) {
	path := utils.GetMusicSettingsPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultMusicSettings(), nil
		}
		return MusicSettings{}, err
	}
	var s MusicSettings
	if err := json.Unmarshal(data, &s); err != nil {
		return MusicSettings{}, err
	}
	if s.Mode == "" {
		s.Mode = "netease"
	}
	if s.Tracks == nil {
		s.Tracks = []MusicTrack{}
	}
	return s, nil
}

func saveMusicSettings(s MusicSettings) error {
	path := utils.GetMusicSettingsPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// GetMusicSettings 公开读取站点默认背景音乐
func GetMusicSettings(c *gin.Context) {
	s, err := loadMusicSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "读取失败: " + err.Error(),
		})
		return
	}
	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   s,
	})
}

// PutMusicSettings 登录后保存站点默认背景音乐
func PutMusicSettings(c *gin.Context) {
	var body MusicSettings
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "参数错误",
		})
		return
	}
	body.Mode = strings.TrimSpace(body.Mode)
	if body.Mode != "netease" && body.Mode != "song" && body.Mode != "files" {
		body.Mode = "netease"
	}
	body.PlaylistId = strings.TrimSpace(body.PlaylistId)
	if body.Tracks == nil {
		body.Tracks = []MusicTrack{}
	}
	if err := saveMusicSettings(body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "保存失败: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "保存成功",
		"data":    body,
	})
}

func metingAPIBases() []string {
	bases := []string{
		"https://api.injahow.cn/meting",
		"https://api.qijieya.cn/meting",
		"https://api.obdo.cc/meting",
	}
	if cfg, err := ini.Load("config.ini"); err == nil {
		if v := strings.TrimSpace(cfg.Section("music").Key("MetingAPI").String()); v != "" {
			custom := strings.TrimRight(v, "/?")
			out := []string{custom}
			for _, b := range bases {
				if b != custom {
					out = append(out, b)
				}
			}
			return out
		}
	}
	return bases
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func metingLooksOK(status int, body []byte) bool {
	if status != http.StatusOK || len(body) < 2 {
		return false
	}
	s := strings.TrimSpace(string(body))
	if !(strings.HasPrefix(s, "{") || strings.HasPrefix(s, "[")) {
		return false
	}
	return strings.Contains(s, "url")
}

// ProxyMeting 代理网易云解析，避免浏览器 CORS；多源回退
// type=song|playlist 返回 JSON；type=url 跟随跳转代理音频流
func ProxyMeting(c *gin.Context) {
	typ := strings.TrimSpace(c.Query("type"))
	id := strings.TrimSpace(c.Query("id"))
	if typ != "song" && typ != "playlist" && typ != "url" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type 须为 song / playlist / url"})
		return
	}
	if id == "" || !isDigits(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 id"})
		return
	}

	if typ == "url" {
		proxyMetingAudio(c, id)
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	var lastStatus int
	var lastBody []byte
	var lastErr string

	for _, base := range metingAPIBases() {
		url := fmt.Sprintf("%s/?server=netease&type=%s&id=%s", base, typ, id)
		resp, err := client.Get(url)
		if err != nil {
			lastErr = err.Error()
			continue
		}
		body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
		resp.Body.Close()
		if err != nil {
			lastErr = err.Error()
			continue
		}
		lastStatus = resp.StatusCode
		lastBody = body
		if metingLooksOK(resp.StatusCode, body) {
			c.Header("Cache-Control", "public, max-age=300")
			c.Header("X-Meting-Upstream", base)
			c.Data(http.StatusOK, "application/json; charset=utf-8", body)
			return
		}
	}

	if lastBody != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":  "网易云解析源均不可用",
			"status": lastStatus,
			"detail": string(lastBody),
		})
		return
	}
	c.JSON(http.StatusBadGateway, gin.H{
		"error":  "网易云解析源均不可用",
		"detail": lastErr,
	})
}

func proxyMetingAudio(c *gin.Context, id string) {
	client := &http.Client{
		Timeout: 20 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 8 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
	var lastErr string
	for _, base := range metingAPIBases() {
		url := fmt.Sprintf("%s/?server=netease&type=url&id=%s", base, id)
		resp, err := client.Get(url)
		if err != nil {
			lastErr = err.Error()
			continue
		}
		ct := resp.Header.Get("Content-Type")
		// 拿到音频或成功的二进制流
		if resp.StatusCode >= 200 && resp.StatusCode < 400 &&
			(strings.Contains(ct, "audio") || strings.Contains(ct, "octet-stream") || strings.Contains(ct, "mpeg") || resp.ContentLength > 1024) {
			defer resp.Body.Close()
			c.Header("Cache-Control", "public, max-age=300")
			c.Header("X-Meting-Upstream", base)
			if ct != "" {
				c.Header("Content-Type", ct)
			} else {
				c.Header("Content-Type", "audio/mpeg")
			}
			c.Status(http.StatusOK)
			_, _ = io.Copy(c.Writer, io.LimitReader(resp.Body, 30<<20))
			return
		}
		// 若是短跳转/文本失败则试下一个源
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		resp.Body.Close()
		lastErr = fmt.Sprintf("%s status=%d ct=%s body=%s", base, resp.StatusCode, ct, string(body))
	}
	c.JSON(http.StatusBadGateway, gin.H{
		"error":  "音频地址解析失败",
		"detail": lastErr,
	})
}

// UploadMusic 鉴权上传音频到 Pictures/music/
func UploadMusic(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法获取文件: " + err.Error()})
		return
	}
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

	rel := "music/" + uniqueFilename
	c.JSON(http.StatusOK, gin.H{
		"status":   200,
		"path":     rel,
		"imageUrl": "/Pictures/" + rel,
		"url":      "/Pictures/" + rel,
	})
}
