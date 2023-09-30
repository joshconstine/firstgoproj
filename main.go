package main

import (
	"fmt"
	"strconv"
	"net/http"
	"os"
	"io"
    "github.com/gorilla/mux"
)


func main() {
	port := 8080
	// Convert the integer port to a string.
	portStr := strconv.Itoa(port)

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		htmlFile, err := os.Open("public/index.html")
		if err != nil {
			http.Error(w, "Unable to read HTML file", http.StatusInternalServerError)
			return
		}
		defer htmlFile.Close()

		// Copy the HTML content to the response writer.
		_, err = io.Copy(w, htmlFile)
		if err != nil {
			http.Error(w, "Unable to copy HTML content to response", http.StatusInternalServerError)
			return
		}
	})
	r.HandleFunc("/recipes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		htmlFile, err := os.Open("public/recipes.html")
		if err != nil {
			http.Error(w, "Unable to read HTML file", http.StatusInternalServerError)
			return
		}
		defer htmlFile.Close()

		// Copy the HTML content to the response writer.
		_, err = io.Copy(w, htmlFile)
		if err != nil {
			http.Error(w, "Unable to copy HTML content to response", http.StatusInternalServerError)
			return
		}
	})
	r.HandleFunc("/recipes/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
        id := vars["id"]
		fmt.Fprintf(w, "You've requested the recipe: %s\n", id)

		// w.Header().Set("Content-Type", "text/html")
		// htmlFile, err := os.Open("public/recipes.html")
		// if err != nil {
		// 	http.Error(w, "Unable to read HTML file", http.StatusInternalServerError)
		// 	return
		// }
		// defer htmlFile.Close()

		// // Copy the HTML content to the response writer.
		// _, err = io.Copy(w, htmlFile)
		// if err != nil {
		// 	http.Error(w, "Unable to copy HTML content to response", http.StatusInternalServerError)
		// 	return
		// }
	})
	
	fmt.Printf("Server is listening on port %d...\n", port)
	r.HandleFunc("/hello-world", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello World"))
	})
	http.ListenAndServe(":"+portStr, r)
}