package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var tmpl *template.Template

func init() {
	// Parse HTML templates
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Route handlers
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/calculate", calculateHandler)

	// Start server
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl.ExecuteTemplate(w, "layout.html", nil)
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form values
	num1Str := r.FormValue("num1")
	num2Str := r.FormValue("num2")
	operator := r.FormValue("operator")

	// Convert string inputs to float
	num1, err1 := strconv.ParseFloat(num1Str, 64)
	num2, err2 := strconv.ParseFloat(num2Str, 64)

	// Check for conversion errors
	if err1 != nil || err2 != nil {
		tmpl.ExecuteTemplate(w, "result.html", map[string]interface{}{
			"Error": "Please enter valid numbers",
		})
		return
	}

	// Perform calculation
	var result float64
	var errorMsg string

	switch operator {
	case "add":
		result = num1 + num2
	case "subtract":
		result = num1 - num2
	case "multiply":
		result = num1 * num2
	case "divide":
		if num2 == 0 {
			errorMsg = "Cannot divide by zero"
		} else {
			result = num1 / num2
		}
	default:
		errorMsg = "Invalid operation"
	}

	// Return result
	if errorMsg != "" {
		tmpl.ExecuteTemplate(w, "result.html", map[string]interface{}{
			"Error": errorMsg,
		})
	} else {
		tmpl.ExecuteTemplate(w, "result.html", map[string]interface{}{
			"Result": result,
		})
	}
}

