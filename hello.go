package main

import (
	"database/sql"
	"fmt"
	"html/template"
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

	// print hello world
	fmt.Println("listening on port 8080")

	http.HandleFunc("/", app.index)
	http.ListenAndServe(":8080", nil)

}

func (app *App) index(w http.ResponseWriter, r *http.Request) {
	// fetch data from groceries table in database
	rows, err := app.DB.Query("SELECT * FROM groceries")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := rowsToGroceries(rows)

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func rowsToGroceries(rows *sql.Rows) []Grocery {
	if rows == nil {
		log.Println("rows is nil")
		return nil
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil
	}
	var groceries []Grocery
	for rows.Next() {
		var grocery Grocery
		err := rows.Scan(&grocery.ID, &grocery.Name, &grocery.InCart)
		if err != nil {
			log.Println(err)
			continue
		}
		groceries = append(groceries, grocery)
	}
	fmt.Println(groceries)
	return groceries
}
