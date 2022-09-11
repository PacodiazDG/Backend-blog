package Logs

//
const (
	CriticalError = 0
	HardError     = 1
	MediumError   = 2
	LowError      = 3
)

type ErrorApi struct {
	Status string
}

type Error struct {
	ErrorApi error
	Errordbg error
}

func Filter() {

}
