package validation

func TruncateString(text string, lent int) string {
	if len(text) <= lent {
		return text
	}
	return text[:lent]
}

func SliceStringContains(v []string, find any) bool {
	for _, v2 := range v {
		if v2 == find {
			return true
		}
	}
	return false
}
