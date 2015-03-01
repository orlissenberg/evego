package emdr

import (
	"encoding/json"
	elastic "github.com/mattbaird/elastigo/lib"
	db "github.com/orlissenberg/evego/dbdump"
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

func ReadRegion(id string) (region db.EveRegion, err error) {
	c := elastic.NewConn()
	c.Hosts = []string{"localhost"}

	region = db.EveRegion{}
	err = c.GetSource("eve", "region", id, nil, &region)

	return
}

func ReadSolarSystem(id string) (system db.EveSolarSystem, err error) {
	c := elastic.NewConn()
	c.Hosts = []string{"localhost"}

	system = db.EveSolarSystem{}
	err = c.GetSource("eve", "solarsystem", id, nil, &system)

	return
}

func (writer *ElasticEmdrWriter) WriteOrder(message []byte) (err error) {
	order := new(EmdrOrderMessage)
	json.Unmarshal(message, order)
	order.mapRows()

	c := elastic.NewConn()
	c.Hosts = []string{"localhost"}

	for _, s := range order.RowSets {
		for _, o := range s.DataRows {
			_, err = c.Index("eve", "order", "", nil, o)
		}
	}

	log.Println("Order")
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

func (writer *ElasticEmdrWriter) DeleteAll() (err error) {
	log.Fatalln("NOT_IMPLEMENTED")

	return
}

