package Blog

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/PacodiazDG/Backend-blog/Modules/Security"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

type PostController struct {
	Conf       Queryconf
	Collection string
}

// InitControllerPost
func InitControllerPost() *PostController {
	return &PostController{}
}

// SetCollection Este método cambia de coleccion en la base de datos
func (v *PostController) SetCollection(Collection string) *PostController {
	v.Conf.Collection = Collection
	return v
}

// Obtine el feed de las ultimas publicaciones
func (v *PostController) FeedFast() ([]FeedStrcture, error) {
	return v.Conf.GetFeed(0, bson.M{"Visible": true, "Password": ""})
}

// FindPost Api
func (v *PostController) FindPost(c *gin.Context) {
	var next int64 = 0
	var err error
	var query bson.M
	search := regexp.QuoteMeta(c.Query("q"))
	if len(search) > 800 {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	skip := c.Query("next")
	if skip != "" {
		next, err = strconv.ParseInt(skip, 10, 64)
		if err != nil || next < 0 {
			c.AbortWithStatus(http.StatusNotAcceptable)
			return
		}
	}
	err = Security.TokenValid(c.Request)
	var visibility = true
	if err == nil {
		visibility = false
	}
	query = bson.M{"Title": bson.M{"$regex": primitive.Regex{
		Pattern: ".*" + search + ".*", Options: "gi"}},
		"Visible": visibility, "Password": ""}
	Feed1, err := v.Conf.GetFeed(next, query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status": "Internal Server Error"})
		return
	}
	query = bson.M{"Description": bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "gi"}},
		"Visible": visibility, "Password": ""}
	Feed2, err := v.Conf.GetFeed(next, query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status": "Internal Server Error"})
		return
	}
	Feed1 = append(Feed1, Feed2...)
	flag := make(map[string]bool)
	var FinalFeed []FeedStrcture
	for _, name := range Feed1 {
		if !flag[name.ID] {
			flag[name.ID] = true
			FinalFeed = append(FinalFeed, name)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"Post": FinalFeed,
	})
}

// Feed es el Feed principal
func (v *PostController) Feed(c *gin.Context) {
	var next int64 = 0
	var err error
	skip := c.Query("next")
	if skip != "" && skip != "0" || len(FastFeed) == 0 {
		next, err = strconv.ParseInt(skip, 10, 64)
		if err != nil || next < 0 {
			c.AbortWithStatus(http.StatusNotAcceptable)
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"Post": FastFeed})
		return
	}
	query := bson.M{}
	err = Security.TokenValid(c.Request)
	if err != nil {
		query = bson.M{"Visible": true, "Password": ""}
	}
	Feed, err := v.Conf.GetFeed(next, query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Post": Feed,
	})
}

// Post retorna el post solicitado
func (v *PostController) Post(c *gin.Context) {
	var Cache PostSimpleStruct
	PostID, err := primitive.ObjectIDFromHex(c.Param("ObjectId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Stauts": "Id not valid"})
		return
	}
	for i := range *CacheRamPost {
		if (*CacheRamPost)[i].ID == c.Param("ObjectId") {
			Cache = (*CacheRamPost)[i]
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"Post":        Cache,
				"Performance": true,
			})
			go v.Conf.Addviews(PostID)
			return
		}
	}
	result, err := (v.Conf).ModelGetArticle(PostID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Stauts": err.Error()})
		return
	}
	_, err = Security.VerifyToken((c.Request))
	if result.Password != "" && result.Password != c.Query("Hash") && err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	go v.Conf.Addviews(PostID)
	c.JSON(http.StatusOK, gin.H{
		"Post":        result,
		"Performance": false,
	})
}

// Visibility Cambia la visiblidad de un  post
func (v *PostController) Visibility(c *gin.Context) {
	_, err := Security.GetinfoToken(Security.ExtractToken(c.Request))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Token not valid."})
		return
	}
	VisibleStatus, err := strconv.ParseBool(c.Query("visible"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}
	PostID, err := primitive.ObjectIDFromHex(c.Param("ObjectId"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}
	ProcessData := PostSimpleStruct{
		Visible: VisibleStatus,
	}
	_, err = v.Conf.ModelUpdate(&ProcessData, PostID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"": ""})
		return
	}
	c.JSON(http.StatusOK, "ok")
}

// Retorna los post mas vistos
func (v *PostController) SetTop() ([]PostSimpleStruct, error) {
	return v.Conf.GetTOP()
}
