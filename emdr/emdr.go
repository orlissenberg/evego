package emdr

import (
	"strings"
	"io"
	"io/ioutil"
	"compress/zlib"
	"time"
	"strconv"
	"log"
	zmq "github.com/pebbe/zmq2"
)

type EmdrClient struct {
	*zmq.Socket
	Address string
}

func NewEmdr(address string) (client *EmdrClient, err error) {
	client = new(EmdrClient)
	client.Address = address
	client.Socket, err = zmq.NewSocket(zmq.SUB)

	return
}

type EmdrWriter interface {
	Write(message []byte) (err error)
	WriteOrder(message []byte) (err error)
	WriteHistory(message []byte) (err error)
	DeleteAll() (err error)
}

func UnixTimeStampString() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func DumpToFile(data []byte) {
	ioutil.WriteFile("error_order_"+UnixTimeStampString()+".json", data, 0644)
}

func (client *EmdrClient) Start(writer EmdrWriter) {
	errCount := 0

	err := client.Connect(client.Address)
	client.SetSubscribe("")
	if err != nil {
		log.Fatalln("Failed to connect: " + err.Error())
	}

	for {
		msg, err := client.Recv(0)

		if err == nil {
			var output []byte
			output, err = ZlibDecode(msg)

			if err == nil {
				writer.Write(output)
			}
		}

		if err != nil {
			if (errCount < 10) {
				log.Println(err.Error)
				errCount++;
			} else {
				log.Fatalln(err.Error)
			}
		}
	}
}

func ZlibDecode(encoded string) (decoded []byte, err error) {
	var stringReader io.Reader
	stringReader = strings.NewReader(encoded)

	var pipeline io.ReadCloser
	pipeline, err = zlib.NewReader(stringReader)

	if err == nil {
		defer pipeline.Close()
		decoded, err = ioutil.ReadAll(pipeline)
	}

	return
}
