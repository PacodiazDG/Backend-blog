package validation

import "os"

// Check if the file exists
func FileExists(FileName string) bool {
	_, err := os.Stat(FileName)
	return !os.IsNotExist(err)
}
