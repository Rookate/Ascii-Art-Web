package main

import (
	"fmt"
	"net/http"
)

const port = "8080"

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/result", Result)
	http.HandleFunc("/download", DownloadFile)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Printf("Server started on http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
