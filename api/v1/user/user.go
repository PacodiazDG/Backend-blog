package user

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	logs "github.com/PacodiazDG/Backend-blog/components/logs"
	"github.com/PacodiazDG/Backend-blog/components/security"
	"github.com/PacodiazDG/Backend-blog/components/validation"
	database "github.com/PacodiazDG/Backend-blog/database"
	"github.com/PacodiazDG/Backend-blog/extensions/redisbackend"
	"github.com/PacodiazDG/Backend-blog/smtpm"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	type templateLoginAlert struct {
		IpAddrs   string
		UserAgent string
		Name      string
	}
	var err error
	var result UserStrcture
	var u LoginRequestStruct
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	collection := *database.Database.Collection("Users")
	err = collection.FindOne(context.TODO(), bson.M{"Email": u.Email}).Decode(&result)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Status": "you entered an incorrect username or password "})
		return
	}
	if result.Banned {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "User is banned"})
		return
	}
	uuidtoken, err := uuid.NewRandom()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(u.Password))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Status": "you entered an incorrect username or password "})
		return
	}
	TokenInfo := security.TokenStrocture{
		Email:       result.Email,
		ID:          (result.ID),
		Permissions: result.Permissions,
		Uuid:        uuidtoken,
	}
	objID, _ := primitive.ObjectIDFromHex(result.ID)
	t := time.Now()
	data := IpAddrUser{
		IDuser:    objID,
		IpAddrs:   security.GetIP(c),
		DateOut:   t.Local().Add(time.Hour * time.Duration(168)),
		Date:      t,
		UserAgent: c.Request.Header.Get("user-agent"),
		Uuidtoken: uuidtoken.String(),
	}
	IPAddrUser(&data)
	token, err := security.CreateToken(TokenInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Status": "token creation failed"})
		return
	}
	var tpl bytes.Buffer
	std1 := templateLoginAlert{security.GetIP(c), security.GetUserAgent(c), ""}
	TemplateL, err := template.ParseFiles("./templates/Mail/login.tmpl")
	if err != nil {
		logs.WriteLogs(err, logs.MediumError)
		c.JSON(http.StatusOK, gin.H{"Status": "Internal Server Error"})
		return
	}
	if err = TemplateL.Execute(&tpl, std1); err != nil {
		logs.WriteLogs(err, logs.MediumError)
		c.JSON(http.StatusOK, gin.H{"Status": "Internal Server Error"})
		return
	}
	smtpm.Send([]string{u.Email}, "Your account was accessed from a new IP address", tpl.String())
	c.JSON(http.StatusOK, gin.H{"Token": "Bearer " + token})
}

func CreateAccount(c *gin.Context) {
	var result UserStrcture
	_, err := security.CheckTokenPermissions([]rune{security.UserManagement}, c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&result); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	password := []byte(result.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	result.Password = string(hashedPassword)
	if !validation.IsValidEmail(result.Email) {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Status": "Error", "Descripton:": "Invalid "})
		return
	}
	if !validation.IsvalidNormalstring(result.Username) || !validation.IsvalidNormalstring(result.FirstName) || !validation.IsvalidNormalstring(result.LastName) {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Status": "Invalid Username"})
		return
	}
	collection := *database.Database.Collection("Users")
	var Email bson.M
	collection.FindOne(context.TODO(), bson.M{"Email": result.Email}).Decode(&Email)
	if Email != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Status": "This email is already registered"})
		return
	}
	if len(result.Permissions) < 9 {
		result.Permissions += strings.Repeat("_", 9-len(result.Permissions))
	}
	if !security.ValidationPermissions(result.Permissions) {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Status": "Permission not valid"})
		return
	}
	Data := UserStrcture{
		Email:       result.Email,
		Image:       "https://ui-avatars.com/api/?name=" + url.QueryEscape(result.Username),
		Password:    result.Password,
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Permissions: result.Permissions,
		Username:    result.Username,
		Banned:      false,
		Created_at:  time.Now(),
	}

	collection.InsertOne(context.Background(), Data)
	c.JSON(http.StatusOK, gin.H{"Status": "User created"})
}

func Updateinfo(c *gin.Context) {
	var result UserStrcture
	if err := c.ShouldBindJSON(&result); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	token, err := security.VerifyToken(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Token Not valid"})
		return
	}
	jwtinfo := token.Claims.(jwt.MapClaims)
	if result.Password != "" {
		password := []byte(result.Password)
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		result.Password = string(hashedPassword)
	}
	result = UserStrcture{
		Username:  result.Username,
		Password:  result.Password,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Email:     result.Email,
		Image:     result.Image,
	}
	IDuser, err := primitive.ObjectIDFromHex(jwtinfo["Userid"].(string))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status": "ID type no valid."})
		return
	}
	_, err = UpdateUserInfo(IDuser, &result)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"Status": err.Error(),
		})

		return
	}
	c.JSON(http.StatusOK, gin.H{"Status": "User Updated"})
}

func DelateaAccount(c *gin.Context) {
	if c.Query("Confirm") != "yes" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Status": "it is necessary to confirm the elimination"})
		return
	}
	token, err := security.VerifyToken(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Token not valid"})
		return
	}
	jwtinfo := token.Claims.(jwt.MapClaims)
	IDuser, _ := primitive.ObjectIDFromHex(jwtinfo["Userid"].(string))
	var GetInfo BasicInfo
	err = DelateAccount(IDuser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Status": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Status": GetInfo,
	})

}

func UserInfo(c *gin.Context) {
	token, err := security.VerifyToken(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Token not valid"})
		return
	}
	jwtinfo := token.Claims.(jwt.MapClaims)
	IDuser, _ := primitive.ObjectIDFromHex(jwtinfo["Userid"].(string))
	var GetInfo BasicInfo
	err = GetBasicInfo(IDuser, &GetInfo)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Info": GetInfo,
	})

}

func Iploggeduser(c *gin.Context) {
	token, err := security.VerifyToken(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Token not valid"})
		return
	}
	jwtinfo := token.Claims.(jwt.MapClaims)
	IDuser, _ := primitive.ObjectIDFromHex(jwtinfo["Userid"].(string))
	ipaddres, err := GetIpaddes(IDuser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status": err})
	}
	c.JSON(http.StatusOK, gin.H{
		"Session": ipaddres,
	})

}

func CheckToken(c *gin.Context) {
	_, err := security.VerifyToken(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Token not valid"})
		return
	}
	c.AbortWithStatus(http.StatusOK)
}

func DelateSession(c *gin.Context) {
	uuID := c.Param("token")
	if _, err := primitive.ObjectIDFromHex(uuID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status": "ObjectID not valid"})
		return
	}
	if redisbackend.SetBanToken(uuID, "session deleted by user") != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status": "an error occurred trying to ban the token"})
		return
	}
	if RemoveIPAddrUser(uuID) != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status": "an error occurred trying to ban the token"})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"Status": "uuid removed"})
}

func Signout(c *gin.Context) {
	jwtinfo, err := security.GetinfoToken(security.ExtractToken(c.Request))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Status": "Token not valid"})
		return
	}
	uuID := jwtinfo["idtoken"].(string)
	if redisbackend.SetBanToken(uuID, "session deleted by user") != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status": "an error occurred trying to ban the token"})
		return
	}
	if RemoveIPAddrUser(uuID) != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status": "an error occurred trying to ban the token"})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"Status": "uuid removed"})
}

// Renew the token that is about to expire and delete the token that is about to expire.
func TokenRenewal(c *gin.Context) {
	jwtinfo, err := security.GetinfoToken(security.ExtractToken(c.Request))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Status": "Token not valid"})
		return
	}
	Uuid, err := uuid.Parse(jwtinfo["idtoken"].(string))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	TokenInfo := security.TokenStrocture{
		Email:       jwtinfo["Email"].(string),
		ID:          jwtinfo["Userid"].(string),
		Permissions: jwtinfo["authority"].(string),
		Uuid:        Uuid,
	}
	token, err := security.CreateToken(TokenInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Status": "token creation failed"})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"Token": "Bearer " + token})
}

func SearchMyPublications(c *gin.Context) {
	Next, err := strconv.Atoi(c.Query("next"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Status": "Token not valid"})
		return
	}
	jwtinfo, err := security.GetinfoToken(security.ExtractToken(c.Request))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Status": "Token not valid"})
		return
	}
	Publications, err := GetPublications(jwtinfo["Userid"].(string), Next)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Status": "Token not valid"})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"Post": Publications})
}
