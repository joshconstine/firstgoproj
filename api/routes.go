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

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "public/index.html") 
    }).Methods("GET")		
    r.HandleFunc("/ingredients", func(w http.ResponseWriter, r *http.Request) {
        GetIngredientsTemplate(w, r, db)
    }).Methods("GET")	
	r.HandleFunc("/recipes", func(w http.ResponseWriter, r *http.Request) {
        GetRecipeTemplate(w, r, db)
    }).Methods("GET")	
	r.HandleFunc("/recipes/{id}", func(w http.ResponseWriter, r *http.Request) {
        GetRecipeById(w, r, db)
    }).Methods("GET")	
	r.HandleFunc("/create-recipe", func(w http.ResponseWriter, r *http.Request) {
        GetCreateRecipeTemplate(w, r, db)
    }).Methods("GET")	

	r.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
        GetListTemplate(w, r, db)
    }).Methods("GET")	
	r.HandleFunc("/generate-list", func(w http.ResponseWriter, r *http.Request) {
        GetGenerateListTemplate(w, r, db)
    }).Methods("POST")	


    apiRouter.HandleFunc("/ingredients", func(w http.ResponseWriter, r *http.Request) {
        CreateIngredient(w, r, db)
    }).Methods("POST")
    apiRouter.HandleFunc("/ingredients/delete", func(w http.ResponseWriter, r *http.Request) {
        DeleteIngredient(w, r, db)
    }).Methods("POST")
    apiRouter.HandleFunc("/ingredients/update", func(w http.ResponseWriter, r *http.Request) {
        UpdateIngredient(w, r, db)
    }).Methods("POST")


	apiRouter.HandleFunc("/recipes", func(w http.ResponseWriter, r *http.Request) {
        CreateRecipe(w, r, db)
    }).Methods("POST")
	apiRouter.HandleFunc("/recipes/delete", func(w http.ResponseWriter, r *http.Request) {
        DeleteRecipe(w, r, db)
    }).Methods("POST")
	apiRouter.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
        SendList(w, r)
    }).Methods("POST")
}
