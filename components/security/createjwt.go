package security

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// structure of a user token
type TokenStrocture struct {
	Email       string
	ID          string
	Uuid        uuid.UUID
	Permissions string
}

// Create a token with a user's basic data
func CreateToken(TokenInfo TokenStrocture) (string, error) {
	jwtCreate := jwt.MapClaims{}
	jwtCreate["Email"] = TokenInfo.Email
	jwtCreate["Userid"] = TokenInfo.ID
	jwtCreate["idtoken"] = TokenInfo.Uuid.String()
	jwtCreate["authority"] = TokenInfo.Permissions
	jwtCreate["exp"] = time.Now().Add(time.Hour * 168).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, jwtCreate)
	token, err := at.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
