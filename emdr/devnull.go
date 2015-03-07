package emdr

import (
	"fmt"
)

type DevNullWriter struct {}

func (writer *DevNullWriter) Write(message []byte) (err error) {
	fmt.Print(".")
	return
}

func (writer *DevNullWriter) WriteOrder(message []byte) (err error) {
	return
}

func (writer *DevNullWriter) WriteHistory(message []byte) (err error) {
	return
}

func (writer *DevNullWriter) DeleteAll() (err error){
	return
}

