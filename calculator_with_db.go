
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Operation string

var db *sql.DB

const (
	jam    Operation = "+"
	menha  Operation = "-"
	tagsim Operation = "/"
	zarb   Operation = "*"
)

type Request struct {
	A         int `json:"a"`
	B         int `json:"b"`
	Operation `json:"operation"`
}

type Respond struct {
	Result int `json:"result"`
}

func Jam(a int, b int) int {
	return a + b
}
func Menha(a int, b int) int {
	return a - b
}
func Tagsim(a int, b int) int {
	return a / b
}
func Zarb(a int, b int) int {
	return a * b
}

func InitSQLite() error {
	var err error
	db, err = sql.Open("sqlite3", "./calculator.db")
	if err != nil {
		return err
	}

	query := `
	CREATE TABLE IF NOT EXISTS calculations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		a INTEGER NOT NULL,
		b INTEGER NOT NULL,
		operation TEXT NOT NULL,
		result INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(query)
	return err
}

func DoCalc(ctx context.Context, req Request) (int, error) {
	var result int

	switch req.Operation {
	case "+":
		result = Jam(req.A, req.B)
	case "-":
		result = Menha(req.A, req.B)
	case "*":
		result = Zarb(req.A, req.B)
	case "/":
		if req.B == 0 {
			return 0, errors.New("division by zero")
		}
		result = Tagsim(req.A, req.B)
	default:
		return 0, errors.New("invalid operation")
	}

	query := `
        INSERT INTO calculations (a, b, operation, result)
        VALUES (?, ?, ?, ?)
    `

	_, err := db.ExecContext(
		ctx,
		query,
		req.A,
		req.B,
		req.Operation,
		result,
	)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	result, err := DoCalc(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Respond{Result: result})
}

func main() {
	if err := InitSQLite(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/calc", CalcHandler)
	log.Println("Server running on :4687")
	log.Fatal(http.ListenAndServe(":4687", nil))
}
