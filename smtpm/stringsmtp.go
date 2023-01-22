package smtpm

import "strings"

// Clears a String from CRLF
func cleanCRLFSmtp(text string) string {
	line := strings.TrimSuffix(text, "\n")
	return line
}
