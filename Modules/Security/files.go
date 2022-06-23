package Security

/// MIMEList
const (
	FileTypeGif = "image/gif"
	FileTypePNG = "image/png"
	FileTypeJPG = "image/jpeg"
)

/// IsImageMIME validad si es imagen por el mime
func IsImageMIME(MIME string) bool {
	for _, v := range []string{FileTypeGif, FileTypePNG, FileTypeJPG} {
		if MIME == v {
			return true
		}
	}
	return false
}
