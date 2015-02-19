package emdr

import (
	"strings"
	"io"
	"io/ioutil"
	"compress/zlib"
	"time"
	"strconv"
	"os"
	"fmt"
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

type EmdrWriter interface {
	Write(message []byte) (err error)
	WriteOrder(message []byte) (err error)
	WriteHistory(message []byte) (err error)
	DeleteAll() (err error)
}

type EmdrMessage struct {
	ResultType string "resultType"
}

func UnixTimeStampString() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func (client *EmdrClient) Start(writer EmdrWriter) {
	//	duration := time.Millisecond * time.Duration(300)
	//	time.Sleep(duration)
	//	data, _ := ioutil.ReadFile("zlib_encoded.json")

	errCount := 0
	for {
		msg, err := client.sub.Recv(0)

		if err == nil {
			var output []byte
			output, err = ZlibDecode(msg)

			if err == nil {
				writer.Write(output)
			}
		}

		if err != nil {
			if (errCount < 10) {
				fmt.Printf("E%d.", errCount)
				errCount++
			} else {
				fmt.Println("Flippin' table, goodbye")
				os.Exit(1)
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

func (client *EmdrClient) Connect() (err error) {
	err = client.sub.Connect(client.address)
	client.sub.SetSubscribe("")
	return
}

func (client *EmdrClient) Close() {
	client.sub.Close()
}
