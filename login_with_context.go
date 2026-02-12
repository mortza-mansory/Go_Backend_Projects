package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type User struct {
	Username string
	Password string
}

var validUser = User{
	Username: "admin",
	Password: "password123",
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	username := r.FormValue("username")
	password := r.FormValue("password")

	resultChan := make(chan bool)
	go validateUser(ctx, username, password, resultChan)

	select {
	case <-ctx.Done():
		http.Error(w, "درخواست timeout خورد", http.StatusRequestTimeout)
		return
	case isValid := <-resultChan:
		if isValid {
			fmt.Fprintf(w, "خوش آمدید %s! لاگین موفقیت آمیز بود.", username)
		} else {
			http.Error(w, "نام کاربری یا رمز عبور اشتباه است", http.StatusUnauthorized)
		}
	}
}

func validateUser(ctx context.Context, username, password string, resultChan chan<- bool) {

	select {
	case <-ctx.Done():
		resultChan <- false
		return
	case <-time.After(2 * time.Second):
		if username == validUser.Username && password == validUser.Password {
			resultChan <- true
		} else {
			resultChan <- false
		}
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", loginHandler).Methods("POST")

	fmt.Println("سرور در حال اجرا روی پورت 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
