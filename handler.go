package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
)

var text string
var banner string
var output string

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHeader(w, "error 404", http.StatusNotFound)
		return
	}
	t, err := template.ParseFiles("templates/home.html")
	if err != nil {
		errorHeader(w, "Unable to load template", http.StatusInternalServerError)
		fmt.Println("Error parsing template:", err)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		errorHeader(w, "Unable to execute template", http.StatusInternalServerError)
		fmt.Println("Error executing template:", err)
	}
}

func Result(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/result.html")
	if err != nil {
		errorHeader(w, "Unable to load template", http.StatusInternalServerError)
		fmt.Println("Error parsing template:", err)
		return
	}
	if r.Method != http.MethodPost {
		errorHeader(w, "Bad request", http.StatusBadRequest)
		return
	}
	text = r.FormValue("text")
	banner = r.FormValue("banner")
	output, err = Ascii_Art(text, banner)
	fmt.Println(banner)
	fmt.Println(text)

	if err != nil {
		errorHeader(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	err = t.Execute(w, map[string]interface{}{
		"Result": output,
	})
	if err != nil {
		errorHeader(w, "Unable to execute template", http.StatusInternalServerError)
		fmt.Println("Error executing template:", err)
		return
	}
}

func Ascii_Art(text, banner string) (string, error) {

	// Read the ASCII art template file
	input, err := os.ReadFile(banner + ".txt")
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		return "", err
	}
	theme := strings.Split(string(input), "\n")
	// Read the string to be converted to ASCII arts
	var result strings.Builder
	char := strings.Split(text, "\n")
	// Parcourez chaque ligne de la chaîne d'entrée
	for i := 0; i < len(char); i++ {
		if char[i] == "" {
			result.WriteString("\n")
			continue
		}
		lines := make([]string, 9) // Initialize an empty slice for each line
		for k := 1; k < 9; k++ {
			line := ""
			for l := 0; l < len(char[i]); l++ {
				//---------------/!\
				if char[i][l] == '\r' || char[i][l] == '\n' {
					continue // Ignore carriage return and line feed characters
				}
				if int(char[i][l])-32 < 0 || int(char[i][l])-32 >= len(theme)/9 {
					return "", fmt.Errorf("character out of range: %v", char[i][l])
				}
				line += theme[(int(char[i][l]-32))*9+k]
			}
			lines[k-1] = line // Store the line in the slice
		}
		//---------------/!\
		result.WriteString(strings.Join(lines, "\n")) // Join the lines with newline character
		result.WriteString("\n\n")                    // Add extra newline after each line
		// Générer l'art ASCII pour chaque caractère de la ligne
	}
	return result.String(), nil
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	result := []byte(output)
	name := strings.Split(text, " ")[0]
	w.Header().Set("Content-Disposition", "attachment; filename=result-"+name+".txt")
	w.Header().Set("Content-Type", "text/plain")
	w.Write(result)
}

func errorHeader(w http.ResponseWriter, msg string, n int) {
	w.Header().Set("Content-Type", "text/html")
	t, err := template.ParseFiles("templates/error.html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}
	err = t.Execute(w, map[string]interface{}{
		"error_message": msg,
		"error_code":    n,
	})
	if err != nil {
		fmt.Println("Error executing template:", err)
	}
}
