package router

import (
	"net/http"

	admin "github.com/PacodiazDG/Backend-blog/Api/v1/Admin"
	"github.com/PacodiazDG/Backend-blog/Api/v1/Blog"
	"github.com/PacodiazDG/Backend-blog/Api/v1/Fileupload"
	"github.com/PacodiazDG/Backend-blog/Api/v1/Sitemap"
	"github.com/PacodiazDG/Backend-blog/Api/v1/User"
	"github.com/PacodiazDG/Backend-blog/Middlewares"
	"github.com/gin-gonic/gin"
)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "Index.tmpl", gin.H{
		"title": "Ok",
	})
}

func BackendRouter(router *gin.Engine) {
	router.GET("/sitemap.xml", Sitemap.SiteMapxml)
	router.GET("/", index)
	router.GET("/Image/blog/:ImageName", Fileupload.BlogImageUpload)
	router.Static("/assets/", "./Serverfiles")
	router.LoadHTMLGlob("Templates/www/*")
	router.HEAD("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	v1 := router.Group("/v1")
	{
		BlogRouter := v1.Group("/blog")
		{
			Blogs := Blog.InitControllerPost()
			Blogs.SetCollection("Post")
			BlogAdminRouter := BlogRouter.Group("/auth")
			BlogAdminRouter.Use(Middlewares.NeedAuthentication)
			{
				BlogAdminRouter.POST("/InsetPost", Blogs.InsertPost)
				BlogAdminRouter.DELETE("/DelatePost/:ObjectId", Blogs.DelatePost)
				BlogAdminRouter.PUT("/UpdatePost/:ObjectId", Blogs.UpdatePost)
				BlogAdminRouter.POST("/PostFiles", Blogs.FileSystemImage)
				BlogAdminRouter.GET("/GetMyPost", Blogs.MyPosts)
			}
			BlogRouter.GET("/p/:ObjectId", Blogs.Post)
			BlogRouter.GET("/feed", Blogs.Feed)
			BlogRouter.GET("/find", Blogs.FindPost)
			BlogRouter.GET("/visibility/:ObjectId", Blogs.Visibility)
		}

		DraftsRouter := v1.Group("/drafts")
		{
			DraftsRouter.Use(Middlewares.NeedAuthentication)
			Drafts := Blog.InitControllerPost()
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
				MyUserAhut.GET("/My", User.UserInfo)
				MyUserAhut.GET("/Iploggeduser", User.Iploggeduser)
				MyUserAhut.GET("/DelateAccount", User.DelateaAccount)
				MyUserAhut.PUT("/My", User.Updateinfo)
				MyUserAhut.GET("/removeToken/:uudi", User.DelateSession)
				MyUserAhut.GET("/CheckToken", func(c *gin.Context) {
					c.AbortWithStatus(http.StatusOK)
				})
			}
			MyUser.POST("/login", User.Login)
			MyUser.POST("/RecoveryAccount", User.RecoveryAccount)
			MyUser.GET("/RecoveryAccount/:Token", User.ValidateRecoveryAccount)
		}
		Adminsite := v1.Group("/admin/")
		Adminsite.Use(Middlewares.NeedAuthentication)
		{
			Adminsite.POST("/CreateAccount", User.CreateAccount)
			Adminsite.POST("/BanToken", User.Updateinfo)
			Adminsite.GET("/Ban/:UserID", admin.BanUser)
			Adminsite.GET("/Unban/:UserID", admin.UnbanUser)
			Adminsite.GET("/Cacherefresh", admin.ManualUpdateFeed)
			Adminsite.POST("/UserManagement", admin.UserManagement)
			Adminsite.GET("/GetUsers", admin.GetUsers)
			Adminsite.GET("/DelateAcount/:UserID", admin.DelateUser)
		}
		FileSystem := v1.Group("/Uplads/")
		FileSystem.Use(Middlewares.NeedAuthentication)
		{
			FileSystem.POST("/Getfile")
		}
	}
}
