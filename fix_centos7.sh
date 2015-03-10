#!/usr/bin/env bash

if [ $# -gt 0  ]
then
    if [ $1 == "undo" ]
    then
        sed -i -- 's/zmq3/zmq2/g' setup.sh
        sed -i -- 's/zmq3/zmq2/g' ./src/emdr/emdr.go
        mv ./src/sqlserver/transfer.tmp ./src/sqlserver/transfer.go
    fi
else
    sed -i -- 's/zmq2/zmq3/g' setup.sh
    sed -i -- 's/zmq2/zmq3/g' ./src/emdr/emdr.go
    mv ./src/sqlserver/transfer.go ./src/sqlserver/transfer.tmp
fi
