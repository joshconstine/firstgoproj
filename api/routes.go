package api

import (
	"github.com/gorilla/mux"
)

// InitRoutes initializes routes and handlers.
func InitRoutes(r *mux.Router) {
	// Create a subrouter for the "/api" path.
	apiRouter := r.PathPrefix("/api").Subrouter()

	// Define routes for recipe handling.
	apiRouter.HandleFunc("/recipes", GetRecipes).Methods("GET")
	apiRouter.HandleFunc("/recipes/{id}", GetRecipe).Methods("GET")
	apiRouter.HandleFunc("/recipes/create", CreateRecipe).Methods("POST")
}
