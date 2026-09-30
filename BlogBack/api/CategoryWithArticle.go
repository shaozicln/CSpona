package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type CategoryWithArticles struct {
	Category
	Articles []Article
}

// CategoryWithArticle 兼容旧接口；默认与 PageHome(articles) 一致，排除测试/漫游地
// scope=wanderland 仅漫游地；scope=all 返回全部（管理用）
func CategoryWithArticle(c *gin.Context) {
	scope := strings.TrimSpace(c.Query("scope"))
	if scope == "" {
		scope = "articles"
	}

	var categories []Category
	db.Find(&categories)

	categoryMap := make(map[uint][]Article)
	for _, category := range categories {
		if scope == "wanderland" {
			if !isWanderlandCategory(category) {
				continue
			}
		} else if scope != "all" && isArticleListExcludedCategory(category) {
			continue
		}
		var articles []Article
		db.Where("category_id = ?", category.Id).Find(&articles)
		categoryMap[category.Id] = articles
	}

	var categoriesWithArticles []CategoryWithArticles
	for _, category := range categories {
		articles, ok := categoryMap[category.Id]
		if !ok {
			continue
		}
		categoriesWithArticles = append(categoriesWithArticles, CategoryWithArticles{
			Category: category,
			Articles: articles,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": categoriesWithArticles})
}
