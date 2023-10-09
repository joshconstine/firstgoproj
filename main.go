package main

import (
	"fmt"
	"os"
	"io"
	"path/filepath"

	"github.com/google/uuid"
	"strconv"
	"net/http"
    "github.com/gorilla/mux"
    "firstgoprog/api" // Replace "firstgoprog" with your actual module name.
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/joho/godotenv/autoload"
)


func UploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 * 1024 * 1024) // 10 MB limit
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Generate a unique filename using a UUID
	fileExt := filepath.Ext(header.Filename)
	newFilename := uuid.New().String() + fileExt

	// Create a new file in the "public/static" directory with the unique filename
	newFilePath := filepath.Join("public/static/images", newFilename)
	newFile, err := os.Create(newFilePath)
	if err != nil {
		http.Error(w, "Failed to create a new file", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	// Copy the uploaded file to the new file
	_, err = io.Copy(newFile, file)
	if err != nil {
		http.Error(w, "Failed to copy the file", http.StatusInternalServerError)
		return
	}

	// Respond with the unique filename or other relevant information
	fmt.Fprintf(w, "File uploaded successfully with filename: %s", newFilename)
}

func main() {
	port := 8080
	// Convert the integer port to a string.
	portStr := strconv.Itoa(port)
	
	r := mux.NewRouter()
	r.HandleFunc("/upload", UploadHandler)

	db, err := sql.Open("mysql", "root:daddy@(db:3306)/food?parseTime=true")
	
    if err != nil {
		// log.Fatal(err)
		fmt.Print("error connecting to db")
    }
    if err := db.Ping(); err != nil {
		fmt.Printf("Error %d...\n", err)
    }
	staticDir := "/images/"
	r.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("./public/images"))))

	// Your other application routes go here...
    // Use the functions from the 'api' package to define routes.
	api.InitRoutes(r, db)

	fmt.Printf("Server is listening on port %d...\n", port)

	http.ListenAndServe(":"+portStr, r)
}