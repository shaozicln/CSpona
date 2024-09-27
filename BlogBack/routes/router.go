package routes

import (
	"BlogBack/api"
	"BlogBack/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 函数名开头大写：公用方法；小写：私有方法，只有当前包内可以应用
func InitRouter() *gin.Engine {
	gin.SetMode(utils.AppMode)
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	r.Use(Cors(config))

	router := r.Group("api")
	{
		//...登陆验证路由... Login.go
		router.POST("/login", api.Login)

		//...用户路由... User.go
		router.GET("/users/:id", api.GetUserWithId)
		router.GET("/users", api.GetUserWithUsername)
		router.POST("/users", api.PostUser)
		router.PUT("/users/:username", api.PutUser)
		router.DELETE("/users/:id", api.DeleteUser)

		// ...文章路由... Article.go
		router.GET("/search", api.SearchArticle)
		router.POST("/articles", api.PostArticle)
		router.PUT("/articles/:id", api.PutArticle)
		router.DELETE("/articles/:id", api.DeleteArticle)

		// ...分类路由... Category.go
		router.GET("/categories", api.GetCategory)
		router.POST("/categories", api.PostCatehgory)
		router.PUT("/categories/:id", api.PutCategory)
		router.DELETE("/categories/:id", api.DeleteCategory)

		//...主页文章展示路由... CategoryWithArticle.go
		router.GET("/categories-with-articles", api.CategoryWithArticle)

		//...文章内容展示+评论路由... ContentAndComment.go
		router.GET("/categories", api.GetContent)
		router.POST("/categories", api.PostContent)
		router.DELETE("/categories/:id", api.DeleteContent)
		router.GET("/categories/:id", api.GetComment)

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
	}

	r.SetTrustedProxies([]string{utils.DbHost})
	r.Run(utils.HttpPort)
	return r
}
