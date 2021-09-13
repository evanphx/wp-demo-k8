package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8081"
	}

	dbAddr := os.Getenv("DATABASE_URL")

	var err error
	if dbAddr != "" {
		dbAddr = "database=hashiconf-demo"
	}

	log.Println("creating database handle")

	db, err = sql.Open("postgres", dbAddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening for requests. addr=%s\n", addr)

	http.ListenAndServe(addr, http.HandlerFunc(handle))
}

func handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("request: %s %s", r.Method, r.URL.Path)

	err := db.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error servicing request: %s\n", err)
		return
	}

	fmt.Println(w, "app is ok!")
}