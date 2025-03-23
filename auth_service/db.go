package authservice

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(ctx context.Context, connString string) (*mongo.Client, error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		log.Println("No deadline found")
		return nil, errors.New("No deadline found")
	}
	if time.Now().After(deadline) {
		log.Println("Deadline exceeded")
		return nil, errors.New("Deadline exceeded")
	}
	//  create a new client and connect to the server
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connString))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//  ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Println("Connected to MongoDB!")
	//  return the client
	return client, nil
}
