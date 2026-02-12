package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Operation string

const (
	jam    Operation = "+"
	menha  Operation = "-"
	tagsim Operation = "/"
	zarb   Operation = "*"
)

type Request struct {
	A int
	B int
	Operation
}

type Respond struct {
	Result int
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

func CalcHandller(w http.ResponseWriter, r *http.Request) {
	context.Context.
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid json body", http.StatusBadRequest)
		return
	}

	var result int

	if req.Operation != jam && req.Operation != menha && req.Operation != tagsim && req.Operation != zarb {
		http.Error(w, "Only + - * / is allowed for field Operation.", http.StatusBadRequest)
		return
	}

	if req.Operation == "+" {
		result = Jam(req.A, req.B)
	}
	if req.Operation == "-" {
		result = Menha(req.A, req.B)

	}
	if req.Operation == "*" {
		result = Zarb(req.A, req.B)

	}
	if req.Operation == "/" {
		result = Tagsim(req.A, req.B)

	}

	resp := Respond{
		Result: result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

func main() {
	http.HandleFunc("/cal", CalcHandller)

	err := http.ListenAndServe(":4687", nil)
	if err != nil {
		fmt.Print("Error on startation")
	}
	fmt.Printf("Backend is running")
}
