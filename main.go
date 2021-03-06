package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	env string
)

func main() {
	env = os.Getenv("APP_ENV")
	if env == "" {
		log.Printf("No app env set, presuming development")
		env = "development"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	dbAddr := os.Getenv("DATABASE_URL")

	var err error
	if dbAddr == "" {
		dbAddr = "postgres://localhost/hashiconf-demo?sslmode=disable"
	}

	log.Printf("creating database handle: %s\n", dbAddr)

	db, err = sql.Open("postgres", dbAddr)
	if err != nil {
		log.Fatal(err)
	}

	addr := ":" + port

	log.Printf("listening for requests. addr=%s\n", addr)

	http.ListenAndServe(addr, http.HandlerFunc(handle))
}

func handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("request: %s %s", r.Method, r.URL.Path)

	err := db.Ping()
	if err != nil {
		log.Printf("error pinging database: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error servicing request: %s\n", err)
		return
	}

	fmt.Fprintf(w, "app is ok!\nenv=%s\n", env)
}
