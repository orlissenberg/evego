package emdr

import (
	"fmt"
	"strings"
	"io"
	"io/ioutil"
	"compress/zlib"
	zmq "github.com/pebbe/zmq2"
	"time"
	"strconv"
	"encoding/json"
	"os"
	elastic "github.com/mattbaird/elastigo/lib"
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
}

type ElasticEmdrWriter struct {}

type EmdrMessage struct {
	ResultType string "resultType"
}

func (writer *ElasticEmdrWriter) Write(message []byte) (err error) {
	var v EmdrMessage
	json.Unmarshal(message, &v)

	switch v.ResultType {
	case "orders":
		err = writer.WriteOrder(message)
	case "history":
		err = writer.WriteHistory(message)
	}

	return
}

func UnixTimeStampString() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

type EmdrOrderRowSet struct {
	TypeId      int64 "typeID"
	//	TypeName    string
	RegionID    int64 "regionID"
	GeneratedAt string "generatedAt"
	Rows        [][]interface{} "rows" // Format not Elasticsearch friendly
	DataRows    []map[string]interface{}
}

type EmdrOrderMessage struct {
	ResultType  string "resultType"
	UploadKeys  []map[string]string "uploadKeys"
	Generator   map[string]string "generator"
	CurrentTime string "currentTime"
	Version     string "version"
	RowSets     []EmdrOrderRowSet "rowsets"
	Columns     []string "columns"
}

func (writer *ElasticEmdrWriter) WriteOrder(message []byte) (err error) {
	c := elastic.NewConn()
	c.Hosts = []string{"localhost"}
	fmt.Println("Order: " + UnixTimeStampString())

	order := new(EmdrOrderMessage)
	json.Unmarshal(message, order)

	// Rewrite the data sets to a key:value format.
	// Loop the sets.
	for setIndex, setValue := range order.RowSets {
		order.RowSets[setIndex].DataRows = make([]map[string]interface{}, len(order.RowSets[setIndex].Rows))

		// Loop the rows.
		for rowIndex, rowValue := range setValue.Rows {
			// Create a new mapping.
			mapping := make(map[string]interface{})
			for index, name := range order.Columns {
				mapping[name] = rowValue[index]
			}
			order.RowSets[setIndex].DataRows[rowIndex] = mapping
		}

		// Remove the old data
		order.RowSets[setIndex].Rows = nil
	}

	_, err = c.Index("eve", "order", "", nil, order)

	if err != nil {
		// Dump data to file system on error for inspection.
		ioutil.WriteFile("error_order_"+strconv.FormatInt(time.Now().Unix(), 10)+".json", message, 0644)
	}

	return
}

func (writer *ElasticEmdrWriter) WriteHistory(message []byte) (err error) {
	c := elastic.NewConn()
	c.Hosts = []string{"localhost"}
	fmt.Println("History: " + UnixTimeStampString())
	_, err = c.Index("eve", "history", "", nil, string(message))

	return
}

func (client *EmdrClient) Start(writer EmdrWriter) {
	errCount := 0

	for {
		msg, err := client.sub.Recv(zmq.SNDMORE)

		if err == nil {
			var output []byte
			output, err = ZlibDecode(msg)

			if err == nil {
				err = writer.Write(output)
			}
		}

		if err != nil {
			if (errCount < 10) {
				errCount++
			} else {
				os.Exit(1)
			}
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
