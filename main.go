package main

import (
	"fmt"
	"strconv"
	"net/http"
	"os"
	"io"
    "github.com/gorilla/mux"
    "firstgoprog/api" // Replace "firstgoprog" with your actual module name.
	"database/sql"
    _ "github.com/go-sql-driver/mysql"

)


func main() {
	port := 8080
	// Convert the integer port to a string.
	portStr := strconv.Itoa(port)

	r := mux.NewRouter()

    // Use the functions from the 'api' package to define routes.
	api.InitRoutes(r)

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
	db, err := sql.Open("mysql", "root:daddy@(127.0.0.1:3306)/?parseTime=true")
    if err != nil {
        // log.Fatal(err)
		fmt.Print("error connecting to db")
    }
    if err := db.Ping(); err != nil {
	fmt.Printf("Error %d...\n", err)
    }

	fmt.Printf("Server is listening on port %d...\n", port)
	
	http.ListenAndServe(":"+portStr, r)
}