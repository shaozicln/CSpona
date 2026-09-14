package api

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"regexp"
	"strings"
	"time"

	"BlogBack/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const visitorCookieName = "cspona_vid"

var visitorIDRe = regexp.MustCompile(`(?i)^[a-f0-9-]{8,64}$`)

// ArticleView 文章浏览去重
// viewer_key: 登录用 "u:{userId}"，匿名用 "v:{uuid}"，兜底 "i:{ipUaHash}"
type ArticleView struct {
	Id        uint      `gorm:"primaryKey;autoIncrement"`
	ArticleId uint      `gorm:"column:article_id;uniqueIndex:uk_article_viewer;not null"`
	ViewerKey string    `gorm:"column:viewer_key;type:varchar(64);uniqueIndex:uk_article_viewer;not null"`
	UserId    *uint     `gorm:"column:user_id;index"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime"`
}

func (ArticleView) TableName() string {
	return "article_views"
}

// applyArticleView 按访客身份去重并可能 +1 浏览量；返回最新 ViewCount
func applyArticleView(c *gin.Context, article *Article) {
	if article == nil || article.Id == 0 {
		return
	}
	key, uid := resolveViewer(c)
	if key == "" {
		return
	}
	if recordArticleView(article.Id, key, uid) {
		article.ViewCount++
		return
	}
	var fresh Article
	if db.Select("view_count").First(&fresh, article.Id).Error == nil {
		article.ViewCount = fresh.ViewCount
	}
}

func resolveViewer(c *gin.Context) (viewerKey string, userID *uint) {
	if uid, ok := middleware.GetUserID(c); ok && uid > 0 {
		id := uid
		return "u:" + itoaUint(uid), &id
	}

	vid := strings.TrimSpace(c.GetHeader("X-Visitor-Id"))
	if vid == "" {
		if cookie, err := c.Cookie(visitorCookieName); err == nil {
			vid = strings.TrimSpace(cookie)
		}
	}
	if visitorIDRe.MatchString(vid) {
		return "v:" + strings.ToLower(vid), nil
	}

	// 服务端签发匿名 ID（HttpOnly Cookie，1 年）
	vid = newVisitorID()
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     visitorCookieName,
		Value:    vid,
		Path:     "/",
		MaxAge:   365 * 24 * 3600,
		HttpOnly: false, // 前端也可读，便于带 X-Visitor-Id
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	})
	return "v:" + vid, nil
}

func recordArticleView(articleID uint, viewerKey string, userID *uint) bool {
	if articleID == 0 || viewerKey == "" {
		return false
	}
	var existing ArticleView
	err := db.Where("article_id = ? AND viewer_key = ?", articleID, viewerKey).First(&existing).Error
	if err == nil {
		return false
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return false
	}
	v := ArticleView{
		ArticleId: articleID,
		ViewerKey: viewerKey,
		UserId:    userID,
	}
	if createErr := db.Create(&v).Error; createErr != nil {
		return false
	}
	_ = db.Model(&Article{}).Where("id = ?", articleID).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
	return true
}

func newVisitorID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		sum := sha256.Sum256([]byte(time.Now().String()))
		return hex.EncodeToString(sum[:16])
	}
	// uuid-like
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	hexs := hex.EncodeToString(b)
	return hexs[0:8] + "-" + hexs[8:12] + "-" + hexs[12:16] + "-" + hexs[16:20] + "-" + hexs[20:32]
}

func itoaUint(n uint) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
