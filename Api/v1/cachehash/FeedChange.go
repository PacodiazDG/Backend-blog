package cachehash

import (
	"net/http"

	"github.com/PacodiazDG/Backend-blog/Modules/Security"
	"github.com/gin-gonic/gin"
)

var FeedHash string

func FeedHashApi(c *gin.Context) {
	if FeedHash == "" {
		Hasran, err := Security.GenerateRandomString(35)
		if err != nil {
			panic(err)
		}
		data := []byte(Hasran)
		FeedHash = string((data))
	}

	c.JSON(http.StatusOK, gin.H{
		"Post": FeedHash})
}

func SetHash() {
	Hasran, err := Security.GenerateRandomString(35)
	if err != nil {
		panic(err)
	}
	data := []byte(Hasran)
	FeedHash = string((data))
}
