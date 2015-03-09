#!/bin/sh

if [ -f "eve" ]
then
    rm eve
fi

go build eve.go
chmod 755 eve

if [ $# -eq 0 ]
then
    ./eve -verbose=vvv emdr
else
    ./eve $@
fi
