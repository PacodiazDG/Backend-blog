package main

import (
	"os"
	"strconv"

	"github.com/PacodiazDG/Backend-blog/api/router"
	Middlewares "github.com/PacodiazDG/Backend-blog/middlewares"
	"github.com/PacodiazDG/Backend-blog/modules/configinit"
	"github.com/PacodiazDG/Backend-blog/modules/validation"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

func main() {
	configinit.Conf()
	lisent := os.Getenv("HostandPort")
	if lisent == "" {
		lisent = ":8080"
	}
	PemFile := os.Getenv("PemFile")
	KeyFile := os.Getenv("KeyFile")
	Server := gin.New()
	Server.Use(Middlewares.GlobalHeader)
	router.BackendRouter(Server)
	println(lisent)
	if validation.FileExists(PemFile) && validation.FileExists(KeyFile) {
		Server.RunTLS(lisent, PemFile, KeyFile)
	} else {
		color.Red("[Error] Files Not found  \n Status of PemFile: " + strconv.FormatBool(validation.FileExists(PemFile)) + "\n Status of KeyFile: " + strconv.FormatBool(validation.FileExists(KeyFile)))
		color.Yellow("[Warning] The server is working without SSL")
		Server.Run(lisent)
	}

}
