package main

import (
	"fmt"
	"strconv"
	"os"
	"net/http"
    "github.com/gorilla/mux"
    "firstgoprog/api" // Replace "firstgoprog" with your actual module name.
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/joho/godotenv/autoload"
	"github.com/srinathgs/mysqlstore"
)
var store *mysqlstore.MySQLStore

func establishdbConnection(user string, password string, host string, port string, database string) (*sql.DB, error) {
	db, err := sql.Open("mysql", user+":"+password+"@("+host+":"+port+")/"+database+"?parseTime=true")
	if err != nil {
		return nil, err
	}
	return db, nil
}


func main() {
	port := 8080
	// Convert the integer port to a string.
	portStr := strconv.Itoa(port)
	
	
	var err error
    store, err = mysqlstore.NewMySQLStore("root:daddy@(db:3306)/food?parseTime=true", "sessions", "/", 3600, []byte(os.Getenv("SESSION_KEY")))
    if err != nil {
      panic(err)
    }
    defer store.Close()
	

	r := mux.NewRouter()
	
	db, err := establishdbConnection("root", "daddy", "db", "3306", "food")
	
    if err != nil {
		// log.Fatal(err)
		fmt.Print("error connecting to db")
    }
    if err := db.Ping(); err != nil {
		fmt.Printf("Error %d...\n", err)
    }

	// Your other application routes go here...
    // Use the functions from the 'api' package to define routes.
	api.InitRoutes(r, db, store)

	fmt.Printf("Server is listening on port %d...\n", port)

	http.ListenAndServe(":"+portStr, r)
}