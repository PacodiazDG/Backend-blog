package security

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PacodiazDG/Backend-blog/extensions/redisbackend"
	"github.com/golang-jwt/jwt"
)

func PermissionsManager() {

}

// ExtractToken
func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// GetinfoToken Solo obtine la informacion del token
func GetinfoToken(tokenStr string) (jwt.MapClaims, error) {
	token, _ := jwt.Parse(tokenStr, nil)
	return token.Claims.(jwt.MapClaims), nil
}

// VerifyToken Esto verifica que el token sea valido a partir Request
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	token, err := jwt.Parse(ExtractToken(r), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	TokenInfo := token.Claims.(jwt.MapClaims)
	if redisbackend.CheckBan(TokenInfo["Userid"].(string), TokenInfo["idtoken"].(string)) {
		return nil, errors.New("token banned")
	}
	return token, nil
}

// TokenValid Solo valida que el token sea válido y que este vijente
func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if token.Claims.Valid() != nil && !token.Valid {
		return err
	}
	return nil
}

// VerifyAuthority verifica que la autoridad corresponda con el token
func VerifyAuthority(Authoritytoken string, AuthoritySys ...rune) bool {
	finded := 0
	for _, v := range Authoritytoken {
		for _, k := range AuthoritySys {
			if v != k {
				finded++
			}
		}
	}
	return finded == len(AuthoritySys)
}

// CheckTokenPermissions a partir de un  Request valida el token y los permisos solicitados
func CheckTokenPermissions(Need []rune, r *http.Request) (jwt.MapClaims, error) {
	jwtinfo, err := GetinfoToken(ExtractToken(r))
	if err != nil {
		return nil, errors.New("token not valid")
	}
	if !OnlyCheckpermissions((jwtinfo["authority"].(string)), Need) {
		return nil, errors.New("need more permissions")
	}
	return jwtinfo, nil
}
