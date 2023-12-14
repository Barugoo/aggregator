package main

import (
	"log"
	"net/http"

	"github.com/Barugoo/twaiv/internal/aggregator"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	ServerAddr string `envconfig:"SERVER_ADDR" default:":8080"`
}

func main() {
	var c config
	envconfig.Process("", &c)

	log.Println("running http server on: " + c.ServerAddr)

	inputCh := make(chan *aggregator.RunSession)
	agg, err := aggregator.NewInMemory(inputCh, aggregator.WeekBucketSize)
	if err != nil {
		log.Fatalf("unable to init in-memory aggregator: %v", err)
	}

	http.HandleFunc("/analyze", AggregateWorkoutsHandle(inputCh, agg))
	http.ListenAndServe(c.ServerAddr, nil)
}
