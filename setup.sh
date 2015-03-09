#!/bin/sh

# Go
export GOROOT=/opt/go
export PATH=$PATH:/opt/go/bin

# Project path
export GOPATH=/vagrant/projects/evego

#go get github.com/pebbe/zmq3
go get github.com/pebbe/zmq2
go get github.com/mattbaird/elastigo
go get gopkg.in/yaml.v2
go get github.com/mattn/go-sqlite3
go get github.com/jinzhu/gorm
go get github.com/denisenkom/go-mssqldb
go get github.com/garyburd/redigo/redis
