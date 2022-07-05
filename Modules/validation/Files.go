package validation

import "os"

func FileExists(FileName string) bool {
	_, err := os.Stat(FileName)
	return !os.IsNotExist(err)
}
