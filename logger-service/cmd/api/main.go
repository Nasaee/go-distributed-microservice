package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"logger/data"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const webPort = "8082"

type Application struct {
	DB     *mongo.Client
	Models data.Models
}

func main() {
	log.Println("Starting logger Service")

	db := connectToDB()
	if db == nil {
		log.Fatal("Couldnot connect to Mongo")
	}

	app := Application{
		DB:     db,
		Models: data.New(db),
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Printf("Starting logger service on port %s", webPort)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}

func connectToDB() *mongo.Client {
	uri := os.Getenv("MONGO_URI")

	for range 10 {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		client, err := mongo.Connect(options.Client().ApplyURI(uri))
		if err == nil {
			err = client.Ping(ctx, nil)
			cancel()

			if err == nil {
				log.Println("Conncted to Mongo")
				return client
			}
		}

		cancel()

		log.Println("Waiting for Mongo")
		time.Sleep(2 * time.Second)
	}

	return nil
}
