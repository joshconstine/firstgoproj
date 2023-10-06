package api

import (
	"github.com/gorilla/mux" 
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
	"net/http"
)

// InitRoutes initializes routes and handlers.
func InitRoutes(r *mux.Router, db *sql.DB) {
	// Create a subrouter for the "/api" path.
	apiRouter := r.PathPrefix("/api").Subrouter()

	// Define routes for recipe handling.
	r.HandleFunc("/recipes", func(w http.ResponseWriter, r *http.Request) {
        GetRecipes(w, r, db)
    }).Methods("GET")	
	apiRouter.HandleFunc("/recipes", func(w http.ResponseWriter, r *http.Request) {
        CreateRecipe(w, r, db)
    }).Methods("POST")
	// apiRouter.HandleFunc("/recipes/{id}", GetRecipe).Methods("GET")
	// apiRouter.HandleFunc("/recipes/create", CreateRecipe).Methods("POST")
}
