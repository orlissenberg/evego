#!/bin/sh

rm -Rf ./pkg
rm -Rf ./bin
rm -Rf ./src/github.com
rm -Rf ./src/gopkg.in
rm -Rf ./src/golang.org

if [ -f "eve" ]
then
    rm eve
fi

if [ -f "wineve.exe" ]
then
    rm wineve.exe
fi