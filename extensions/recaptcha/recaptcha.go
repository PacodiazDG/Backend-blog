package recaptcha

import (
	"errors"
	"net/http"
	"os"
	"strings"
)

func CheckRecaptcha(ip, response string) (bool, error) {
	url := "https://www.google.com/recaptcha/api/siteverify"
	r := os.Getenv("rSecret")
	if r == "" {
		return false, errors.New("not configured recaptcha")
	}
	payload := strings.NewReader("secret=" + os.Getenv("rSecret") + "&response=" + response + "&remoteip=" + ip + "")
	cdscds, _ := http.NewRequest("POST", url, payload)
	_ = cdscds
	return false, nil
}
