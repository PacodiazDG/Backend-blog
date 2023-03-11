package router

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":       "test",
		"description": "test",
		"site_name":   "test",
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
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":       result.Title,
		"description": result.Description,
		"site_name":   "Blog",
	})
}
func P404(c *gin.Context) {
	c.HTML(http.StatusNotFound, "index.html", gin.H{
		"title":       "Page not found",
		"description": "Page not found",
		"site_name":   "Blog",
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
	router.GET("/404", index)
	router.GET("/Pages", page)
	router.Static("/static/", "./www-data/static")
	router.NoRoute(index)

}
