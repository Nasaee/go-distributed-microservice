package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "8080"

type Application struct{}

func main() {
	app := Application{}
	log.Printf("Starting broker service on port %s", webPort)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
