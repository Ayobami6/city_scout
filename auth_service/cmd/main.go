package main

import (
	authservice "auth_service"
	"context"
	"log"
	"time"
)

func main() {
	addr := authservice.GetEnv("ADDR", "localhost:8181")
	dbUrl := authservice.GetEnv("DB_URL", "mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	dbClient, err := authservice.ConnectDB(ctx, dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := dbClient.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	api := authservice.NewAPIServer(addr, dbClient)
	if err := api.Start(); err != nil {
		log.Fatal(err)
	}

}
