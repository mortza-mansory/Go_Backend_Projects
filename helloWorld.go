package main

import (
	"fmt"
	"net/http"
)

func HelloHandller(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 20; i++ {
		fmt.Fprintf(w, "Hello World!")
	}

}

func main() {
	http.HandleFunc("/", HelloHandller)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
