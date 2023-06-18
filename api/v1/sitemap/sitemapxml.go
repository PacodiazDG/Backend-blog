package sitemap

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"strconv"

	database "github.com/PacodiazDG/Backend-blog/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// Generates a map (sitemapindex) with "sitemap" directions
func siteMapLoc(CountDoc int64) string {
	Meta := "<sitemapindex xmlns=\"http://www.google.com/schemas/sitemap/0.84\">\n"
	for i := int64(0); i < CountDoc/10+1; i++ {
		Meta += "<sitemap><loc>" + os.Getenv("Siteurl") + "sitemap.xml?next=" + strconv.FormatInt(i*10, 10) + "</loc></sitemap>"
	}
	Meta += "</sitemapindex>"
	return Meta
}

func SiteMapxml(c *gin.Context) {
	var results []bson.M
	collection := *database.Database.Collection("Post")
	var next int64 = 0
	var err error
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"_id": 1})
	findOptions.SetLimit(10)
	findOptions.SetProjection(bson.M{"_id": 1, "Title": 1})
	skip := c.Query("next")
	c.Header("content-type", "text/xml")
	if skip != "" {
		next, err = strconv.ParseInt(skip, 10, 64)
		if err != nil || next < 0 {
			c.AbortWithStatus(http.StatusNotAcceptable)
			return
		}
	} else {
		CountDoc, err := collection.CountDocuments(c.Request.Context(), bson.M{"Visible": true})
		if err != nil {
			panic(err)
		}
		c.String(http.StatusOK, siteMapLoc(CountDoc))
		return
	}
	findOptions.SetSkip(next)
	cursor, err := collection.Find(context.Background(), bson.M{"Visible": true}, findOptions)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		panic(err)
	}
	if err = cursor.All(context.Background(), &results); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		panic(err)
	}
	siteUrl := os.Getenv("Siteurl")
	if os.Getenv("ExternSiteMap") == "true" {
		siteUrl = os.Getenv("SecondSite")
	}

	Meta := "<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd\">"
	for _, value := range results {

		u, err := url.Parse(siteUrl + "Pages?id=" + value["_id"].(primitive.ObjectID).Hex())
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		Meta += "<url><loc>" + u.String() + "</loc></url>"
	}
	Meta += "</urlset>"
	c.String(http.StatusOK, Meta)
}
