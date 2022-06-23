package validation

func TruncateString(text string, lent int) string {
	if len(text) <= lent {
		return text
	}
	return text[:lent]
}
