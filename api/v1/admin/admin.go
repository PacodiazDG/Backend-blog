package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/PacodiazDG/Backend-blog/api/v1/blog"
	"github.com/PacodiazDG/Backend-blog/api/v1/user"
	"github.com/PacodiazDG/Backend-blog/extensions/redisbackend"
	logs "github.com/PacodiazDG/Backend-blog/modules/logs"
	"github.com/PacodiazDG/Backend-blog/modules/security"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// BanUser
func BanUser(c *gin.Context) {
	_, err := security.CheckTokenPermissions([]rune{security.BanUser}, c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": err.Error()})
		return
	}
	IDuser, err := primitive.ObjectIDFromHex(c.Param("UserID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Status": "ID type no valid."})
		return
	}
	Reason := c.DefaultQuery("Reason", "Other")
	result, err := user.UpdateUserInfo(IDuser, &user.UserStrcture{Banned: true})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Status": err.Error()})
		return
	}
	if result.MatchedCount == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Status": "User Not found"})
		return
	}
	err = redisbackend.SetBan(redisbackend.UserRedisJson{
		ID:           IDuser,
		Blocked:      true,
		Reason:       "Banned",
		Date:         time.Now(),
		LoginAttempt: 0,
		Details:      Reason,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"Status": "Banned from db"})
		return
	}
	c.AbortWithStatusJSON(200, gin.H{"Status": "User Baned"})
}

// DelateUser

func DelateUser(c *gin.Context) {
	_, err := security.CheckTokenPermissions([]rune{security.UserManagement}, c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": err.Error()})
		return
	}
	IDuser, err := primitive.ObjectIDFromHex(c.Param("UserID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Status": "ID type no valid."})
		return
	}
	err = user.DelateAccount(IDuser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Status": err.Error()})
		return
	}
	c.AbortWithStatusJSON(200, gin.H{"Status": "User Deleted"})
}

// UnbanUser
func UnbanUser(c *gin.Context) {
	_, err := security.CheckTokenPermissions([]rune{security.BanUser}, c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": err.Error()})
		return
	}
	Info := user.UserStrcture{
		Banned: false,
	}
	IDuser, err := primitive.ObjectIDFromHex(c.Param("UserID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status": "ID type no valid."})
		return
	}
	result, err := user.UpdateUserInfo(IDuser, &Info)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status": "Error reported in logs"})
		return
	}
	if result.MatchedCount == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Status": "User Not found"})
		return
	}
	err = redisbackend.RemoveBan(IDuser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"Status": "Banned from db"})
		return
	}
	c.AbortWithStatusJSON(200, gin.H{"Status": "User Baned"})
}

// ChangeAbout
func ChangeAbout(c *gin.Context) {
	c.AbortWithStatus(200)
}

// ChangeInfoForUser Cambia
func ChangeInfoForUser(c *gin.Context) {
	_, err := security.CheckTokenPermissions([]rune{security.UserManagement}, c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": err.Error()})
		return
	}
	var Info user.UserStrcture
	if err := c.ShouldBindJSON(&Info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": err.Error()})
		return
	}
	objID, err := primitive.ObjectIDFromHex(Info.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Error: No valid UserID"})
		return
	}
	Info.Permissions = ""
	_, err = user.UpdateUserInfo(objID, &Info)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status: ": err.Error()})
		return
	}
}

// UserManagement
func UserManagement(c *gin.Context) {
	_, err := security.CheckTokenPermissions([]rune{security.UserManagement}, c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": err.Error()})
		return
	}
	var result user.UserStrcture
	if err := c.ShouldBindJSON(&result); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(result.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	Info := user.UserStrcture{
		Username:    result.Username,
		Password:    string(hashedPassword),
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Email:       result.Email,
		Image:       result.Image,
		Permissions: result.Permissions,
	}
	IDuser, err := primitive.ObjectIDFromHex(result.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status": "ID type no valid."})
		return
	}
	if user.EmailIsAvailable(result.Email) {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Status": "This email is already being used."})
		return
	}
	_, err = user.UpdateUserInfo(IDuser, &Info)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Status": "User Updated"})
}

// ManualUpdateFeed
func ManualUpdateFeed(c *gin.Context) {
	_, err := security.CheckTokenPermissions([]rune{security.SiteConfig}, c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": err.Error()})
		return
	}
	blog.SetTopFeed()
	blog.SetTopPost()
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"Status": "Cache Updated"})
}

// ListUsers
func ListofUsers(c *gin.Context) {
	_, err := security.CheckTokenPermissions([]rune{security.UserManagement}, c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": err.Error()})
		return
	}
	var next int64 = 0
	skip := c.Query("next")
	if skip != "" {
		next, err = strconv.ParseInt(skip, 10, 64)
		if err != nil || next < 0 {
			c.AbortWithStatus(http.StatusNotAcceptable)
			return
		}
	}
	next = next * 10
	username := c.Query("Username")
	if username != "" {
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Status": "Id not valid format"})
			return
		}
		next = 0
	}
	listOfUsers, err := ListUsers(next, username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, logs.ErrorApi{Status: err.Error()})
		return
	}
	c.AbortWithStatusJSON(200, gin.H{"Status": listOfUsers})
}
