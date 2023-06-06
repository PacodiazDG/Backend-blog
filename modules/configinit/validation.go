package configinit

import (
	"errors"
	"os"
)

func validation() error {
	if os.Getenv("DB_CONFIG") == "" {
		return errors.New("DB_CONFIG ")
	}
	if os.Getenv("JWT_SECRET") == "" {
		return errors.New("")
	}
	if os.Getenv("TokenExpirationTimeInMinutes") == "" {
		return errors.New("")
	}
	if os.Getenv("DefaultDatabase") == "" {
		return errors.New("")
	}
	if os.Getenv("LogErr") == "" {
		return errors.New("")
	}
	// SMTP credencials
	if os.Getenv("SMTPHost") == "" {
		return errors.New("")
	}
	if os.Getenv("SMTPport") == "" {
		return errors.New("")
	}
	if os.Getenv("SMTPPassword") == "" {
		return errors.New("")
	}
	if os.Getenv("TestSMTP") == "" {
		return errors.New("")
	}
	if os.Getenv("TestEmailSMTP") == "" {
		return errors.New("")
	}
	if os.Getenv("UseSMTP") == "" {
		return errors.New("")
	}
	//Redis credencials
	if os.Getenv("RedisAddr") == "" {
		return errors.New("")
	}
	if os.Getenv("RedisPassword") == "" {
		return errors.New("")
	}
	//UseProxy
	if os.Getenv("UseProxy") == "" {
		return errors.New("")
	}
	if os.Getenv("IpaddressByHeader") == "" {
		return errors.New("")
	}
	if os.Getenv("") == "" {
		return errors.New("")
	}
	if os.Getenv("") == "" {
		return errors.New("")
	}
	if os.Getenv("") == "" {
		return errors.New("")
	}
	if os.Getenv("") == "" {
		return errors.New("")
	}
	if os.Getenv("") == "" {
		return errors.New("")
	}
	if os.Getenv("") == "" {
		return errors.New("")
	}
	if os.Getenv("") == "" {
		return errors.New("")
	}

	return nil
}
