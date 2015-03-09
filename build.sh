#!/usr/bin/env bash

if [ -f "eve" ]
then
    rm eve
fi

go build eve.go
chmod 755 eve

if [ $# -gt 0 ]
then
 ./run.sh $@
fi