package main

import (
	"os"
	"path/filepath"
	"log"
	"io/ioutil"
	"flag"
	"emdr"
	"strings"
	db "sqlite"
	yaml "gopkg.in/yaml.v2"
	elastic "github.com/mattbaird/elastigo/lib"
)

type EveSettings map[string]string

func StartEmdrClient(settings EveSettings) {
	emdrClient, err := emdr.NewEmdr(settings["emdr_relay"])

	if err != nil {
		log.Fatalln("Failed to connect.")
	}

	emdr.ElasticsearchHosts = []string{settings["es_host"]}
	emdrClient.Start(emdr.NewElasticWriter())
}

func StartSqliteTransfer(settings EveSettings) {
	path := settings["sqlite_path"]

	c := elastic.NewConn()
	c.Hosts = []string{settings["es_host"]}

	db.ReadRegions(&db.ElasticRegionWriter{Conn: c}, path)
	db.ReadSolarSystems(&db.ElasticSolarSystemWriter{Conn: c}, path)
}

func ReadSettings() (settings EveSettings) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	var dat []byte;
	dat, err = ioutil.ReadFile(dir+"/eve.yaml")
	settings = make(map[string]string)
	yaml.Unmarshal(dat, &settings)

	return settings
}

func main() {
	// Load commandline flags.
	host := flag.String("es_host", "127.0.0.1:9200", "Elasticsearch host.")
	verbose := flag.String("verbose", "", "Verbosity level.")
	flag.Parse()
	mode := flag.Args()

	// Log settings.
	log.SetFlags(log.Lmicroseconds)
	if *verbose == "" {
		log.SetOutput(ioutil.Discard)
	}

	// Load configuration.
	settings := ReadSettings()
	settings["verbose"] = *verbose
	settings["es_host"] = *host

	// Show runtime settings.
	if strings.HasPrefix(*verbose, "v") {
		log.Println(mode, *host, settings)
	}

	// Execute mode, examples:
	// 	 ./run.sh -verbose=v sqlite
	// 	 ./run.sh -verbose=vvv emdr
	if len(mode) >= 1 {
		switch mode[0] {
		case "emdr":
			StartEmdrClient(settings)
		case "sqlite":
			StartSqliteTransfer(settings)
		}
	}
}
