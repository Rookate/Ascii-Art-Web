package main

import (
	"fmt"
	"net/http"
)

const port = "8080"

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/generator", Generator)
	http.HandleFunc("/result", Result)
	fmt.Printf("Server started on http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
