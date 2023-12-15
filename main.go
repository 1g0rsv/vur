package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Структура для получения текста от пользователя
type UserInput struct {
	Text string `json:"text"`
}

// Функция для включения CORS
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// Функция для подключения к базе данных
func dbConn() (db *sql.DB) {

	dbDriver := "mysql"
	dbUser := "root"     // Замените на своего пользователя
	dbPass := "12345678" // Ваш пароль
	dbName := "mydb"     // Имя вашей базы данных
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1)/"+dbName+"?parseTime=true")

	if err != nil {
		panic(err.Error())
	}
	return db
}

// Функция для сохранения зашифрованного текста
func saveText(w http.ResponseWriter, r *http.Request) {
	enableCors(&w) // Включение CORS для этого обработчика

	log.Println("Received a request to save text") // Логирование приходящего запроса

	if r.Method == "OPTIONS" {
		return // Обработка предварительных запросов CORS
	}

	if r.Method == "POST" {
		var userText UserInput
		err := json.NewDecoder(r.Body).Decode(&userText)
		if err != nil {
			log.Printf("Error decoding request: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hashedText := sha256.Sum256([]byte(userText.Text))
		db := dbConn()
		insert, err := db.Prepare("INSERT INTO texts(encrypted_text) VALUES(?)")
		if err != nil {
			log.Printf("Error preparing database query: %s", err)
			panic(err.Error())
		}
		_, err = insert.Exec(hex.EncodeToString(hashedText[:]))
		if err != nil {
			log.Printf("Error executing database query: %s", err)
			panic(err.Error())
		}

		fmt.Fprintf(w, "Text saved successfully")
		defer db.Close()
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/save", saveText)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
