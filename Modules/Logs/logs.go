package Logs

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

//WriteLogs maneja los errores
func WriteLogs(errors error) {
	s := errors.Error()
	f, err := os.OpenFile("./error.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		makefile("./error.log")
		panic(err)
	}
	defer f.Close()
	currentTime := time.Now()
	Log := currentTime.Format("2006.01.02 15:04:05") + " " + s + "\n"
	if _, err := f.WriteString(Log); err != nil {
		panic(err)
	}
}
func makefile(name string) {
	err := ioutil.WriteFile(name, nil, 0644)
	if err != nil {
		fmt.Printf("Unable to write file: %v", err)
	}
}
