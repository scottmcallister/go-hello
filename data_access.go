package main

import (
	"database/sql"
	"fmt"
	"log"
)

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

func (app *App) addGroceryToDB(name string) error {
	_, err := app.DB.Exec("INSERT INTO groceries (name, inCart) VALUES (?, ?)", name, false)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) deleteGroceryFromDB(id string) error {
	_, err := app.DB.Exec("DELETE FROM groceries WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) toggleGrocery(id string) error {
	_, err := app.DB.Exec("UPDATE groceries SET inCart = NOT inCart WHERE id = ?", id)
	if err != nil {
		return err
	}
	fmt.Println("toggling " + id)
	return nil
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
