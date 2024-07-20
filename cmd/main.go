package main

import (
	"flag"
	"log"

	_ "github.com/codescalersinternships/secretnote-api-spa-eyadhussein/docs"
	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/api"
	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/storage"
)

// @title Secret Note API
// @version 1.0
// @description This is a sample server for managing notes.
// @BasePath /api
// @securityDefinitions.apikey Token
// @in cookie
// @name token
func main() {
	var listenAddr string
	flag.StringVar(&listenAddr, "listen-addr", ":8080", "server listen address")

	flag.Parse()

	store := storage.NewMySQL()

	err := store.Init()

	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(listenAddr, store)

	server.Run()
}
