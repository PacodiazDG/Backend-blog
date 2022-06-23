package validation

import "os"

func IsFileExists(FileName string) bool {
	_, err := os.Stat(FileName)
	return !os.IsNotExist(err)
}
