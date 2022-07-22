package mods

import (
	"bytes"
	"errors"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

func convertImgToBytes(m image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, m, nil)
	if err != nil {
		return nil, errors.New("image not encode")
	}
	return buf.Bytes(), nil
}

func ImageController(c *gin.Context) {
	file := filepath.Clean(c.Param("ImageName"))
	dat, err := os.Open("./Serverfiles/blog/" + file)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	rs := c.Query("rs")
	img, _, err := image.Decode(dat)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Status": "image not decoded"})
		return
	}
	intVar, err := strconv.Atoi(rs)
	if err != nil || intVar <= 50 {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"Status": "rs not valid",
		})
		return
	}
	if rs != "" && intVar <= 2400 {
		imagebyte, err := convertImgToBytes(resize.Resize(uint(intVar), 0, img, resize.Lanczos3))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status": err.Error()})
			return
		}
		c.Data(http.StatusOK, "image/png", imagebyte)
		return
	}
	c.Request.Header.Add("x-request-id", "requestID")
	c.File("./Serverfiles/blog/" + file)
}
