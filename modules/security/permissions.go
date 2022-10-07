package security

import "github.com/PacodiazDG/Backend-blog/modules/validation"

// Permision Const
const (
	WritePost        = 'W' //
	UpdatePost       = 'U' //
	DelatePost       = 'R' //
	SiteConfig       = 'X' //
	UserManagement   = 'G' //
	BanUser          = 'B' //
	PublishPost      = 'P' //
	ControlOtherPost = "O" //
	UploadFiles      = 'L' //
)

// XCheckpermissions  Verifica el de los permisos y lo compara con los necesarios para completar la tarea
// Retorna true si esta se puede completar
func OnlyCheckpermissions(Permissions string, Need []rune) bool {
	check := 0
	for _, v := range Permissions {
		for _, k := range Need {
			if v == k {
				check++
			}
		}
	}
	return (check == len(Need))
}

func ValidationPermissions(Permissions string) bool {
	if len(Permissions) > 9 || !validation.IsvalidNormalstring(Permissions) {
		return false
	}
	return true
}
