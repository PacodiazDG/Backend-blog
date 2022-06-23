package SMTPM

import (
	"encoding/base64"
	"net/smtp"
	"os"

	"github.com/PacodiazDG/Backend-blog/Modules/validation"
)

//ValidadteEmailArray verifica
func ValidadteEmailArray(email []string) bool {
	for _, value := range email {
		println(value)
		if !validation.IsValidEmail(value) {
			return false
		}
	}
	return true
}

//Send envia el correo uno por uno a partir de un array
func Send(to []string, Subject, text string) {
	for _, s := range to {
		if validation.IsValidEmail(s) {
			SendEmail([]string{s}, Subject, text)
		}
	}
}

//SendEmail Envia mail es el framework por default de mercy
func SendEmail(to []string, Subject, text string) {
	if os.Getenv("UseSMTP") == "false" {
		return
	}
	str := base64.StdEncoding.EncodeToString([]byte(text))
	SMTPHost := os.Getenv("SMTPHost")
	SMTPport := os.Getenv("SMTPport")
	SMTPUsername := os.Getenv("SMTPUsername")
	SMTPPassword := os.Getenv("SMTPPassword")
	auth := smtp.PlainAuth("", SMTPUsername, SMTPPassword, SMTPHost)
	msgs := []byte("MIME-Version: 1.0\n" +
		"To: " + to[0] + "\n" +
		"Subject: " + cleanCRLFSmtp(Subject) + "\n" +
		"Content-Type: text/html; charset=utf-8\n" +
		"Content-Transfer-Encoding: base64\n" +
		"\n\n" +
		str +
		"\n")
	err := smtp.SendMail(SMTPHost+":"+SMTPport, auth, SMTPUsername, to, msgs)
	if err != nil {
		panic(err)
	}
}
