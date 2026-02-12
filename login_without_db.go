package main

import (
	"encoding/json"
	"net/http"
)

type UserInfo struct {
	Name string
	Pass string
	Role string
}

var users = map[string]UserInfo{
	"admin": {
		Name: "admin",
		Pass: "yaro",
		Role: "hich kare",
	},
	"user": {
		Name: "user",
		Pass: "user",
		Role: "user",
	},
}

type LoginRequest struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

type UserDTO struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type LoginResponse struct {
	Message string  `json:"message"`
	User    UserDTO `json:"user"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "invalid json body",
		})
		return
	}

	user, exists := users[req.Name]
	if !exists || user.Pass != req.Pass {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{
			Error: "invalid credentials",
		})
		return
	}

	resp := LoginResponse{
		Message: "login successful",
		User: UserDTO{
			Name: user.Name,
			Role: user.Role,
		},
	}

	writeJSON(w, http.StatusOK, resp)
}

func main() {
	http.HandleFunc("/login", LoginHandler)
	http.ListenAndServe(":8080", nil)
}
