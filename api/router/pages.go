package router

import (
	"html"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":         html.EscapeString(os.Getenv("SiteMetaTitle")),
		"description":   html.EscapeString(os.Getenv("SiteMetaDescription")),
		"site_name":     html.EscapeString(os.Getenv("SiteMetaTitle")),
		"gverification": html.EscapeString(os.Getenv("GoogleSite_Verification")),
	})
}
func page(c *gin.Context) {
	PostID, err := primitive.ObjectIDFromHex(c.Query("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Stauts": "Id not valid"})
		return
	}
	result, err := blogs.Conf.GetMetaPost(PostID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Stauts": err.Error()})
		return
	}
	if !result.Visible {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":         "",
			"description":   "",
			"site_name":     "",
			"ogimage":       "",
			"gverification": "",
		})
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":         html.EscapeString(result.Title),
		"description":   html.EscapeString(result.Description),
		"site_name":     html.EscapeString(os.Getenv("SiteMetaTitle")),
		"ogimage":       html.EscapeString(result.Imagen),
		"gverification": html.EscapeString(os.Getenv("GoogleSite_Verification")),
	})
}
func P404(c *gin.Context) {
	c.HTML(http.StatusNotFound, "index.html", gin.H{
		"title":       "Page not found",
		"description": "Page not found",
		"site_name":   html.EscapeString(os.Getenv("SiteMetaTitle")),
	})
}
func PageManagement(router *gin.Engine) {

	if os.Getenv("Pages") != "true" {
		router.GET("/", func(ctx *gin.Context) {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"Status": "Ok"})
		})
		return
	}
	router.LoadHTMLFiles("./www-data/index.html")
	router.GET("/404", P404)
	router.GET("/Pages", page)
	router.Static("/static/", "./www-data/static")
	router.NoRoute(index)

}
