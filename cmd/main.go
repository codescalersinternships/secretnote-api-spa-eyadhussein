package main

import (
	"flag"
	"log"
	"os"

	_ "github.com/codescalersinternships/secretnote-api-spa-eyadhussein/docs"
	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/api"
	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/storage"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

// @title Secret Note API
// @version 1.0
// @description This is a sample server for managing notes.
// @BasePath /api
// @securityDefinitions.apikey Token
// @in cookie
// @name token
func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("failed to load .env file")
	}

	var listenAddr string
	flag.StringVar(&listenAddr, "listen-addr", ":5000", "server listen address")

	flag.Parse()

	dbConfig := storage.NewConfig(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	store := storage.NewMySQL(
		dbConfig,
	)

	err = store.Init()

	if err != nil {
		log.Fatal(err)
	}

	secretKey := os.Getenv("JWT_SECRET_KEY")
	server := api.NewServer(listenAddr, store, secretKey, rate.Limit(1), 1)

	server.Run()
}
