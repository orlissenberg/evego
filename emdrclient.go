package main

import (
	"log"
	// "io/ioutil"
	emdr "github.com/orlissenberg/evego/emdr"
)

func main() {
	emdrClient, err := emdr.NewEmdr("tcp://relay-us-central-1.eve-emdr.com:8050")
	// emdrClient, err := emdr.NewEmdr("tcp://relay-eu-germany-1.eve-emdr.com:8050")
	// emdrClient, err := emdr.NewEmdr("tcp://relay-eu-denmark-1.eve-emdr.com:8050")

	log.SetFlags(log.Lmicroseconds)
	// log.SetOutput(ioutil.Discard)

	if err != nil {
		log.Fatalln("Failed to connect.")
	}

	emdrClient.Start(new(emdr.ElasticEmdrWriter))
}
