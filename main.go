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
	"github.com/google/uuid"
)

// UserInput структура для получения текста от пользователя
type UserInput struct {
	Text string `json:"text"`
}

// enableCors функция для включения CORS
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// dbConn функция для подключения к базе данных
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

// saveText функция для сохранения зашифрованного текста
func saveText(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "POST" {
		var userText UserInput
		err := json.NewDecoder(r.Body).Decode(&userText)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hashedText := sha256.Sum256([]byte(userText.Text))
		db := dbConn()
		newUUID := uuid.New().String()
		insert, err := db.Prepare("INSERT INTO texts(uuid, encrypted_text) VALUES(?, ?)")
		if err != nil {
			panic(err.Error())
		}
		_, err = insert.Exec(newUUID, hex.EncodeToString(hashedText[:]))
		if err != nil {
			panic(err.Error())
		}

		response := map[string]string{"uuid": newUUID}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		defer db.Close()
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// getText функция для получения и отображения текста по UUID
func getText(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	uuid := r.URL.Path[len("/text/"):]

	db := dbConn()
	query, err := db.Prepare("SELECT encrypted_text FROM texts WHERE uuid = ?")
	if err != nil {
		panic(err.Error())
	}
	defer query.Close()

	var encryptedText string
	err = query.QueryRow(uuid).Scan(&encryptedText)
	if err != nil {
		http.Error(w, "Text not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, encryptedText)
	defer db.Close()
}

func main() {
	http.HandleFunc("/save", saveText)
	http.HandleFunc("/text/", getText)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
