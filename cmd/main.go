package main

import (
	"flag"
	"log"

	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/api"
	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/storage"
)

func main() {
	var listenAddr string
	flag.StringVar(&listenAddr, "listen-addr", ":8080", "server listen address")

	flag.Parse()

	store, err := storage.NewMySQL()

	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(listenAddr, store)

	server.Run()
}
