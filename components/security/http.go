package security

import (
	"html"
	"os"

	"github.com/gin-gonic/gin"
)

// Get the ip depending on the configuration
// if the IpaddressByHeader configuration parameter has a value, it will be
// used to obtain the ip of that header
func GetIP(c *gin.Context) string {
	if os.Getenv("UseProxy") == "false" {
		IpAddrs := c.ClientIP()
		return IpAddrs
	}
	return html.EscapeString(c.Request.Header.Get(os.Getenv("IpaddressByHeader")))
}

// Get the user-agent and return it in EscapeString
func GetUserAgent(c *gin.Context) string {
	return html.EscapeString(c.Request.Header.Get("user-agent"))
}
