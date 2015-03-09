#!/bin/sh

if [ $# -eq 0 ]
then
    ./eve -verbose=vvv emdr
else
    ./eve $@
fi
