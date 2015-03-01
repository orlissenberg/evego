package dbdump

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	elastic "github.com/mattbaird/elastigo/lib"
	"log"
	"strconv"
)

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

func (std *StdoutRegionWriter) Write(region EveRegion) {
	fmt.Println(region.String())
}

// Elasticsearch writer
type ElasticRegionWriter struct {}

func (std *ElasticRegionWriter) Write(region EveRegion) {
	c := elastic.NewConn()
	c.Hosts = []string{"localhost"}

	_, err := c.Index("eve", "region", strconv.FormatInt(region.Id, 10), nil, region)
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

func (std *StdoutSolarSystemWriter) Write(SolarSystem EveSolarSystem) {
	fmt.Println(SolarSystem.String())
}

// Elasticsearch writer
type ElasticSolarSystemWriter struct {}

func (std *ElasticSolarSystemWriter) Write(SolarSystem EveSolarSystem) {
	c := elastic.NewConn()
	c.Hosts = []string{"localhost"}

	_, err := c.Index("eve", "solarsystem", strconv.FormatInt(SolarSystem.Id, 10), nil, SolarSystem)
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


