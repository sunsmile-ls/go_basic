package routes

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置成发布模式
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	v1 := r.Group("/api/v1")
	// 注册用户
	v1.POST("/signup", controller.SignUpHandle)
	// 登录
	v1.POST("/login", controller.LoginHandle)
	v1.Use(middlewares.JWTAuthMiddleware()) // 应用JWT认证中间件
	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)
		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.GetPostDetailHandler)
		v1.GET("/posts", controller.GetPostListHandler)
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
