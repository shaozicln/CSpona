package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)
type Comment struct {
	Id        uint      `gorm:"primary key; autoIncrement" column:"id"`
	ArticleId int       `gorm:"column:article_id"`
	UserId    int       `gorm:"column:user_id"`
	Idea      string    `gorm:"column:idea"`
	ParentId  *uint     `gorm:"column:parent_id"`
	CreatedAt time.Time `gorm:"type:timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp"`
	User      User      `gorm:"foreignKey:UserId"`
	Replies   []Comment `gorm:"foreignKey:ParentId"` // 子评论关联（基于parent_id）
	IsPinned  int       `gorm:"column:is_pinned"`
	ReplyId   *int       `gorm:"column:reply_id"` // 新增：指向被直接回复的评论ID
    PinnedAt *time.Time `gorm:"column:pinned_at;type:timestamp;null"` // 置顶时间
}

func (Comment) TableName() string {
	return "comments"
}

// 置顶评论API
func PinComment(c *gin.Context) {
    // 获取评论ID
    commentIDStr := c.Param("id")
    commentID, err := strconv.Atoi(commentIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的评论ID"})
        return
    }

    // 解析请求体（保持原有结构，仅修改字段）
    var request struct {
        IsPinned int `json:"is_pinned"`
    }
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
        return
    }

    // 查询评论是否存在
    var comment Comment
    if err := db.First(&comment, commentID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
        return
    }

    // 更新置顶状态和时间（保持单条Update语句风格）
    if request.IsPinned == 1 {
        // 置顶操作：设置 is_pinned 和 pinned_at
        if err := db.Model(&comment).Updates(map[string]interface{}{
            "is_pinned":  request.IsPinned,
            "pinned_at":  time.Now(),
        }).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "更新置顶状态失败"})
            return
        }
    } else {
        // 取消置顶：只更新 is_pinned，pinned_at 设为 NULL
        if err := db.Model(&comment).Updates(map[string]interface{}{
            "is_pinned":  request.IsPinned,
            "pinned_at":  nil,
        }).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "更新置顶状态失败"})
            return
        }
    }

    // 返回结果（保持原有结构，添加 pinned_at 字段）
    c.JSON(http.StatusOK, gin.H{
        "message": "操作成功",
        "data": gin.H{
            "comment_id": commentID,
            "is_pinned":  request.IsPinned,
            "pinned_at":  comment.PinnedAt, // 返回更新后的时间（可能为 nil）
        },
    })
}

func GetContent(c *gin.Context) {
	id := c.Param("id")
	var article Article
	if err := db.Preload("User").Preload("Comments").Preload("Comments.User").Find(&article, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
		return
	}
	db.Where("id = ?", id).First(&article)
	article.ViewCount++
	db.Save(&article)
	var articleCount int64
	db.Model(&Article{}).Where("user_id = ?", article.UserId).Count(&articleCount)
	article.User.ArticleCount = int(articleCount)
	var commentCount int64
	db.Model(&Comment{}).Where("article_id = ?", id).Count(&commentCount)
	article.CommentCount = uint(commentCount)
	var category Category
	db.Model(&Category{}).Where("id = ?", article.CategoryId).Find(&category)
	article.Category.Name = category.Name
	c.JSON(http.StatusOK, gin.H{"data": article})
}

func PostContent(c *gin.Context) {
    id1 := c.Param("id")
    id2, err := strconv.ParseUint(id1, 10, 64)
    if err != nil {
        c.JSON(400, gin.H{"error": "无效的文章ID"})
        return
    }
    
    var comment Comment
    if err := c.ShouldBindJSON(&comment); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误：" + err.Error()})
        return
    }
    
    comment.ArticleId = int(id2)
    
    if comment.ParentId != nil { // 回复评论
        var parentComment Comment
        if err := db.Where("id = ? AND article_id = ?", *comment.ParentId, comment.ArticleId).First(&parentComment).Error; err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "父评论不存在或不属于当前文章"})
            return
        }

        if parentComment.ParentId != nil {
            comment.ParentId = parentComment.ParentId
            comment.ReplyId = intPtr(int(parentComment.Id)) // 指向被回复的子评论ID
        } else {
            comment.ReplyId = intPtr(int(parentComment.Id)) // 指向被回复的主评论ID
        }
    } else { // 主评论（无 parentId）
        comment.ParentId = nil  
        comment.ReplyId = nil   // 关键修复：主评论无回复目标，设为 NULL
    }
    
    // 创建评论（捕获数据库错误）
    if err := db.Create(&comment).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评论失败：" + err.Error()})
        return
    }
    
    // 预加载用户信息
    var createdComment Comment
    if err := db.Where("id = ?", comment.Id).Preload("User").First(&createdComment).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论信息失败：" + err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"data": createdComment})
}

// 辅助函数：将 int 转为 *int（用于给 ReplyId 赋值）
func intPtr(i int) *int {
    return &i
}
func DeleteContent(c *gin.Context) {
    id1 := c.Param("id")
    id, err := strconv.ParseUint(id1, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的ID格式"})
        return
    }

    // 由于设置了级联删除，直接删除主评论即可
    err = db.Where("id = ?", id).Delete(&Comment{}).Error
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除失败"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"success": true, "message": "删除成功"})
}

func GetComment(c *gin.Context) {
    articleIdStr := c.Param("articleId")
    articleId, err := strconv.Atoi(articleIdStr)
    if err != nil || articleId <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
        return
    }
    
    idea := c.Query("idea")
    username := c.Query("username")
    
    var comments []Comment
    query := db.Where("article_id = ? AND parent_id IS NULL", articleId)
    
    if idea != "" {
        query = query.Where("idea LIKE ?", "%"+idea+"%")
    } else if username != "" {
        var user User
        db.Where("username = ?", username).Find(&user)
        query = query.Where("user_id = ?", user.Id)
    }
    
    // 预加载用户信息和回复（包括回复的用户信息）
    query.Preload("User").
        Preload("Replies").
        Preload("Replies.User").
        Find(&comments)
    
    c.JSON(http.StatusOK, gin.H{"data": comments})
}

func GetCommentWithReplies(c *gin.Context) {
    commentIdStr := c.Param("commentId")
    commentId, err := strconv.ParseUint(commentIdStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的评论ID"})
        return
    }
    
    var comment Comment
    if err := db.Where("id = ?", commentId).
        Preload("User").
        Preload("Replies").
        Preload("Replies.User").
        First(&comment).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "评论未找到"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"data": comment})
}

// 更新评论内容API
func PutComment(c *gin.Context) {
    // 获取评论ID
    commentIDStr := c.Param("id")
    commentID, err := strconv.Atoi(commentIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的评论ID"})
        return
    }

    // 解析请求体
    var request struct {
        Idea string `json:"idea"`
    }
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
        return
    }

    // 验证评论内容
    if request.Idea == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "评论内容不能为空"})
        return
    }

    // 查询评论是否存在
    var comment Comment
    if err := db.First(&comment, commentID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
        return
    }

    // 保存原始评论内容用于返回
    originalIdea := comment.Idea

    // 更新评论内容
    if err := db.Model(&comment).Updates(map[string]interface{}{
        "idea":       request.Idea,
        "updated_at": time.Now(),
    }).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "更新评论失败"})
        return
    }

    // 预加载用户信息
    db.Preload("User").First(&comment, commentID)

    // 返回结果（包含原始内容）
    c.JSON(http.StatusOK, gin.H{
        "message": "更新成功",
        "data": gin.H{
            "comment":       comment,
            "original_idea": originalIdea,
        },
    })
}
