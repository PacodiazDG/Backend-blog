package main

import (
	"os"
	"strconv"

	"github.com/PacodiazDG/Backend-blog/Api/router"
	"github.com/PacodiazDG/Backend-blog/Middlewares"
	"github.com/PacodiazDG/Backend-blog/Modules/ConfigInit"
	"github.com/PacodiazDG/Backend-blog/Modules/validation"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

func main() {
	lisent := os.Getenv("PORT")
	ConfigInit.Conf()
	if lisent == "" {
		lisent = os.Getenv("Port")
	}
	PemFile := os.Getenv("PemFile")
	KeyFile := os.Getenv("KeyFile")
	Server := gin.Default()
	Server.Use(Middlewares.GlobalHeader)
	router.MercyRouter(Server)
	if lisent == "" {
		lisent = ":8080"
	}
	if validation.IsFileExists(PemFile) && validation.IsFileExists(KeyFile) {
		Server.RunTLS(lisent, PemFile, KeyFile)
	} else {
		color.Red("[Error] Files Not found  \n Status of PemFile: " + strconv.FormatBool(validation.IsFileExists(PemFile)) + "\n Status of KeyFile: " + strconv.FormatBool(validation.IsFileExists(KeyFile)))
		color.Yellow("[Warning] The server is working without SSL")
		Server.Run(lisent)
	}

}
