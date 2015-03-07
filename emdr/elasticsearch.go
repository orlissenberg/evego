package emdr

import (
	"encoding/json"
	elastic "github.com/mattbaird/elastigo/lib"
	db "github.com/orlissenberg/evego/dbdump"
	"log"
)

type EveReader struct {
	*elastic.Conn
}

func (reader *EveReader) ReadRegion(id string) (region db.EveRegion, err error) {
	region = db.EveRegion{}
	err = reader.GetSource("eve", "region", id, nil, &region)

	return
}

func (reader *EveReader) ReadSolarSystem(id string) (system db.EveSolarSystem, err error) {
	system = db.EveSolarSystem{}
	err = reader.GetSource("eve", "solarsystem", id, nil, &system)

	return
}

func (reader *EveReader) ReadInvType(id string) (system db.EveInvType, err error) {
	system = db.EveInvType{}
	err = reader.GetSource("eve", "invtype", id, nil, &system)

	return
}

func (reader *EveReader) ReadStation(id string) (station db.EveStation, err error) {
	station = db.EveStation{}
	err = reader.GetSource("eve", "station", id, nil, &station)

	return
}

type ElasticEmdrWriter struct {
	*elastic.Conn
	*EveReader
}

func NewEveReader() *EveReader {
	reader := new(EveReader)
	reader.Conn = elastic.NewConn()
	reader.Hosts = []string{"localhost"}

	return reader
}

func NewElasticWriter() *ElasticEmdrWriter {
	writer := new(ElasticEmdrWriter)
	writer.Conn = elastic.NewConn()
	writer.Hosts = []string{"localhost"}

	writer.EveReader = new(EveReader)
	writer.EveReader.Conn = writer.Conn

	return writer
}

func (writer *ElasticEmdrWriter) Write(message []byte) (err error) {
	var v EmdrMessage
	json.Unmarshal(message, &v)

	switch v.ResultType {
	case "orders":
		err = writer.WriteOrder(message)
	case "history":
		// Soon! ...
		// err = writer.WriteHistory(message)
	}

	return
}

func (writer *ElasticEmdrWriter) WriteOrder(message []byte) (err error) {
	order := new(EmdrOrderMessage)
	json.Unmarshal(message, order)
	order.mapRows(writer.EveReader)

	for _, s := range order.RowSets {
		for _, o := range s.DataRows {
			_, err = writer.Index("eve", "order", "", nil, o)
		}
	}

	log.Println("Order")
	return
}

func (writer *ElasticEmdrWriter) WriteHistory(message []byte) (err error) {
	history := new(EmdrHistoryMessage)
	json.Unmarshal(message, history)
	history.mapRows(writer.EveReader)

	log.Println("History")
	_, err = writer.Index("eve", "history", "", nil, history)

	return
}

func (writer *ElasticEmdrWriter) DeleteAll() (err error) {
	log.Fatalln("NOT_IMPLEMENTED")

	return
}

