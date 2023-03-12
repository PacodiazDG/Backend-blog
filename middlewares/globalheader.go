package middlewares

import (
	"net/http"
	"os"

	"github.com/PacodiazDG/Backend-blog/api/v1/blog"
	"github.com/PacodiazDG/Backend-blog/modules/security"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// Global headers
func GlobalHeader(c *gin.Context) {
	c.Writer.Header().Set("X-Frame-Options", "SAMEORIGIN")
	c.Writer.Header().Set("Referrer-Policy", "same-origin")
	c.Writer.Header().Set("Content-Security-Policy", os.Getenv("ContentSecurityPolicy"))
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

// Middleware checks that the token is valid and searches a blacklist if the token is valid
func NeedAuthentication(c *gin.Context) {
	if c.GetHeader("Authorization") != "" {
		token, err := security.VerifyToken(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Token Not valid"})
			return
		}
		jwtinfo := token.Claims.(jwt.MapClaims)
		if blog.TokenBlackList(jwtinfo["Userid"].(string), jwtinfo["idtoken"].(string)) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Token Forbidden"})
			return
		}
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Status": "Cannot access without token"})
		return
	}
	c.Next()
}
