package Blog

import (
	"crypto/sha512"
	"encoding/hex"
	"io"
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

// InsertPost
func (v *PostController) InsertPost(c *gin.Context) {
	jwtinfo, err := Security.CheckTokenPermissions([]rune{Security.PublishPost}, c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": err.Error()})
		return
	}
	var result PostSimpleStruct
	if err := c.ShouldBindJSON(&result); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	if err := IsValidStruct(&result); err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Status": "invalid structure, " + err.Error()})
		return
	}
	if result.Imagen == "" {
		result.Imagen = "https://cdn.pixabay.com/photo/2015/04/23/21/59/tree-736877_960_720.jpg"
	}
	re := regexp.MustCompile(`([A-Fa-f0-9]{128}(.jpg|.jpeg|.png|.gif))`)
	result = PostSimpleStruct{
		Title:         result.Title,
		Content:       result.Content,
		Tags:          result.Tags,
		Date:          time.Now(),
		Author:        jwtinfo["Userid"].(string),
		Visible:       result.Visible,
		Imagen:        result.Imagen,
		Password:      result.Password,
		Description:   validation.TruncateString((result.Description), 250),
		Views:         0,
		UrlImageFound: re.FindAllString(result.Content, -1),
	}
	Status, err := v.Conf.ModelInsertPost(&result)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status:": "Error inserting the post"})
		return
	}
	c.JSON(http.StatusOK, Status)
	ReflexCache()
}

// DelatePost
func (v *PostController) DelatePost(c *gin.Context) {
	_, err := Security.CheckTokenPermissions([]rune{Security.UploadFiles}, c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": err.Error()})
		return
	}
	objectId, err := primitive.ObjectIDFromHex(c.Param("ObjectId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Status:": "Error ObjectId Invalid"})
		return
	}
	err = v.Conf.DelatePost(objectId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotImplemented, gin.H{"Status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Status": "Successfully removed"})
	ReflexCache()
}

// MyPosts Retorna los post publicados por el usuario
func (v *PostController) MyPosts(c *gin.Context) {
	var next int64 = 0
	var err error
	skip := c.Query("next")
	if skip != "" {
		next, err = strconv.ParseInt(skip, 10, 64)
		if err != nil || next < 0 {
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
	Feed, err := v.Conf.GetFeed(next, query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, Feed)
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
		Description: validation.TruncateString((result.Description), 179),
	}
	_, err = v.Conf.ModelUpdate(&result, PostID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, "ok")
	ReflexCache()
}

// Subir imagenes al servidor
func (PostController) FileSystemImage(c *gin.Context) {
	_, err := Security.CheckTokenPermissions([]rune{Security.UploadFiles}, c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": err.Error()})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Status": "No file is received"})
		return
	}
	infofile, err := file.Open()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	byteContainer, err := io.ReadAll(infofile)
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status": "Unable to save the file"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Url": "/assets/blog/" + hash + ".png"})
}

// Initialize Inizializa un post o un draft
func (v *PostController) Initialize(c *gin.Context) {
	jwtinfo, err := Security.CheckTokenPermissions([]rune{Security.PublishPost}, c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": err.Error()})
		return
	}
	Initialize := PostSimpleStruct{
		Title:       "New post",
		Content:     "write your content here",
		Tags:        []string{"Example"},
		Date:        time.Now(),
		Author:      jwtinfo["Userid"].(string),
		Visible:     false,
		Imagen:      "",
		Description: "write your description here",
		Views:       0,
	}
	Infomodel, err := v.Conf.ModelInsertPost(&Initialize)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status:": "An error occurred initializing the post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Status": Infomodel.InsertedID})
}
