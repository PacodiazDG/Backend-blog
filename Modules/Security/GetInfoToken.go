package Security

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func GetInfoTokenbyHeader(c *gin.Context) (jwt.MapClaims, error) {
	token, err := VerifyToken(c.Request)
	if err != nil {
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil
}
