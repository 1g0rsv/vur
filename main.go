package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

// UserInput структура для получения текста от пользователя
type UserInput struct {
	Text string `json:"text"`
}

var aesKey []byte

func main() {
	// Преобразование шестнадцатеричного ключа в байты
	var err error
	hexKey := os.Getenv("HEX_KEY")
	aesKey, err = hex.DecodeString(hexKey)
	if err != nil {
		log.Fatal("Error decoding AES key:", err)
	}
	if len(aesKey) != 32 {
		log.Fatal("Invalid AES key length: expected 32 bytes, got", len(aesKey))
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	http.HandleFunc("/save", saveText)
	http.HandleFunc("/text/", getText)
	log.Println("Server is starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func encrypt(text string) (string, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	plaintext := []byte(text)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return hex.EncodeToString(ciphertext), nil
}

func decrypt(encryptedText string) (string, error) {
	ciphertext, err := hex.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func saveText(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "POST" {
		var userText UserInput
		err := json.NewDecoder(r.Body).Decode(&userText)
		if err != nil {
			log.Printf("Error decoding request: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		encryptedText, err := encrypt(userText.Text)
		if err != nil {
			log.Printf("Error encrypting text: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db := dbConn()
		newUUID := uuid.New().String()
		insert, err := db.Prepare("INSERT INTO texts(uuid, encrypted_text) VALUES(?, ?)")
		if err != nil {
			log.Printf("Error preparing insert query: %v", err)
			panic(err.Error())
		}
		_, err = insert.Exec(newUUID, encryptedText)
		if err != nil {
			log.Printf("Error executing insert query: %v", err)
			panic(err.Error())
		}

		response := map[string]string{"uuid": newUUID}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		log.Println("Text saved successfully with UUID:", newUUID)
		defer db.Close()
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func getText(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	uuid := r.URL.Path[len("/text/"):]

	db := dbConn()
	var encryptedText string
	err := db.QueryRow("SELECT encrypted_text FROM texts WHERE uuid = ?", uuid).Scan(&encryptedText)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Text with UUID %s not found", uuid)
			http.Error(w, "Text not found", http.StatusNotFound)
		} else {
			log.Printf("Error querying database: %v", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	decryptedText, err := decrypt(encryptedText)
	if err != nil {
		log.Printf("Error decrypting text: %v", err)
		http.Error(w, "Failed to decrypt text", http.StatusInternalServerError)
		return
	}

	// Удаление записи после расшифровки
	_, err = db.Exec("DELETE FROM texts WHERE uuid = ?", uuid)
	if err != nil {
		log.Printf("Error deleting text: %v", err)
		http.Error(w, "Failed to delete text", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, decryptedText)
	log.Printf("Text with UUID %s was read and deleted", uuid)
	defer db.Close()
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := os.Getenv("VUR_USER")
	dbPass := os.Getenv("VUR_PASS")
	dbName := os.Getenv("VUR_DB")
	// db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1)/"+dbName+"?parseTime=true")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(db)/"+dbName+"?parseTime=true")

	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		panic(err.Error())
	}
	return db
}
