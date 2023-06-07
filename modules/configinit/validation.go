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
		return errors.New("LogErr cannot be empty")
	}
	return nil
}
