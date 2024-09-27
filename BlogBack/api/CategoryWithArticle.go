package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CategoryWithArticles struct {
	Category
	Articles []Article
}

func CategoryWithArticle(c *gin.Context) {
	var categories []Category
	db.Find(&categories)
	// 使用 map 来组织每个分类下的文章
	categoryMap := make(map[uint][]Article)
	for _, category := range categories {
		var articles []Article
		db.Where("category_id = ?", category.Id).Find(&articles)
		categoryMap[category.Id] = articles
	}
	// 将 map 转换为 CategoryWithArticles 切片
	var categoriesWithArticles []CategoryWithArticles
	for _, category := range categories {
		articles, ok := categoryMap[category.Id]
		if !ok {
			articles = []Article{}
		}
		categoriesWithArticles = append(categoriesWithArticles, CategoryWithArticles{
			Category: category,
			Articles: articles,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": categoriesWithArticles})
}
