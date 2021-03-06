package main

import (
	"github.com/newrelic/go-agent"
	"github.com/edukorg/go-shortener/config"
	"github.com/edukorg/go-shortener/handler"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	address := "0.0.0.0:" + cfg.Port
	log.Printf("Running on %s", address)
	newRelic, err := newrelic.NewApplication(cfg.NewRelic)
	err = http.ListenAndServe(address, &handler.RedirectHandler{Client: cfg.RedisClient, Logger: cfg.Logger, NewRelic: newRelic})
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
