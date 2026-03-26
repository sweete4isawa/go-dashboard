package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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

	var dbUser string
	var dbPass string

	err := db.QueryRow("SELECT username, password FROM users WHERE username=?", user["username"]).Scan(&dbUser, &dbPass)

	if err != nil || dbPass != user["password"] {
		http.Error(w, "Unauthorized", 401)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": dbUser,
	})

	tokenStr, _ := token.SignedString(jwtKey)

	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenStr,
	})
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "No token", 401)
			return
		}

		tokenStr = strings.Replace(tokenStr, "Bearer ", "", 1)

		_, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			http.Error(w, "Invalid token", 401)
			return
		}

		next(w, r)
	}
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
	http.HandleFunc("/items", authMiddleware(getItems))
	http.HandleFunc("/items/create", authMiddleware(createItem))

	// 🌐 static file (HTML dashboard)
	fs := http.FileServer(http.Dir("./templates"))
	http.Handle("/", fs)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Running on port", port)
	http.ListenAndServe(":"+port, nil)
}
