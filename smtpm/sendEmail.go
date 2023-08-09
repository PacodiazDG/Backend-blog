package smtpm

import (
	"encoding/base64"
	"net/smtp"
	"os"

	"github.com/PacodiazDG/Backend-blog/components/logs"
	"github.com/PacodiazDG/Backend-blog/components/validation"
)

// Verify that the mailing list is valid.
func ValidadteEmailArray(email []string) bool {
	for _, value := range email {
		if !validation.IsValidEmail(value) {
			return false
		}
	}
	return true
}

// send mail one by one
func Send(to []string, Subject, text string) {
	for _, s := range to {
		if validation.IsValidEmail(s) {
			SendEmail([]string{s}, Subject, text)
		}
	}
}

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
		logs.WriteLogs(err, logs.MediumError)
		panic(err)
	}
}
