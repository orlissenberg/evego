package main

import (
	"os"
	"path/filepath"
	"log"
	"io/ioutil"
	"flag"
	"emdr"
	db "sqlite"
	"strings"
	yaml "gopkg.in/yaml.v2"
)

type EveSettings map[string]string

func StartEmdrClient(settings EveSettings) {
	emdrClient, err := emdr.NewEmdr(settings["emdr_relay"])

	if err != nil {
		log.Fatalln("Failed to connect.")
	}

	emdrClient.Start(emdr.NewElasticWriter())
}

func StartSqliteTransfer(settings EveSettings) {
	path := settings["sqlite_path"]
	db.ReadRegions(new(db.ElasticRegionWriter), path)
	db.ReadSolarSystems(new(db.ElasticSolarSystemWriter), path)
}

func ReadSettings() (settings EveSettings) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	var dat []byte;
	dat, err = ioutil.ReadFile(dir + "/eve.yaml")
	settings = make(map[string]string)
	yaml.Unmarshal(dat, &settings)

	return settings
}

func main() {
	// Load commandline flags.
	elastic := flag.String("es_host", "127.0.0.1:9200", "Elasticsearch host.")
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

	// Show runtime settings.
	if strings.HasPrefix(*verbose, "v") {
		log.Println(mode, *elastic, settings)
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
