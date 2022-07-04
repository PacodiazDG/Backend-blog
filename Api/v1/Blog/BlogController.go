package Blog

import (
	"crypto/sha512"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/PacodiazDG/Backend-blog/Modules/Security"
	"github.com/PacodiazDG/Backend-blog/Modules/validation"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

type PostController struct {
	Model      PostModel
	Collection string
}

// InitControllerPost
func InitControllerPost() *PostController {
	return &PostController{}
}

//SetCollection Este método cambia de coleccion en la base de datos
func (v *PostController) SetCollection(Collection string) *PostController {
	v.Model.Collection = Collection
	return v
}

// InsertPost
func (v *PostController) InsertPost(c *gin.Context) {

	jwtinfo, err := Security.GetinfoToken(Security.ExtractToken(c.Request))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Error Token not valid"})
		return
	}
	if !Security.XCheckpermissions((jwtinfo["authority"].(string)), []rune{Security.PublishPost}) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Need more permissions"})
		return
	}
	var result PostSimpleStruct
	if err := c.ShouldBindJSON(&result); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	if err := IsValidStruct(&result); err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Status": "Error",
			"Details": err})
		return
	}
	if result.Imagen == "" {
		result.Imagen = "https://cdn.pixabay.com/photo/2015/04/23/21/59/tree-736877_960_720.jpg"
	}
	re := regexp.MustCompile(`([A-Fa-f0-9]{128}(.jpg|.jpeg|.png|.gif))`)
	matches := re.FindAllString(result.Content, -1)
	result = PostSimpleStruct{
		Title:         result.Title,
		Content:       result.Content,
		Tags:          result.Tags,
		Date:          time.Now(),
		Author:        jwtinfo["Userid"].(string),
		Visible:       result.Visible,
		Imagen:        result.Imagen,
		Password:      result.Password,
		Description:   validation.TruncateString((result.Description), 150) + "...",
		Views:         0,
		UrlImageFound: matches,
	}
	Status, err := v.Model.ModelInsertPost(&result)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Status:": "Error",
			"Details": err,
		})
		return
	}
	c.JSON(http.StatusOK, Status)
	ReflexCache()

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
		if err != nil {
			c.AbortWithStatus(http.StatusNotAcceptable)
			return
		}
	}
	TargetQuery := "Title"
	err = Security.TokenValid(c.Request)
	query = bson.M{TargetQuery: bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "gi"}}}
	var visibility bool = true
	if err == nil {
		visibility = false
	}
	query["Visible"] = visibility
	Feed1, err := v.Model.GetFeed(next, query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status": "Internal Server Error"})
		return
	}

	TargetQuery = "Description"
	query = bson.M{TargetQuery: bson.M{"$regex": primitive.Regex{Pattern: ".*" + search + ".*", Options: "gi"}}}
	query["Visible"] = visibility
	Feed2, err := v.Model.GetFeed(next, query)
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
	c.JSON(200, gin.H{
		"Post": FinalFeed,
	})
}

//DelatePost
func (v *PostController) DelatePost(c *gin.Context) {
	jwtinfo, err := Security.GetinfoToken(Security.ExtractToken(c.Request))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Invalid token"})
		return
	}
	if !Security.XCheckpermissions((jwtinfo["authority"].(string)), []rune{Security.DelatePost}) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Need more permissions"})
		return
	}
	objectId, err := primitive.ObjectIDFromHex(c.Param("ObjectId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Status:": "Error ObjectId Invalid"})
		return
	}
	err = v.Model.DelatePost(objectId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotImplemented, gin.H{"Status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Status": "Successfully removed"})
	ReflexCache()
}

// Feed es el Feed principal
func (v *PostController) Feed(c *gin.Context) {
	var next int64 = 0
	var err error
	skip := c.Query("next")
	if skip != "" && skip != "0" || len(FastFeed) == 0 {
		next, err = strconv.ParseInt(skip, 10, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusNotAcceptable)
			return
		}
	} else {
		c.JSON(200, gin.H{
			"Post": FastFeed,
		})
		return

	}
	query := bson.M{}
	err = Security.TokenValid(c.Request)
	if err != nil {
		query = bson.M{"Visible": true, "Password": ""}
	}
	Feed, err := v.Model.GetFeed(next, query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status": "Internal Server Error"})
		return
	}
	c.JSON(200, gin.H{
		"Post": Feed,
	})
}

// MyPosts Retorna los post publicados por el usuario
func (v *PostController) MyPosts(c *gin.Context) {
	var next int64 = 0
	var err error
	skip := c.Query("next")
	if skip != "" {
		next, err = strconv.ParseInt(skip, 10, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusNotAcceptable)
			return
		}
	}
	jwtinfo, err := Security.GetinfoToken(Security.ExtractToken(c.Request))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Token Not valid"})
		return
	}
	query := bson.M{"Author": jwtinfo["Userid"].(string)}
	Feed, err := v.Model.GetFeed(next, query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status": "Internal Server Error"})
		return
	}
	c.JSON(200, Feed)
}

// Post retorna el post solicitado
func (v *PostController) Post(c *gin.Context) {
	CacheAvailable := false
	var Cache PostSimpleStruct
	PostID, err := primitive.ObjectIDFromHex(c.Param("ObjectId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Stauts": "Id not valid"})
		return
	}
	for i := range *CacheRamPost {
		if (*CacheRamPost)[i].ID == c.Param("ObjectId") {
			Cache = (*CacheRamPost)[i]
			CacheAvailable = true
			break
		}

	}
	if CacheAvailable {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"Post":        Cache,
			"Performance": true,
		})
		go v.Model.Addviews(PostID)
		return
	}
	result, err := (v.Model).ModelGetArticle(PostID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Stauts": err.Error()})
		return
	}
	_, err = Security.VerifyToken((c.Request))
	if result.Password != "" && result.Password != c.Query("Hash") && err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	go v.Model.Addviews(PostID)
	c.JSON(http.StatusOK, gin.H{
		"Post":        result,
		"Performance": false,
	})
}

// UpdatePost
func (v *PostController) UpdatePost(c *gin.Context) {
	jwtinfo, err := Security.GetinfoToken(Security.ExtractToken(c.Request))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Token not valid."})
		return
	}
	var result PostSimpleStruct
	if err := c.ShouldBindJSON(&result); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	PostID, err := primitive.ObjectIDFromHex(c.Param("ObjectId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "Invalid ObjectId")
		return
	}
	result = PostSimpleStruct{
		Title:       result.Title,
		Content:     result.Content,
		Tags:        result.Tags,
		Author:      jwtinfo["Userid"].(string),
		Visible:     result.Visible,
		Imagen:      result.Imagen,
		Password:    result.Password,
		Description: validation.TruncateString((result.Description), 139),
	}
	_, err = v.Model.ModelUpdate(&result, PostID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, "ok")
	ReflexCache()
}

// Initialize Inizializa un post o un draft
func (v *PostController) Initialize(c *gin.Context) {
	jwtinfo, err := Security.GetinfoToken(Security.ExtractToken(c.Request))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"": ""})
		return
	}
	Initialize := PostSimpleStruct{
		Title:       "New post" + time.Now().GoString(),
		Content:     "write your content here",
		Tags:        []string{"Example"},
		Date:        time.Now(),
		Author:      jwtinfo["Userid"].(string),
		Visible:     false,
		Imagen:      "",
		Description: "write your description here",
		Views:       0,
	}
	Infomodel, err := v.Model.ModelInsertPost(&Initialize)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Status:": "Error",
			"Details": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status": Infomodel.InsertedID,
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
	_, err = v.Model.ModelUpdate(&ProcessData, PostID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"": ""})
		return
	}
	c.JSON(http.StatusOK, "ok")
}

// Retorna los post mas vistos
func (v *PostController) SetTop() ([]PostSimpleStruct, error) {
	return v.Model.GetTOP()
}

// Obtine el feed de las ultimas publicaciones
func (v *PostController) FeedFast() ([]FeedStrcture, error) {
	return v.Model.GetFeed(0, bson.M{"Visible": true, "Password": ""})
}

// Subir imagenes al servidor
func (PostController) FileSystemImage(c *gin.Context) {
	jwtinfo, err := Security.GetinfoToken(Security.ExtractToken(c.Request))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Token not valid."})
		return
	}
	if !Security.XCheckpermissions((jwtinfo["authority"].(string)), []rune{Security.UploadFiles}) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Need more permissions"})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}
	infofile, err := file.Open()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	byteContainer, err := ioutil.ReadAll(infofile)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Error reading file"})
		return
	}
	MIME := http.DetectContentType(byteContainer)

	if !Security.IsImageMIME(MIME) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "It is not image"})
		return
	}
	h := sha512.New()
	h.Write(byteContainer)
	hash := hex.EncodeToString(h.Sum(nil))
	if err := c.SaveUploadedFile(file, "./Serverfiles/blog/"+hash+".png"); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Url": "/assets/blog/" + hash + ".png"})
}
