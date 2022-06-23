package notebook

import (
	"net/http"
	"time"

	"github.com/PacodiazDG/Backend-blog/Modules/Security"
	"github.com/PacodiazDG/Backend-blog/Modules/validation"
	"github.com/gin-gonic/gin"
)

func AllNotes() {

}

func SeachNotes() {

}

func InsertNotes(c *gin.Context) {
	jwtinfo, err := Security.GetinfoToken(Security.ExtractToken(c.Request))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Error 0x1584584c"})
		return
	}
	var userJson NoteStruct
	if err := c.BindJSON(&userJson); err != nil {
		c.AbortWithStatus(400)
		return
	}
	note := NoteStruct{
		Title:   validation.TruncateString(userJson.Title, 50),
		Content: userJson.Content,
		Date:    time.Now(),
		Tags:    userJson.Tags,
		Author:  (jwtinfo["Userid"].(string)),
	}
	modelInsertNote(&note)
}
