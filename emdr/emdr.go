package emdr

import (
	"fmt"
	"strings"
	"io"
	"io/ioutil"
	"compress/zlib"
	zmq "github.com/pebbe/zmq2"
)

type EmdrClient struct {
	sub *zmq.Socket
	address string
}

func NewEmdr(address string) (client *EmdrClient, err error) {
	client = new(EmdrClient)
	client.sub, err = zmq.NewSocket(zmq.SUB)
	client.address = address
	return
}

func (client *EmdrClient) Start() {
	for count := 0; count < 10; {
		msg, err := client.sub.Recv(zmq.SNDMORE)

		if err != nil {
			fmt.Println(err)
			count++;
		} else {
			output, _ := ZlibDecode(msg)
			fmt.Println(string(output))
		}
	}
}

func ZlibDecode(encoded string) (decoded []byte, err error) {
	var pipeline io.ReadCloser
	var stringReader io.Reader
	stringReader = strings.NewReader(encoded)
	pipeline, err = zlib.NewReader(stringReader)
	defer pipeline.Close()

	if err == nil {
		decoded, err = ioutil.ReadAll(pipeline)
	}

	return
}

func (client *EmdrClient) Connect() (err error) {
	err = client.sub.Connect(client.address)
	client.sub.SetSubscribe("")
	return
}

func (client *EmdrClient) Close() {
	client.sub.Close()
}
