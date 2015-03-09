package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	elastic "github.com/mattbaird/elastigo/lib"
	"log"
	"strconv"
)

var ElasticsearchHosts []string

type RegionWriter interface {
	Write(region EveRegion)
}

type EveRegion struct {
	Id   int64
	Name string
}

func (region *EveRegion) String() string {
	return strconv.FormatInt(region.Id, 10) + " - " + region.Name
}

// Stdout writer
type StdoutRegionWriter struct {}

func (writer *StdoutRegionWriter) Write(region EveRegion) {
	fmt.Println(region.String())
}

// Elasticsearch writer
type ElasticRegionWriter struct {
	*elastic.Conn
}

func (writer *ElasticRegionWriter) Write(region EveRegion) {
	_, err := writer.Index("eve", "region", strconv.FormatInt(region.Id, 10), nil, region)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(region.String())
}

// Read from SQLite
func ReadRegions(writer RegionWriter, path string) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select regionId, regionName from mapRegions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var region = EveRegion{}
		rows.Scan(&region.Id, &region.Name)
		writer.Write(region)
	}
	rows.Close()
}

type SolarSystemWriter interface {
	Write(SolarSystem EveSolarSystem)
}

type EveSolarSystem struct {
	Id   int64
	Name string
}

func (SolarSystem *EveSolarSystem) String() string {
	return strconv.FormatInt(SolarSystem.Id, 10) + " - " + SolarSystem.Name
}

// Stdout writer
type StdoutSolarSystemWriter struct {}

func (writer *StdoutSolarSystemWriter) Write(SolarSystem EveSolarSystem) {
	fmt.Println(SolarSystem.String())
}

// Elasticsearch writer
type ElasticSolarSystemWriter struct {
	*elastic.Conn
}

func (writer *ElasticSolarSystemWriter) Write(SolarSystem EveSolarSystem) {
	_, err := writer.Index("eve", "solarsystem", strconv.FormatInt(SolarSystem.Id, 10), nil, SolarSystem)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(SolarSystem.String())
}

// Read from SQLite
func ReadSolarSystems(writer SolarSystemWriter, path string) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select SolarSystemId, SolarSystemName from mapSolarSystems")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var SolarSystem = EveSolarSystem{}
		rows.Scan(&SolarSystem.Id, &SolarSystem.Name)
		writer.Write(SolarSystem)
	}
	rows.Close()
}


