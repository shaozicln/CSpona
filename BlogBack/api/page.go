package api

import (
	"BlogBack/utils"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ---- 页面聚合 DTO（字段名保持与现有前端大写风格兼容） ----

type pageArticleCard struct {
	Id        uint      `json:"Id"`
	Title     string    `json:"Title"`
	Img       string    `json:"Img"`
	ImgUrl    string    `json:"ImgUrl"`
	UserId    uint      `json:"UserId"`
	CreatedAt time.Time `json:"CreatedAt"`
	ViewCount uint      `json:"ViewCount"`
}

type pageCategoryBlock struct {
	Id          uint              `json:"Id"`
	Name        string            `json:"Name"`
	Description string            `json:"Description"`
	Img         string            `json:"Img"`
	ImgUrl      string            `json:"ImgUrl"`
	Articles    []pageArticleCard `json:"Articles"`
}

type pageAuthor struct {
	Id           uint   `json:"Id"`
	Username     string `json:"Username"`
	Avatar       string `json:"Avatar"`
	AvatarUrl    string `json:"AvatarUrl"`
	Description  string `json:"Description"`
	ArticleCount int    `json:"ArticleCount"`
}

type pageCommentUser struct {
	Id        uint   `json:"Id"`
	Username  string `json:"Username"`
	Avatar    string `json:"Avatar"`
	AvatarUrl string `json:"AvatarUrl"`
}

type pageComment struct {
	Id        uint             `json:"Id"`
	ArticleId int              `json:"ArticleId"`
	UserId    int              `json:"UserId"`
	Idea      string           `json:"Idea"`
	ParentId  *uint            `json:"ParentId"`
	ReplyId   *int             `json:"ReplyId"`
	IsPinned  int              `json:"IsPinned"`
	PinnedAt  *time.Time       `json:"PinnedAt"`
	CreatedAt time.Time        `json:"CreatedAt"`
	UpdatedAt time.Time        `json:"UpdatedAt"`
	User      pageCommentUser  `json:"User"`
	Replies   []pageComment    `json:"Replies"`
}

type pageArticleDetail struct {
	Id           uint      `json:"Id"`
	Title        string    `json:"Title"`
	Content      string    `json:"Content"`
	CategoryId   uint      `json:"CategoryId"`
	UserId       uint      `json:"UserId"`
	ViewCount    uint      `json:"ViewCount"`
	CommentCount uint      `json:"CommentCount"`
	CreatedAt    time.Time `json:"CreatedAt"`
	UpdatedAt    time.Time `json:"UpdatedAt"`
	Img          string    `json:"Img"`
	ImgUrl       string    `json:"ImgUrl"`
	CategoryName string    `json:"CategoryName"`
	MusicType    string    `json:"MusicType"`
	MusicRef     string    `json:"MusicRef"`
}

// PageHome 文章列表页聚合：分类 + 文章 + 已解析封面 URL
func PageHome(c *gin.Context) {
	var categories []Category
	if err := db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "加载分类失败", "data": nil})
		return
	}

	blocks := make([]pageCategoryBlock, 0, len(categories))
	for _, category := range categories {
		var articles []Article
		db.Where("category_id = ?", category.Id).Find(&articles)
		sort.Slice(articles, func(i, j int) bool {
			if category.Id == 1000 {
				pi := isWanderlandPinned(articles[i])
				pj := isWanderlandPinned(articles[j])
				if pi != pj {
					return pi
				}
			}
			return articles[i].CreatedAt.After(articles[j].CreatedAt)
		})

		cards := make([]pageArticleCard, 0, len(articles))
		for _, a := range articles {
			cards = append(cards, pageArticleCard{
				Id:        a.Id,
				Title:     a.Title,
				Img:       a.Img,
				ImgUrl:    utils.ResolveImageURL(a.Img, ""),
				UserId:    a.UserId,
				CreatedAt: a.CreatedAt,
				ViewCount: a.ViewCount,
			})
		}

		blocks = append(blocks, pageCategoryBlock{
			Id:          category.Id,
			Name:        category.Name,
			Description: category.Description,
			Img:         category.Img,
			ImgUrl:      utils.ResolveImageURL(category.Img, ""),
			Articles:    cards,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": gin.H{"categories": blocks},
	})
}

// PageArticle 文章详情聚合：正文 + TOC + 作者 + 全量评论（图片 URL 已解析）
func PageArticle(c *gin.Context) {
	id := c.Param("id")
	articleID, err := strconv.Atoi(id)
	if err != nil || articleID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "无效的文章ID", "data": nil})
		return
	}

	var article Article
	if err := db.Preload("User").First(&article, articleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "msg": "文章未找到", "data": nil})
		return
	}

	// 浏览量：登录按 user，未登录按访客 ID，同一身份只计一次
	applyArticleView(c, &article)

	var articleCount int64
	db.Model(&Article{}).Where("user_id = ?", article.UserId).Count(&articleCount)

	var commentCount int64
	db.Model(&Comment{}).Where("article_id = ?", articleID).Count(&commentCount)

	var category Category
	db.Where("id = ?", article.CategoryId).Find(&category)

	var comments []Comment
	db.Where("article_id = ? AND parent_id IS NULL", articleID).
		Preload("User").
		Preload("Replies").
		Preload("Replies.User").
		Find(&comments)

	pageComments := mapComments(comments)
	sortComments(pageComments)

	toc := utils.BuildMarkdownTOC(article.Content)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": gin.H{
			"article": pageArticleDetail{
				Id:           article.Id,
				Title:        article.Title,
				Content:      article.Content,
				CategoryId:   article.CategoryId,
				UserId:       article.UserId,
				ViewCount:    article.ViewCount,
				CommentCount: uint(commentCount),
				CreatedAt:    article.CreatedAt,
				UpdatedAt:    article.UpdatedAt,
				Img:          article.Img,
				ImgUrl:       utils.ResolveImageURL(article.Img, ""),
				CategoryName: category.Name,
				MusicType:    article.MusicType,
				MusicRef:     article.MusicRef,
			},
			"toc": toc,
			"author": pageAuthor{
				Id:           article.User.Id,
				Username:     article.User.Username,
				Avatar:       article.User.Avatar,
				AvatarUrl:    utils.ResolveImageURL(article.User.Avatar, "boli.jpg"),
				Description:  article.User.Description,
				ArticleCount: int(articleCount),
			},
			"category": gin.H{
				"Id":   category.Id,
				"Name": category.Name,
			},
			"comments": pageComments,
		},
	})
}

func mapComments(list []Comment) []pageComment {
	out := make([]pageComment, 0, len(list))
	for _, cmt := range list {
		item := mapOneComment(cmt)
		if len(cmt.Replies) > 0 {
			item.Replies = mapComments(cmt.Replies)
			sortComments(item.Replies)
		} else {
			item.Replies = []pageComment{}
		}
		out = append(out, item)
	}
	return out
}

func mapOneComment(cmt Comment) pageComment {
	return pageComment{
		Id:        cmt.Id,
		ArticleId: cmt.ArticleId,
		UserId:    cmt.UserId,
		Idea:      cmt.Idea,
		ParentId:  cmt.ParentId,
		ReplyId:   cmt.ReplyId,
		IsPinned:  cmt.IsPinned,
		PinnedAt:  cmt.PinnedAt,
		CreatedAt: cmt.CreatedAt,
		UpdatedAt: cmt.UpdatedAt,
		User: pageCommentUser{
			Id:        cmt.User.Id,
			Username:  cmt.User.Username,
			Avatar:    cmt.User.Avatar,
			AvatarUrl: utils.ResolveImageURL(cmt.User.Avatar, "boli.jpg"),
		},
		Replies: []pageComment{},
	}
}

func sortComments(list []pageComment) {
	sort.SliceStable(list, func(i, j int) bool {
		a, b := list[i], list[j]
		if a.IsPinned == 1 && b.IsPinned != 1 {
			return true
		}
		if a.IsPinned != 1 && b.IsPinned == 1 {
			return false
		}
		if a.IsPinned == 1 && b.IsPinned == 1 {
			at := time.Time{}
			bt := time.Time{}
			if a.PinnedAt != nil {
				at = *a.PinnedAt
			}
			if b.PinnedAt != nil {
				bt = *b.PinnedAt
			}
			return at.After(bt)
		}
		return a.CreatedAt.After(b.CreatedAt)
	})
}

// 漫游地置顶：开放公告（历史 ID=95），标题兜底匹配
func isWanderlandPinned(a Article) bool {
	if a.Id == 95 {
		return true
	}
	return strings.Contains(a.Title, "正式开放") && strings.Contains(a.Title, "Wanderland")
}
