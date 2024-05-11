package main

import (
	"html/template"
	"net/http"
)

func (app *App) Index(w http.ResponseWriter, r *http.Request) {
	data := app.getGroceriesFromDB()

	RenderTemplate(w, []string{"index.html", "list.html"}, data)
}

func (app *App) Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	_ = app.addGroceryToDB(name)
	data := app.getGroceriesFromDB()

	RenderTemplate(w, []string{"list.html"}, data)
}

func (app *App) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	_ = app.deleteGroceryFromDB(id)
	data := app.getGroceriesFromDB()

	RenderTemplate(w, []string{"list.html"}, data)
}

func RenderTemplate(w http.ResponseWriter, tmpls []string, data interface{}) {
	t, err := template.ParseFiles(tmpls...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
