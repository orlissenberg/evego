
SET GOPATH="%cd%"
echo %PATH%

del wineve.exe
go build wineve.go

SET HOST="192.168.33.48:9200"

::wineve.exe sqlserver -verbose=v !! WILL NOT WORK !!
wineve.exe -verbose="v" -es_host=%HOST% sqlserver
::wineve.exe sqlserver