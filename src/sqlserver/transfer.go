package sqlserver

import (
	"strings"
	"log"
	"strconv"
	"encoding/json"
	_ "github.com/denisenkom/go-mssqldb"
	gorm "github.com/jinzhu/gorm"
	elastic "github.com/mattbaird/elastigo/lib"
)

var verbose string

func Transfer(settings map[string]string) {
	verbose = settings["verbose"]

	db, _ := gorm.Open("mssql", settings["sqlserver_connection"])
	db.LogMode(true)
	db.DB()

	// Then you could invoke `*sql.DB`'s functions with it
	db.DB().Ping()

	// Disable table name's pluralization
	db.SingularTable(true)

	c := elastic.NewConn()
	c.Hosts = []string{settings["es_host"]}

	// Hiding dev code 〜(￣△￣〜)
	if false {
		firstType := EveInvType{}
		query := db.Find(&firstType, 19)
		firstType.ElasticWrite(c)

		firstStation := EveStation{}
		query = db.Find(&firstStation, 60000004)
		log.Println(firstStation.String())
		log.Println(query.Error)
	}

	// Import inv types
	var types []EveInvType
	db.Find(&types)
	for _, t := range types {
		t.ElasticWrite(c)
	}

	// Import stations
	var stations []EveStation
	db.Find(&stations)
	for _, s := range stations {
		s.ElasticWrite(c)
	}
}

func (e EveStation) TableName() string {
	return "eve_data.dbo.staStations"
}

func (station *EveStation) String() string {
	if verbose == "vvv" {
		data, _ := json.MarshalIndent(station, "", "    ");
		return string(data)
	}

	return station.Name
}

func (station *EveStation) ElasticWrite(c *elastic.Conn) {
	_, err := c.Index("eve", "station", strconv.FormatInt(int64(station.Id), 10), nil, station)
	if err != nil {
		log.Fatal(err)
	}

	if strings.HasPrefix(verbose, "v") {
		log.Println(station.String())
	}
}

func (e EveInvType) TableName() string {
	return "eve_data.dbo.invTypes"
}

func (invType *EveInvType) String() string {
	if verbose == "vvv" {
		data, _ := json.MarshalIndent(invType, "", "    ");
		return string(data)
	}

	return invType.Name
}

func (invType *EveInvType) ElasticWrite(c *elastic.Conn) {
	_, err := c.Index("eve", "invtype", strconv.FormatInt(int64(invType.Id), 10), nil, invType)
	if err != nil {
		log.Fatal(err)
	}

	if strings.HasPrefix(verbose, "v") {
		log.Println(invType.String())
	}
}
