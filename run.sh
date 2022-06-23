#!/bin/bash

Red="\[\033[0;31m\]"  
if (( $EUID == 0 )); then
   echo "${RED} For security reasons, the use of the root user has been disabled" 
   exit 1
fi

file=./init
if [ -e "$file" ]; then
    rm ./init -f
    go build init.go
    ./init
else 
    go build init.go
    ./init
fi 

