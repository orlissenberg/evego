package main

import (
	"strings"
	"os"
	"path/filepath"
	"log"
	"io/ioutil"
	"flag"
	"sqlserver"
	yaml "gopkg.in/yaml.v2"
)

type EveSettings map[string]string

func StartSQLServer(settings EveSettings) {
	sqlserver.Transfer(settings)
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
	settings["es_host"] = *elastic

	// Show runtime settings.
	if strings.HasPrefix(settings["verbose"], "v") {
		log.Println(mode, *elastic, settings)
	}

	// Execute mode.
	if len(mode) >= 1 {
		switch mode[0] {
		case "sqlserver":
			StartSQLServer(settings)
		}
	}
}
