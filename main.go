package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
)

var db *sql.DB
var jwtKey = []byte("secret")

func connectDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
}

// LOGIN
func login(w http.ResponseWriter, r *http.Request) {
	var user map[string]string
	json.NewDecoder(r.Body).Decode(&user)

	if user["username"] == "admin" && user["password"] == "admin" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user": "admin",
		})

		tokenStr, _ := token.SignedString(jwtKey)
		json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
		return
	}

	http.Error(w, "Unauthorized", 401)
}

// GET ITEMS
func getItems(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.Query("SELECT id, name FROM items")

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

// CREATE ITEM
func createItem(w http.ResponseWriter, r *http.Request) {
	var item map[string]string
	json.NewDecoder(r.Body).Decode(&item)

	db.Exec("INSERT INTO items(name) VALUES(?)", item["name"])
	w.Write([]byte("OK"))
}

func main() {

	connectDB()

	http.HandleFunc("/login", login)
	http.HandleFunc("/items", getItems)
	http.HandleFunc("/items/create", createItem)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Running on port", port)
	http.ListenAndServe(":"+port, nil)
}
