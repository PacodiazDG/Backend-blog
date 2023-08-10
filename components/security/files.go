package security

const (
	FileTypeGif = "image/gif"
	FileTypePNG = "image/png"
	FileTypeJPG = "image/jpeg"
)

// validated if image by mime
func IsImageMIME(MIME string) bool {
	for _, v := range []string{FileTypeGif, FileTypePNG, FileTypeJPG} {
		if MIME == v {
			return true
		}
	}
	return false
}
