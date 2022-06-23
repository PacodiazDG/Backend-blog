package router

import (
	"net/http"

	admin "github.com/PacodiazDG/Backend-blog/Api/v1/Admin"
	"github.com/PacodiazDG/Backend-blog/Api/v1/Blog"
	mods "github.com/PacodiazDG/Backend-blog/Api/v1/Mods"
	notebook "github.com/PacodiazDG/Backend-blog/Api/v1/Notebook"
	"github.com/PacodiazDG/Backend-blog/Api/v1/Sitemap"
	"github.com/PacodiazDG/Backend-blog/Api/v1/User"
	"github.com/PacodiazDG/Backend-blog/Api/v1/cachehash"
	database "github.com/PacodiazDG/Backend-blog/Database"
	"github.com/PacodiazDG/Backend-blog/Middlewares"
	"github.com/gin-gonic/gin"
)

func test(c *gin.Context) {
	val, err := database.RedisCon.Get("key").Result()
	if err != nil {
		panic(err)
	}
	c.XML(http.StatusOK, gin.H{"message": val, "status": http.StatusOK})
}
func index(c *gin.Context) {
	c.HTML(http.StatusOK, "Index.tmpl", gin.H{
		"title": "Ok",
	})
}

// Routerings
func MercyRouter(router *gin.Engine) {
	router.GET("/sitemap.xml", Sitemap.SiteMapxml)
	router.GET("/", index)
	router.GET("/Image/blog/:ImageName", mods.ImageController)
	//Decrapped <=1
	router.Static("/assets/", "./Serverfiles")

	router.LoadHTMLGlob("Templates/www/*")
	router.HEAD("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	v1 := router.Group("/v1")
	{
		//Blog router
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
			//manualUpdateFeed
		}
		//Drafts

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
		Notes := v1.Group("/notes")
		{
			Notes.Use(Middlewares.NeedAuthentication)

			Notes.POST("/InsetPost", notebook.InsertNotes)
			Notes.GET("/find", notebook.InsertNotes)
			Notes.GET("/GetNotes", notebook.InsertNotes)

		}
		// My user grup
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
		//
		Performance := v1.Group("/Cache")
		{
			// check for changes in the feed otherwise update the browser cache Middlewares.VerifyToken()
			Performance.GET("/FeedHashApi", cachehash.FeedHashApi)
		}

		Adminsite := v1.Group("/admin/")
		Adminsite.Use(Middlewares.NeedAuthentication)
		{
			Adminsite.POST("/CreateAccount", User.CreateAccount)
			Adminsite.POST("/BanToken", User.Updateinfo)     //2/3/22
			Adminsite.GET("/Ban/:UserID", admin.BanUser)     //2/3/22
			Adminsite.GET("/Unban/:UserID", admin.UnbanUser) //2/3/22
			Adminsite.POST("/PrivilegeElevation", test)      //2/3/22
			Adminsite.GET("/Cacherefresh", admin.ManualUpdateFeed)
			Adminsite.POST("/UserManagement", admin.UserManagement)  //2/3/22
			Adminsite.GET("/DumpUsers", admin.ListUsers)             //2/3/22
			Adminsite.GET("/DelateAcount/:UserID", admin.DelateUser) //2/3/22
		}

		//SecAudit
		FileSystem := v1.Group("/Uplads/")
		FileSystem.Use(Middlewares.NeedAuthentication)
		{
			FileSystem.POST("/Getfile")
		}
	}
}
