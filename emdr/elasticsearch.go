package emdr

import (
	"encoding/json"
	elastic "github.com/mattbaird/elastigo/lib"
	"log"
)

type ElasticEmdrWriter struct {}

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

func (writer *ElasticEmdrWriter) WriteOrder(message []byte) (err error) {
	order := new(EmdrOrderMessage)
	json.Unmarshal(message, order)
	order.mapRows()

	c := elastic.NewConn()
	c.Hosts = []string{"localhost"}
	log.Println("Order")
	_, err = c.Index("eve", "order", "", nil, order)

	return
}

func (writer *ElasticEmdrWriter) WriteHistory(message []byte) (err error) {
	history := new(EmdrHistoryMessage)
	json.Unmarshal(message, history)
	history.mapRows()

	c := elastic.NewConn()
	c.Hosts = []string{"localhost"}
	log.Println("History")
	_, err = c.Index("eve", "history", "", nil, history)

	return
}

func (writer *ElasticEmdrWriter) DeleteAll() (err error){
	log.Fatalln("NOT_IMPLEMENTED")

	return
}

