package SMTPM

import "strings"

func cleanCRLFSmtp(text string) string {
	line := strings.TrimSuffix(text, "\n")
	return line
}
