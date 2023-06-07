package configinit

import (
	"errors"
	"os"
)

func validation() error {
	if os.Getenv("DB_CONFIG") == "" {
		return errors.New("DB_CONFIG cannot be empty")
	}
	if os.Getenv("JWT_SECRET") == "" {
		return errors.New("JWT_SECRET cannot be empty")
	}
	if os.Getenv("TokenExpirationTimeInMinutes") == "" {
		return errors.New("TokenExpirationTimeInMinutes cannot be empty")
	}
	if os.Getenv("DefaultDatabase") == "" {
		return errors.New("DefaultDatabase cannot be empty")
	}
	if os.Getenv("LogErr") == "" {
		return errors.New("LogErr cannot be empty")
	}
	return nil
}
