package main

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type Article struct {
	Id           uint      `gorm:"primary key"`
	Title        string    `gorm:"type:varchar(255)"`
	Content      string    `gorm:"type:text"`
	CategoryId   uint      `gorm:"type:int"`
	UserId       uint      `gorm:"type:int"`
	ViewCount    uint      `gorm:"type:int"`
	CommentCount uint      `gorm:"type:int"`
	CreatedAt    time.Time `gorm:"type:timestamp"`
	UpdatedAt    time.Time `gorm:"type:timestamp"`
	Img          string    `gorm:"type:varchar(255)"`
	User         User
	Comments     []Comment
	Category     Category
}
type Category struct {
	Id          uint
	Name        string
	Description string
	Img         string
}
type CategoryWithArticles struct {
	Category
	Articles []Article
}
type Comment struct {
	Id        uint      `gorm:"primary key; autoIncrement" column:"id"`
	ArticleId int       `gorm:"column:article_id"`
	UserId    int       `gorm:"column:user_id"`
	Idea      string    `gorm:"column:idea"`
	CreatedAt time.Time `gorm:"type:timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp"`
	User      User
}
type Message struct {
	Id        uint
	UserId    int
	Content   string
	CreatedAt time.Time
	User      User
}
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
type Friend struct {
	Id           uint   `gorm:"type:int" column:"id"`
	Name         string `gorm:"type:varchar(255)" column:"name"`
	Web          string `gorm:"type:varchar(255)" column:"web"`
	Introduction string `gorm:"type:varchar(255)" column:"introduction"`
	Img          string `gorm:"type:varchar(255)" column:"img"`
}

type Advice struct {
	Id       uint   `gorm:"primary key; autoIncrement" column:"id"`
	Username string `gorm:"type:varchar(255)" column:"username"`
	Email    string `gorm:"type:varchar(255)" column:"email"`
	Type     string `gorm:"type:varchar(255)" column:"type"`
	Content  string `gorm:"type:varchar(255)" column:"content"`
}

type Application struct {
	Id           uint   `gorm:"primary key; autoIncrement" column:"id"`
	Username     string `gorm:"type:varchar(255)" column:"username"`
	Email        string `gorm:"type:varchar(255)" column:"email"`
	Name         string `gorm:"type:varchar(255)" column:"name"`
	Web          string `gorm:"type:varchar(255)" column:"web"`
	Introduction string `gorm:"type:varchar(255)" column:"introduction"`
	Img          string `gorm:"type:varchar(255)" column:"img"`
}

func main() {
	dsn := "root:2005n04y20A!!!@tcp(127.0.0.1:3306)/kh?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("error:%v", err)
		return
	}

	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	r.Use(cors.New(corsConfig))

	//...主页文章展示路由...
	r.GET("/categories-with-articles", func(c *gin.Context) {
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
	})

	// ...文章路由...
	r.GET("/search", func(c *gin.Context) {
		title := c.Query("title")
		id := c.Query("id")
		content := c.Query("content")
		var articles []Article
		if title != "" {
			db.Where("title LIKE ?", "%"+title+"%").Find(&articles)
		} else if id != "" {
			db.Where("id = ?", id).Find(&articles)
		} else if content != "" {
			db.Where("content LIKE ?", "%"+content+"%").Find(&articles)
		} else {
			db.Find(&articles)
		}
		c.JSON(http.StatusOK, gin.H{"data": articles})
	})
	r.POST("/articles", func(c *gin.Context) {
		var article Article
		if err := c.ShouldBindJSON(&article); err == nil {
			db.Create(&article)
			c.JSON(200, gin.H{"data": article})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})
	r.PUT("/articles/:id", func(c *gin.Context) {
		id := c.Param("id")

		// 绑定请求体到Article结构体
		var article Article
		if err := c.BindJSON(&article); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		id1, _ := strconv.ParseUint(id, 10, 64)
		article.Id = uint(id1)
		// 执行更新操作
		if result := db.Model(&Article{}).Where("id = ?", id).Updates(Article{
			Id:         article.Id,
			Title:      article.Title,
			Content:    article.Content,
			CategoryId: article.CategoryId,
			UserId:     article.UserId,
		}); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"data": article})
		}
	})
	r.DELETE("/articles/:id", func(c *gin.Context) {
		var article Article
		db.Delete(&article, c.Param("id"))
		c.JSON(http.StatusOK, gin.H{"data": article})
	})

	// ...分类路由...
	r.GET("/categories", func(c *gin.Context) {
		name := c.Query("name")
		id := c.Query("id")
		var categories []Category
		if name != "" {
			db.Where("name LIKE ?", "%"+name+"%").Find(&categories)
		} else if id != "" {
			db.Where("id = ?", id).Find(&categories)
		} else {
			db.Find(&categories)
		}
		c.JSON(http.StatusOK, gin.H{"data": categories})
	})
	r.POST("/categories", func(c *gin.Context) {
		var category Category
		if err := c.ShouldBindJSON(&category); err == nil {
			db.Create(&category)
			c.JSON(200, gin.H{"data": category})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})
	r.PUT("/categories/:id", func(c *gin.Context) {
		id := c.Param("id")

		// 绑定请求体到Article结构体
		var category Category
		if err := c.BindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		id1, _ := strconv.ParseUint(id, 10, 64)
		category.Id = uint(id1)
		// 执行更新操作
		if result := db.Model(&Category{}).Where("id = ?", id).Updates(Category{
			Id:          category.Id,
			Name:        category.Name,
			Description: category.Description,
		}); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"data": category})
		}
	})
	r.DELETE("/categories/:id", func(c *gin.Context) {
		var category Category
		db.Delete(&category, c.Param("id"))
		c.JSON(http.StatusOK, gin.H{"data": category})
	})

	//...文章内容展示+评论路由...
	r.GET("/path-to-article/:id", func(c *gin.Context) {
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
	})
	r.POST("/path-to-article/:id", func(c *gin.Context) {
		id1 := c.Param("id")
		id2, err := strconv.ParseUint(id1, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "无效的文章ID"})
			return
		}
		var comment Comment
		if err := c.ShouldBindJSON(&comment); err == nil {
			comment.ArticleId = int(id2)
			db.Create(&comment)
			db.Where("article_id = ?", id2).Preload("User").Find(&comment)
			c.JSON(200, gin.H{"data": comment})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		fmt.Println(comment.ArticleId)
	})
	r.DELETE("/path-to-article/:id", func(c *gin.Context) {
		var comments []Comment
		id1 := c.Param("id")
		id, err := strconv.ParseUint(id1, 10, 64) // 解析ID，并处理可能的错误
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的ID格式"})
			return
		}

		// 检查是否存在该ID的记录
		db.Where("id=?", id).First(&comments)
		// 如果记录存在，则删除记录
		err = db.Where("id=?", id).Delete(&comments).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "删除成功"})
	})
	r.GET("/comments/:articleId", func(c *gin.Context) {
		articleIdStr := c.Param("articleId")
		articleId, err := strconv.Atoi(articleIdStr)
		if err != nil || articleId <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
			return
		}
		idea := c.Query("idea")
		username := c.Query("username")
		var comments []Comment
		if idea != "" {
			db.Where("article_id = ?  AND  idea LIKE ?", articleId, "%"+idea+"%").Preload("User").Find(&comments)
		} else if username != "" {
			var user User
			db.Where("username = ?", username).Find(&user)
			db.Where("article_id = ? AND user_id = ?", articleId, user.Id).Preload("User").Find(&comments)
		} else {
			var user User
			db.Where("username = ?", username).Find(&user)
			db.Where("article_id = ?", articleId).Preload("User").Find(&comments)
		}

		c.JSON(http.StatusOK, gin.H{"data": comments})
	})

	//...用户路由...
	r.GET("/users", func(c *gin.Context) {
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
	})
	r.POST("/users", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err == nil {
			db.Create(&user)
			c.JSON(200, gin.H{"data": user})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})
	r.PUT("/users/:username", func(c *gin.Context) {
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
	})
	r.DELETE("/users/:id", func(c *gin.Context) {
		var user User
		db.Delete(&user, c.Param("id"))
		c.JSON(http.StatusOK, gin.H{"data": user})
	})
	r.GET("/users/:id", func(c *gin.Context) {
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
	})

	//...登陆验证路由...
	r.POST("/login", func(c *gin.Context) {
		var cinuser User
		if err := c.ShouldBindJSON(&cinuser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
			return
		}
		var user User
		result := db.Where("username = ?", cinuser.Username).First(&user)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				// 用户不存在
				c.JSON(200, gin.H{"message": "喵喵喵？注册了吗就来登录？"})
			} else {
				// 查询出错
				c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
			}
			return
		}
		if user.Password == cinuser.Password {
			// 登录成功
			articleCount := int64(0)
			db.Model(&Article{}).Where("user_id = ?", user.Id).Count(&articleCount)
			user.ArticleCount = int(articleCount)
			c.JSON(http.StatusOK, gin.H{"message": "登录成功", "id": user.Id, "qx": user.RoleQx, "username": user.Username, "email": user.Email, "password": user.Password, "avatar": user.Avatar, "articleCount": user.ArticleCount})
		} else {
			// 密码错误
			c.JSON(200, gin.H{"message": "密码错误"})
		}
	})

	//...留言板路由...
	r.GET("/commentBoard", func(c *gin.Context) {
		var messages []Message
		db.Find(&messages)
		db.Preload("User").Find(&messages)
		var messageCount int64
		db.Model(&Message{}).Count(&messageCount)
		c.JSON(http.StatusOK, gin.H{
			"data":         messages,
			"messageCount": uint(messageCount),
		})
	})
	r.POST("/commentBoard", func(c *gin.Context) {
		var message Message
		if err := c.ShouldBindJSON(&message); err == nil {
			db.Create(&message)
			c.JSON(200, gin.H{"data": message})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})
	r.DELETE("/commentBoard/:id", func(c *gin.Context) {
		var messages []Message
		id1 := c.Param("id")
		id, err := strconv.ParseUint(id1, 10, 64) // 解析ID，并处理可能的错误
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的ID格式"})
			return
		}

		// 检查是否存在该ID的记录
		db.Where("id=?", id).First(&messages)
		// 如果记录存在，则删除记录
		err = db.Where("id=?", id).Delete(&messages).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "删除成功"})
	})

	//...友链网址...
	r.GET("/friendsWeb", func(c *gin.Context) {
		id := c.Query("id")
		name := c.Query("name")
		web := c.Query("web")
		var friends []Friend
		if name != "" {
			db.Where("Name LIKE ?", "%"+name+"%").Find(&friends)
		} else if web != "" {
			db.Where("Web = ?", web).Find(&friends)
		} else if id != "" {
			db.Where("Id = ?", id).Find(&friends)
		} else {
			db.Find(&friends)
		}
		c.JSON(http.StatusOK, gin.H{"data": friends})
	})
	r.POST("/friendsWeb", func(c *gin.Context) {
		var friend Friend
		if err := c.ShouldBindJSON(&friend); err == nil {
			db.Create(&friend)
			c.JSON(200, gin.H{"data": friend})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})
	r.PUT("/friendsWeb/:id", func(c *gin.Context) {
		id := c.Param("id")
		// 绑定请求体到Article结构体
		var friend Friend
		id1, _ := strconv.ParseUint(id, 10, 64)
		friend.Id = uint(id1)
		if err := c.BindJSON(&friend); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		// 执行更新操作
		if result := db.Model(&Friend{}).Where("Id = ?", id1).Updates(Friend{
			Id:           friend.Id,
			Name:         friend.Name,
			Web:          friend.Web,
			Introduction: friend.Introduction,
			Img:          friend.Img,
		}); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"data": friend})
		}
	})
	r.DELETE("/friendsWeb/:id", func(c *gin.Context) {
		var friends []Friend
		id := c.Param("id")
		id1, _ := strconv.ParseUint(id, 10, 64)
		// 检查是否存在该ID的记录
		db.Where("Id=?", id1).First(&friends)
		// 如果记录存在，则删除记录
		err = db.Where("Id=?", id1).Delete(&friends).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "删除成功"})
	})

	//...友链申请...
	r.GET("/application", func(c *gin.Context) {
		id := c.Query("id")
		username := c.Query("username")
		web := c.Query("web")
		var applications []Application
		if username != "" {
			db.Where("name LIKE ?", "%"+username+"%").Find(&applications)
		} else if web != "" {
			db.Where("web = ?", web).Find(&applications)
		} else if id != "" {
			db.Where("Id = ?", id).Find(&applications)
		} else {
			db.Find(&applications)
		}
		c.JSON(http.StatusOK, gin.H{"data": applications})
	})
	r.POST("/application", func(c *gin.Context) {
		var application Application
		if err := c.ShouldBindJSON(&application); err == nil {
			db.Create(&application)
			c.JSON(200, gin.H{"data": application})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})
	r.DELETE("/application/:id", func(c *gin.Context) {
		var applications []Application
		id := c.Param("id")
		id1, _ := strconv.ParseUint(id, 10, 64)
		// 检查是否存在该ID的记录
		db.Where("Id=?", id1).First(&applications)
		// 如果记录存在，则删除记录
		err = db.Where("Id=?", id1).Delete(&applications).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "删除成功"})
	})

	//...建议...
	r.GET("/advice", func(c *gin.Context) {
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
	})
	r.POST("/advice", func(c *gin.Context) {
		var advice Advice
		if err := c.ShouldBindJSON(&advice); err == nil {
			db.Create(&advice)
			c.JSON(200, gin.H{"data": advice})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.Run(":8081")
}
