package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"logger/data"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const webPort = "8081"

type Application struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting logger Service")

	db := connectToDB()
	if db == nil {
		log.Fatal("Couldnot connect to Postgres")
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

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	// Retry/Polling Mechanism เพื่อรอให้ Postgres พร้อมใช้งานอย่างสมบูรณ์ก่อนที่จะเริ่มทำงาน
	for range 10 {
		db, err := sql.Open("pgx", dsn)
		if err == nil && db.Ping() == nil {
			log.Println("Connected to Postgres!")
			return db
		}
		log.Println("Waiting for postgres")
		time.Sleep(time.Second * 2)
	}
	return nil
}
