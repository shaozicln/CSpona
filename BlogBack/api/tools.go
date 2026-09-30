package api

import (
	"BlogBack/utils"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// ToolSlide 工具箱轮播单项（服务器 JSON 可直接改）
type ToolSlide struct {
	Img           string `json:"img"`
	Title         string `json:"title"`
	LayerTitle    string `json:"layerTitle"`
	Subtitle      string `json:"subtitle"`
	Description   string `json:"description"`
	RightTitle    string `json:"rightTitle"`
	RightSubtitle string `json:"rightSubtitle"`
	Link          string `json:"link"`
	// 兼容旧字段名
	ImgUrl string `json:"imgUrl,omitempty"`
	Route  string `json:"route,omitempty"`
}

// GetToolsList 公开读取工具箱列表（改服务器上的 tools-list.json 即可生效）
func GetToolsList(c *gin.Context) {
	path := utils.GetToolsListPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusOK, gin.H{
				"status":  200,
				"data":    []ToolSlide{},
				"message": "tools-list 文件不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "读取失败: " + err.Error(),
		})
		return
	}

	var items []ToolSlide
	if err := json.Unmarshal(data, &items); err != nil {
		// 兼容 { "items": [...] }
		var wrap struct {
			Items []ToolSlide `json:"items"`
		}
		if err2 := json.Unmarshal(data, &wrap); err2 != nil || len(wrap.Items) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  500,
				"message": "tools-list.json 格式错误: " + err.Error(),
			})
			return
		}
		items = wrap.Items
	}

	for i := range items {
		if strings.TrimSpace(items[i].Img) == "" && items[i].ImgUrl != "" {
			items[i].Img = items[i].ImgUrl
		}
		if strings.TrimSpace(items[i].Link) == "" && items[i].Route != "" {
			items[i].Link = items[i].Route
		}
		items[i].ImgUrl = items[i].Img
		items[i].Route = items[i].Link
	}

	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   items,
	})
}
