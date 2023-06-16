package router

import (
	"net/http"
	"os"

	"github.com/PacodiazDG/Backend-blog/api/v1/admin"
	"github.com/PacodiazDG/Backend-blog/api/v1/blog"
	"github.com/PacodiazDG/Backend-blog/api/v1/sitemap"
	"github.com/PacodiazDG/Backend-blog/api/v1/user"
	Middlewares "github.com/PacodiazDG/Backend-blog/middlewares"
	"github.com/gin-gonic/gin"
)

var blogs = blog.InitControllerPost()

func BackendRouter(router *gin.Engine) {
	blogs.SetCollection("Post")
	if os.Getenv("Storage") == "true" {
		router.Static("/assets/", "./Serverfiles")
	}
	PageManagement(router)
	router.GET("/sitemap.xml", sitemap.SiteMapxml)
	router.GET("/pages", index)
	router.HEAD("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	v1 := router.Group("/v1")
	{
		v1.Use(Middlewares.ApiInfo)
		BlogRouter := v1.Group("/blog")
		{
			BlogAdminRouter := BlogRouter.Group("/auth")
			BlogAdminRouter.Use(Middlewares.NeedAuthentication)
			{
				BlogAdminRouter.POST("/InsetPost", blogs.InsertPost)
				BlogAdminRouter.DELETE("/DelatePost/:ObjectId", blogs.DelatePost)
				BlogAdminRouter.PUT("/UpdatePost/:ObjectId", blogs.UpdatePost)
				BlogAdminRouter.POST("/UploadImage", blogs.UploadImage)
				BlogAdminRouter.GET("/GetMyPost", blogs.MyPosts)
			}
			BlogRouter.GET("/p/:ObjectId", blogs.Post)
			BlogRouter.GET("/feed", blogs.Feed)
			BlogRouter.GET("/find", blogs.FindPost)
			BlogRouter.GET("/visibility/:ObjectId", blogs.Visibility)
			BlogRouter.GET("/RecommendedPost/:ObjectId", blogs.RecommendedPost)
		}

		DraftsRouter := v1.Group("/drafts")
		{
			DraftsRouter.Use(Middlewares.NeedAuthentication)
			Drafts := blog.InitControllerPost()
			Drafts.SetCollection("Drafts")
			DraftsRouter.POST("/InsetPost", Drafts.InsertPost)
			DraftsRouter.DELETE("/DelatePost/:ObjectId", Drafts.DelatePost)
			DraftsRouter.PUT("/UpdatePost/:ObjectId", Drafts.UpdatePost)
			DraftsRouter.GET("/initialize", Drafts.Initialize)
			DraftsRouter.GET("/GetMyDrafts", Drafts.MyPosts)
			DraftsRouter.GET("/p/:ObjectId", Drafts.Post)
			DraftsRouter.GET("/find", Drafts.FindPost)
		}

		MyUser := v1.Group("/user")
		{
			MyUserAhut := MyUser.Group("/Ahut")
			{
				MyUserAhut.Use(Middlewares.NeedAuthentication)
				MyUserAhut.GET("/My", user.UserInfo)
				MyUserAhut.GET("/Iploggeduser", user.Iploggeduser)
				MyUserAhut.GET("/DelateAccount", user.DelateaAccount)
				MyUserAhut.PUT("/My", user.Updateinfo)
				MyUserAhut.GET("/removeToken/:token", user.DelateSession)
				MyUserAhut.GET("/signout", user.Signout)

				MyUserAhut.GET("/CheckToken", func(c *gin.Context) {
					c.AbortWithStatus(http.StatusOK)
				})
				MyUserAhut.GET("/TokenRenewal", user.TokenRenewal)
			}
			MyUser.POST("/login", user.Login)
			MyUser.POST("/RecoveryAccount", user.RecoveryAccount)
			MyUser.GET("/RecoveryAccount/:Token", user.ValidateRecoveryAccount)
		}
		Adminsite := v1.Group("/admin/")
		Adminsite.Use(Middlewares.NeedAuthentication)
		{
			Adminsite.POST("/CreateAccount", user.CreateAccount)
			Adminsite.POST("/BanToken", user.Updateinfo)
			Adminsite.GET("/Ban/:UserID", admin.BanUser)
			Adminsite.GET("/Unban/:UserID", admin.UnbanUser)
			Adminsite.GET("/Cacherefresh", admin.ManualUpdateFeed)
			Adminsite.POST("/UserManagement", admin.UserManagement)
			Adminsite.GET("/GetUsers", admin.ListofUsers)
			Adminsite.GET("/DelateAcount/:UserID", admin.DelateUser)
		}
		FileSystem := v1.Group("/Uplads/")
		FileSystem.Use(Middlewares.NeedAuthentication)
		{
			FileSystem.POST("/Getfile")
		}
	}
}
