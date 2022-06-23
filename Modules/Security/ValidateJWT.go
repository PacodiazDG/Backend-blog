package Security

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

func PermissionsManager() {

}

//ExtractToken
func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

//GetinfoToken Solo obtine la informacion del token
func GetinfoToken(tokenStr string) (jwt.MapClaims, error) {
	token, _ := jwt.Parse(tokenStr, nil)
	return token.Claims.(jwt.MapClaims), nil
}

//VerifyToken Esto verifica que el token sea valido a partir del cifrado de jwt de lo contrario regresa un error
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

	return token, nil
}

// TokenValid Solo valida que el token sea válido y que este vijente
func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
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
