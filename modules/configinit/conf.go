package configinit

import (
	"bytes"
	"os"
	"strconv"
	"text/template"

	"github.com/PacodiazDG/Backend-blog/api/v1/blog"
	database "github.com/PacodiazDG/Backend-blog/database"
	services "github.com/PacodiazDG/Backend-blog/services"
	SMTPM "github.com/PacodiazDG/Backend-blog/smtpm"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const banner string = `
================================================================================
=====================+--------------------------------+=========================
=====================|              INIT.            |=========================
=====================+--------------------------------+=========================
================================================================================
`

type test struct {
	Name  string
	Marks int
	Id    string
}

func Conf() {
	println(banner)
	if _, err := os.Stat("./Serverfiles/"); os.IsNotExist(err) {
		FolderErr := os.MkdirAll("./Serverfiles/", os.ModePerm)
		if FolderErr != nil {
			panic(FolderErr)
		}
		FolderErr = os.MkdirAll("./Serverfiles/blog/", os.ModePerm)
		if FolderErr != nil {
			panic(FolderErr)
		}
	}
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	i, err := strconv.Atoi(os.Getenv("TokenExpirationTimeInMinutes"))
	if err != nil {
		color.Red("[Error] Token expiration number is not a valid number ")
		panic(err)
	}
	if i > 120000 {
		color.Yellow("It is not advisable to have such a large number at the expiration of the token")
	}
	if os.Getenv("ginReleaseMode") == "true" {
		gin.SetMode(gin.ReleaseMode)
	}
	// Load template and
	if os.Getenv("TestSMTP") == "true" {
		var tpl bytes.Buffer
		std1 := test{"test", 94, "test"}
		t, err := template.ParseFiles("./Templates/Mail/Test SMTP confg.tmpl")
		if err != nil {
			panic(err)
		}
		if err := t.Execute(&tpl, std1); err != nil {
			panic(err)
		}
		result := tpl.String()
		SMTPM.Send([]string{"example@example.com", "example@example.com"}, "Test", result)
	}
	database.Initdb()
	database.InitRedis()
	blog.SetLastStories()
	services.AutoSetCacheTop()
	services.ImagebackupService()
}
