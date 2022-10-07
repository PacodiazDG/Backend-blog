package recaptcha

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

//En desarrollo
func CheckRecaptcha(ip, response string) (bool, error) {
	url := "https://www.google.com/recaptcha/api/siteverify"
	r := os.Getenv("rSecret")
	if r == "" {
		return false, errors.New("not configured recaptcha")
	}
	payload := strings.NewReader("secret=" + os.Getenv("rSecret") + "&response=" + response + "&remoteip=" + ip + "")
	req, _ := http.NewRequest("POST", url, payload)
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(req)
	fmt.Println(string(body))
	return false, nil
}
