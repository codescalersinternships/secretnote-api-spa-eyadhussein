package main

import (
	"flag"
	"log"
	"os"

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

	dbConfig := storage.NewConfig(
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	store := storage.NewMySQL(
		dbConfig,
	)

	err := store.Init()

	if err != nil {
		log.Fatal(err)
	}

	secretKey := os.Getenv("JWT_SECRET_KEY")
	server := api.NewServer(listenAddr, store, secretKey)

	server.Run()
}
