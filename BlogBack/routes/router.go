package routes

import (
	"BlogBack/api"
	"BlogBack/middleware"
	"BlogBack/utils"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(utils.AppMode)
	r := gin.Default()
	r.Use(middleware.Cors())
	router := r.Group("/api")
	{
		// 公开：登录 / 注册 / 读接口 / 页面聚合
		router.POST("/login", api.Login)
		router.POST("/logout", api.Logout)
		router.POST("/users", api.PostUser) // 注册
		router.POST("/password-reset", api.ResetPassword)

		router.GET("/users/:id", api.GetUserWithId)
		router.GET("/users", api.GetUserWithUsername)

		router.GET("/search", api.SearchArticle)
		router.GET("/categories", api.GetCategory)
		router.GET("/categories-with-articles", api.CategoryWithArticle)
		router.GET("/page/home", api.PageHome)
		router.GET("/page/article/:id", middleware.AuthOptional(), api.PageArticle)

		router.GET("/path-to-article/:id", middleware.AuthOptional(), api.GetContent)
		router.GET("/comments/:articleId", api.GetComment)
		router.GET("/commentsReplies/:commentId", api.GetCommentWithReplies)

		router.GET("/commentBoard", api.GetMessage)
		router.GET("/friendsWeb", api.GetFriend)
		router.GET("/application", api.GetApplication)
		router.POST("/application", api.PostApplication) // 友链申请对访客开放
		router.GET("/advice", api.GetAdvice)
		router.POST("/advice", api.PostAdvice) // 反馈对访客开放
		router.GET("/about-me", api.GetAboutMe)
		router.GET("/music/settings", api.GetMusicSettings)
		router.GET("/music/meting", api.ProxyMeting)

		// 会话探测：未登录也返回 200，避免访客浏览时控制台 401
		router.GET("/auth/me", middleware.AuthOptional(), api.AuthMe)

		// 需登录（写操作 / 管理）
		auth := router.Group("/")
		auth.Use(middleware.AuthRequired())
		{
			auth.PUT("/users/:id", api.PutUser)
			auth.DELETE("/users/:id", api.DeleteUser)

			auth.POST("/articles", api.PostArticle)
			auth.PUT("/articles/:id", api.PutArticle)
			auth.DELETE("/articles/:id", api.DeleteArticle)
			auth.POST("/upload-image", api.UploadImage)
			auth.POST("/music/upload", api.UploadMusic)
			auth.PUT("/music/settings", api.PutMusicSettings)

			auth.POST("/categories", api.PostCatehgory)
			auth.PUT("/categories/:id", api.PutCategory)
			auth.DELETE("/categories/:id", api.DeleteCategory)

			auth.POST("/path-to-article/:id", api.PostContent)
			auth.DELETE("/path-to-article/:id", api.DeleteContent)
			auth.PUT("/comments/:id/pin", api.PinComment)
			auth.PUT("/comments/:id", api.PutComment)

			auth.POST("/commentBoard", api.PostMessage)
			auth.DELETE("/commentBoard/:id", api.DeleteMessage)

			auth.POST("/friendsWeb", api.PostFriend)
			auth.PUT("/friendsWeb/:id", api.PutFriend)
			auth.DELETE("/friendsWeb/:id", api.DeleteFriend)
			auth.DELETE("/application/:id", api.DeleteApplication)

			auth.PUT("/about-me", api.PutAboutMe)

			auth.GET("/reminder/:user_id", api.GetUserReminders)
			auth.POST("/reminder", api.PostReminder)
			auth.PUT("/reminder/:id", api.PutReminder)
			auth.DELETE("/reminder/:id", api.DeleteReminder)
			auth.POST("/reminder/:id/email", api.SendReminderEmail)
		}
	}
	return r
}
