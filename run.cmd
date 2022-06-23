
 IF EXIST "init.exe" (
  del init.exe
  echo "Removing init.exe"
) 
 
go build init.go
init.exe
