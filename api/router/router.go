package router

import (
	"net/http"

	"github.com/PacodiazDG/Backend-blog/api/v1/admin"
	"github.com/PacodiazDG/Backend-blog/api/v1/blog"
	"github.com/PacodiazDG/Backend-blog/api/v1/sitemap"
	"github.com/PacodiazDG/Backend-blog/api/v1/user"
	Middlewares "github.com/PacodiazDG/Backend-blog/middlewares"
	"github.com/gin-gonic/gin"
)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "Index.tmpl", gin.H{
		"title": "Ok",
	})
}

func BackendRouter(router *gin.Engine) {
	router.GET("/sitemap.xml", sitemap.SiteMapxml)
	router.GET("/", index)
	// router.GET("/Image/blog/:ImageName", Fileupload.BlogImageUpload)
	router.Static("/assets/", "./Serverfiles")
	router.LoadHTMLGlob("Templates/www/*")
	router.HEAD("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	v1 := router.Group("/v1")
	{
		BlogRouter := v1.Group("/blog")
		{
			Blogs := blog.InitControllerPost()
			Blogs.SetCollection("Post")
			BlogAdminRouter := BlogRouter.Group("/auth")
			BlogAdminRouter.Use(Middlewares.NeedAuthentication)
			{
				BlogAdminRouter.POST("/InsetPost", Blogs.InsertPost)
				BlogAdminRouter.DELETE("/DelatePost/:ObjectId", Blogs.DelatePost)
				BlogAdminRouter.PUT("/UpdatePost/:ObjectId", Blogs.UpdatePost)
				BlogAdminRouter.POST("/UploadImage", Blogs.UploadImage)
				BlogAdminRouter.GET("/GetMyPost", Blogs.MyPosts)
			}
			BlogRouter.GET("/p/:ObjectId", Blogs.Post)
			BlogRouter.GET("/feed", Blogs.Feed)
			BlogRouter.GET("/find", Blogs.FindPost)
			BlogRouter.GET("/visibility/:ObjectId", Blogs.Visibility)
			BlogRouter.GET("/RecommendedPost/:ObjectId", Blogs.RecommendedPost)
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
				MyUserAhut.GET("/removeToken/:uudi", user.DelateSession)
				MyUserAhut.GET("/CheckToken", func(c *gin.Context) {
					c.AbortWithStatus(http.StatusOK)
				})
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
