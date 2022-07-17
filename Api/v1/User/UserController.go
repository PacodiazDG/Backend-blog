package User

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	database "github.com/PacodiazDG/Backend-blog/Database"
	"github.com/PacodiazDG/Backend-blog/Extensions/RedisBackend"
	"github.com/PacodiazDG/Backend-blog/Modules/Logs"
	"github.com/PacodiazDG/Backend-blog/Modules/Security"
	"github.com/PacodiazDG/Backend-blog/Modules/validation"
	"github.com/PacodiazDG/Backend-blog/SMTPM"
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
		c.AbortWithStatusJSON(http.StatusNotFound, "User not Found")
		return
	}
	if result.Banned {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Details": "User is banned"})
		return
	}
	uuidtoken, err := uuid.NewRandom()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(u.Password))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, "Password not Found")
		return
	}
	TokenInfo := Security.TokenStrocture{
		Email:       result.Email,
		ID:          (result.ID),
		Permissions: result.Permissions,
		Uuid:        uuidtoken,
	}
	objID, _ := primitive.ObjectIDFromHex(result.ID)
	t := time.Now()
	data := IpAddrUser{
		IDuser:    objID,
		IpAddrs:   Security.GetIP(c),
		DateOut:   t.Local().Add(time.Hour * time.Duration(168)),
		Date:      t,
		UserAgent: c.Request.Header.Get("user-agent"),
		Uuidtoken: uuidtoken.String(),
	}
	IPAddrUser(&data)
	token, err := Security.CreateToken(TokenInfo)
	if err != nil {
		Logs.WriteLogs(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var tpl bytes.Buffer
	std1 := templateLoginAlert{Security.GetIP(c), Security.GetUserAgent(c), ""}
	TemplateL, err := template.ParseFiles("./Templates/Mail/Login.tmpl")
	if err != nil {
		Logs.WriteLogs(err)
		c.JSON(http.StatusOK, gin.H{"Token": "Bearer " + token})
		return
	}
	if err = TemplateL.Execute(&tpl, std1); err != nil {
		Logs.WriteLogs(err)
		c.JSON(http.StatusOK, gin.H{"Token": "Bearer " + token})
		return
	}
	SMTPM.Send([]string{u.Email}, "Your account was accessed from a new IP address", tpl.String())
	//Security.BanedToken(token)
	c.JSON(http.StatusOK, gin.H{"Token": "Bearer " + token})
}

func CreateAccount(c *gin.Context) {
	var result UserStrcture
	jwtinfo, err := Security.GetinfoToken(Security.ExtractToken(c.Request))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Error 0x1584584c"})
		return
	}
	if !Security.XCheckpermissions((jwtinfo["authority"].(string)), []rune{Security.UserManagement}) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Need more permissions"})
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
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Status": "Error",
			"Description": "Invalid Username"})
		return
	}
	collection := *database.Database.Collection("Users")
	var Email bson.M
	collection.FindOne(context.TODO(), bson.M{"Email": result.Email}).Decode(&Email)
	if Email != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Status": "Error",
			"Description": "This email is already registered"})
		return
	}
	if len(result.Permissions) < 9 {
		result.Permissions += strings.Repeat("_", 9-len(result.Permissions))
	}
	if !Security.ValidationPermissions(result.Permissions) {
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
	token, err := Security.VerifyToken(c.Request)
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
	//
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
	token, err := Security.VerifyToken(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Details": "Token not valid"})
		return
	}
	jwtinfo := token.Claims.(jwt.MapClaims)
	IDuser, _ := primitive.ObjectIDFromHex(jwtinfo["Userid"].(string))
	var GetInfo BasicInfo
	err = DelateAccount(IDuser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Error code": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Info": GetInfo,
	})

}

func UserInfo(c *gin.Context) {
	token, err := Security.VerifyToken(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Details": "Token not valid"})
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
	token, err := Security.VerifyToken(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Details": "Token not valid"})
		return
	}
	jwtinfo := token.Claims.(jwt.MapClaims)
	IDuser, _ := primitive.ObjectIDFromHex(jwtinfo["Userid"].(string))
	ipaddres, err := GetIpaddes(IDuser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err})
	}
	c.JSON(http.StatusOK, gin.H{
		"Session": ipaddres,
	})

}

func CheckToken(c *gin.Context) {
	_, err := Security.VerifyToken(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Details": "Token not valid"})
		return
	}
	c.AbortWithStatus(http.StatusOK)
}

func DelateSession(c *gin.Context) {
	uudi := c.Param("uudi")
	_, err := uuid.Parse(uudi)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"Status": "uuid is not valid"})
		return
	}
	if RedisBackend.InsertBanidtoken(uudi) != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Details": "an error occurred trying to ban the token"})
		return
	}
	if RemoveIPAddrUser(uudi) != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Details": "an error occurred trying to ban the token"})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"Status": "uuid removed"})
}
