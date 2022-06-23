package Middlewares

import (
	"net/http"
	"os"

	"github.com/PacodiazDG/Backend-blog/Api/v1/Blog"
	"github.com/PacodiazDG/Backend-blog/Modules/Security"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func GlobalHeader(c *gin.Context) {
	c.Writer.Header().Set("X-Frame-Options", "SAMEORIGIN")
	c.Writer.Header().Set("Referrer-Policy", "same-origin")
	c.Writer.Header().Set("Content-Security-Policy", "default-src 'self' 'unsafe-inline'")
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("Cross"))
	if c.Request.Method == http.MethodOptions {
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS,PUT,DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization,Content-Type")
		c.AbortWithStatus(http.StatusOK)
		return
	}
	c.Next()
}

func NeedAuthentication(c *gin.Context) {
	if c.GetHeader("Authorization") != "" {
		token, err := Security.VerifyToken(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Token Not valid"})
			return
		}
		jwtinfo := token.Claims.(jwt.MapClaims)
		if Blog.TokenBlackList(jwtinfo["Userid"].(string), jwtinfo["idtoken"].(string)) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Token Forbidden"})
			return
		}
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Status": "Cannot access without token"})
		return
	}
	c.Next()
}
