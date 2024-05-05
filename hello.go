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
	http.HandleFunc("/add", app.add)
	http.ListenAndServe(":8080", nil)

}

func (app *App) index(w http.ResponseWriter, r *http.Request) {
	// fetch data from groceries table in database
	data := app.getGroceriesFromDB()

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *App) getGroceriesFromDB() []Grocery {
	rows, err := app.DB.Query("SELECT * FROM groceries")
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	data := rowsToGroceries(rows)

	return data
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

func (app *App) add(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	_, err := app.DB.Exec("INSERT INTO groceries (name, inCart) VALUES (?, ?)", name, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := app.getGroceriesFromDB()
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
