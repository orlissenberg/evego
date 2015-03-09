
SET GOPATH="%cd%"
echo %PATH%

del wineve.exe
go build wineve.go

::wineve.exe sqlserver -verbose=v !! WILL NOT WORK !!
wineve.exe -verbose="v" -es_host="192.168.33.48:9200" sqlserver
::wineve.exe sqlserver