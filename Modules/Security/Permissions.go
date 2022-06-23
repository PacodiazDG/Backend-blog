package Security

import "github.com/PacodiazDG/Backend-blog/Modules/validation"

/*
EDITOR|

*/
//Permision Const
const (
	WritePost        = 'W' //FF
	UpdatePost       = 'U' //vfdvdfF
	DelatePost       = 'R' //
	SiteConfig       = 'X' //
	UserManagement   = 'G' //
	BanUser          = 'B' //
	PublishPost      = 'P' //
	ControlOtherPost = "O" //
	UploadFiles      = 'L' //
)

//XCheckpermissions Verifica el de los permisisos y lo compara con los necesarios para completar la tarea
//.Retorna true si esta se puede completar
func XCheckpermissions(Permissions string, Need []rune) bool {
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
	println(!validation.IsvalidNormalstring(Permissions))
	if len(Permissions) > 9 || !validation.IsvalidNormalstring(Permissions) {
		return false
	}
	return true
}
