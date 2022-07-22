package User

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"os"
	"time"

	database "github.com/PacodiazDG/Backend-blog/Database"
	"github.com/PacodiazDG/Backend-blog/Modules/Security"
	"github.com/PacodiazDG/Backend-blog/Modules/validation"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func salvetoken(Token, Email string) error {
	date := time.Now()
	dbConfig := *database.Database

	collection := dbConfig.Collection("RecoveryAccount")
	_, err := collection.InsertOne(context.TODO(), bson.M{"Token": Token, "Used": false, "Date": date.String(), "Email": Email})
	if err != nil {
		return err
	}
	return nil
}

//RecoveryAccount
func RecoveryAccount(c *gin.Context) {
	Email := c.PostForm("EmailRecovery")
	if Email == "" || !validation.IsValidEmail(Email) {
		c.AbortWithStatus(http.StatusNotAcceptable)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, "invalid parameter")

		return
	}
	dbConfig := *database.Database
	collection := dbConfig.Collection("Users")
	var result UserStrcture
	NoFound := collection.FindOne(context.TODO(), bson.M{"Email": Email}).Decode(&result)
	if NoFound != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, "User not Found")
		return
	}

	token, err := Security.GenerateRandomString(50)
	if err != nil {
		panic(err)
	}
	err = salvetoken(token, Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Error"})
		return
	}
	URL := os.Getenv("Siteurl") + "/?token=" + token
	var tpl bytes.Buffer
	std1 := RecoveryAccountStrcture{"test", URL}
	t, err := template.ParseFiles("./Templates/Mail/recoveryaccount.tmpl")
	if err != nil {
		panic(err)
	}
	if err := t.Execute(&tpl, std1); err != nil {
		panic(err)
	}
	results := tpl.String()
	//SMTPM.Send([]string{Email}, "Test", result)
	_ = results
}
