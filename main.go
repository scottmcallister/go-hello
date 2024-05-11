package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	DB *sql.DB
}

type Grocery struct {
	ID     int
	Name   string
	InCart bool
}

func main() {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := &App{DB: db}

	http.HandleFunc("/", app.Index)
	http.HandleFunc("/add", app.Add)
	http.HandleFunc("/delete", app.Delete)

	fmt.Println("listening on port 8080")
	http.ListenAndServe(":8080", nil)

}
