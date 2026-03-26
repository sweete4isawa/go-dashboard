package handlers

import (
	"encoding/json"
	"net/http"

	"go-dashboard/config"
)

func GetItems(w http.ResponseWriter, r *http.Request) {
	rows, _ := config.DB.Query("SELECT id, name FROM items")

	var items []map[string]interface{}

	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)

		items = append(items, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}

	json.NewEncoder(w).Encode(items)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	var item map[string]string
	json.NewDecoder(r.Body).Decode(&item)

	config.DB.Exec("INSERT INTO items(name) VALUES(?)", item["name"])

	w.Write([]byte("OK"))
}
