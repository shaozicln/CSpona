package routes

import (
	"BlogBack/api"
	"BlogBack/middleware"
	"BlogBack/utils"
	"github.com/gin-gonic/gin"
)

// 函数名开头大写：公用方法；小写：私有方法，只有当前包内可以应用
func InitRouter() *gin.Engine {
	gin.SetMode(utils.AppMode)
	r := gin.Default()
	r.Use(middleware.Cors())
	router := r.Group("/api")
	{
		//...登陆验证路由... Login.go
		router.POST("/login", api.Login)

		//...用户路由... User.go
		router.GET("/users/:id", api.GetUserWithId)
		router.GET("/users", api.GetUserWithUsername)
		router.POST("/users", api.PostUser)
		router.PUT("/users/:id", api.PutUser)
		router.DELETE("/users/:id", api.DeleteUser)

		// ...文章路由... Article.go
		router.GET("/search", api.SearchArticle)
		router.POST("/articles", api.PostArticle)
		router.PUT("/articles/:id", api.PutArticle)
		router.DELETE("/articles/:id", api.DeleteArticle)

		// ...图片上传路由... Article.go
		router.POST("/upload-image", api.UploadImage)

		// ...分类路由... Category.go
		router.GET("/categories", api.GetCategory)
		router.POST("/categories", api.PostCatehgory)
		router.PUT("/categories/:id", api.PutCategory)
		router.DELETE("/categories/:id", api.DeleteCategory)

		//...主页文章展示路由... CategoryWithArticle.go
		router.GET("/categories-with-articles", api.CategoryWithArticle)

		//...文章内容展示+评论路由... ContentAndComment.go
		router.GET("/path-to-article/:id", api.GetContent)
		router.POST("/path-to-article/:id", api.PostContent)
		router.DELETE("/path-to-article/:id", api.DeleteContent)
		router.GET("/comments/:articleId", api.GetComment)
		router.GET("/commentsReplies/:commentId", api.GetCommentWithReplies)
		router.PUT("/comments/:id/pin", api.PinComment)
		router.PUT("/comments/:id", api.PutComment)

		//...留言板路由... Message.go
		router.GET("/commentBoard", api.GetMessage)
		router.POST("/commentBoard", api.PostMessage)
		router.DELETE("/commentBoard/:id", api.DeleteMessage)

		//...友链网址... Friend.go
		router.GET("/friendsWeb", api.GetFriend)
		router.POST("/friendsWeb", api.PostFriend)
		router.PUT("/friendsWeb/:id", api.PutFriend)
		router.DELETE("/friendsWeb/:id", api.DeleteFriend)

		//...友链申请... Application.go
		router.GET("/application", api.GetApplication)
		router.POST("/application", api.PostApplication)
		router.DELETE("/application/:id", api.DeleteApplication)

		//...建议... Advice.go
		router.GET("/advice", api.GetAdvice)
		router.POST("/advice", api.PostAdvice)

		//...任务提醒... Reminder.go
		router.GET("/reminder/:user_id", api.GetUserReminders) // 获取指定用户的提醒
		router.POST("/reminder", api.PostReminder)
		router.PUT("/reminder/:id", api.PutReminder)
		router.DELETE("/reminder/:id", api.DeleteReminder)
		router.POST("/reminder/:id/email", api.SendReminderEmail)
	}
	return r
}
