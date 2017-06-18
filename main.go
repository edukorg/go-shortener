package main

import (
	"fmt"
	"github.com/guilhermef/go-shortener/config"
	"github.com/guilhermef/go-shortener/handler"
	"log"
	"net/http"
	"os"

	newrelic "github.com/newrelic/go-agent"
)

func mustGetEnv(key string) string {
	if val := os.Getenv(key); "" != val {
		return val
	}
	panic(fmt.Sprintf("environment variable %s unset", key))
}

func getApp() newrelic.Application {
	if os.Getenv("NEW_RELIC_LICENSE_KEY") == "" {
		return nil
	}

	nrCfg := newrelic.NewConfig(mustGetEnv("NEW_RELIC_APP_NAME"), mustGetEnv("NEW_RELIC_LICENSE_KEY"))
	nrCfg.Logger = newrelic.NewDebugLogger(os.Stdout)

	var err error
	var app newrelic.Application
	app, err = newrelic.NewApplication(nrCfg)
	if err != nil {
		panic(err)
	}

	return app
}

func main() {
	app := getApp()

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	address := "0.0.0.0:" + cfg.Port
	log.Printf("Running on %s", address)
	if app != nil {
		http.HandleFunc(newrelic.WrapHandleFunc(app, "/", (&handler.RedirectHandler{Client: cfg.RedisClient, Logger: cfg.Logger}).ServeHTTP))
	} else {
		http.Handle("/", &handler.RedirectHandler{Client: cfg.RedisClient, Logger: cfg.Logger})
	}
	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
