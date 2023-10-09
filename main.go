package main

import (
	"fmt"
	

	"strconv"

	"net/http"
    "github.com/gorilla/mux"
    "firstgoprog/api" // Replace "firstgoprog" with your actual module name.
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/joho/godotenv/autoload"

)

func main() {
	port := 8080
	// Convert the integer port to a string.
	portStr := strconv.Itoa(port)
	
	

	r := mux.NewRouter()
	// r.HandleFunc("/upload", UploadHandler(uploader))
	


	db, err := sql.Open("mysql", "root:daddy@(db:3306)/food?parseTime=true")
	
    if err != nil {
		// log.Fatal(err)
		fmt.Print("error connecting to db")
    }
    if err := db.Ping(); err != nil {
		fmt.Printf("Error %d...\n", err)
    }

	// Your other application routes go here...
    // Use the functions from the 'api' package to define routes.
	api.InitRoutes(r, db)

	fmt.Printf("Server is listening on port %d...\n", port)

	http.ListenAndServe(":"+portStr, r)
}