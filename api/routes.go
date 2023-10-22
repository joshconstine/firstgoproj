package api

import (
	"github.com/gorilla/mux" 
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
	"net/http"	
    "fmt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/srinathgs/mysqlstore"
)
func NewUploader() *s3manager.Uploader {
	ACCESS_KEY:= "AKIAX6ZNEPNPAR6OXRRO"
	SECRET_KEY:= "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"
	s3Config := &aws.Config{
		Region:      aws.String("us-west-1"),
		Credentials: credentials.NewStaticCredentials(ACCESS_KEY, SECRET_KEY, ""),
	}

	s3Session := session.New(s3Config)

	uploader := s3manager.NewUploader(s3Session)
	fmt.Printf("Created new S3 Uploder")
	
	return uploader
}

  func sessTest(w http.ResponseWriter, r *http.Request, store *mysqlstore.MySQLStore) {
    username := r.FormValue( "username")
    password := r.FormValue( "password")
    fmt.Printf("username: %s\n", username)
    fmt.Printf("password: %s\n", password)

    session, err := store.Get(r, username)
    session.Values["bar"] = "baz"
    session.Values["baz"] = "foo"
    err = session.Save(r, w)
    fmt.Printf("%#v\n", session)
    fmt.Println(err)
  }



// InitRoutes initializes routes and handlers.
func InitRoutes(r *mux.Router, db *sql.DB, store *mysqlstore.MySQLStore ) {
	// Create a subrouter for the "/api" path.
	apiRouter := r.PathPrefix("/api").Subrouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "public/index.html") 
    }).Methods("GET")		
    r.HandleFunc("/sign-in", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "public/auth.html") 
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
    r.HandleFunc("/update-ingredients", func(w http.ResponseWriter, r *http.Request) {
        UpdateIngredientsHandler(w, r, db)
    }).Methods("GET")	


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
    apiRouter.HandleFunc("/recipes/description", func(w http.ResponseWriter, r *http.Request) {
        UpdateDescription(w, r, db)
    }).Methods("POST")
	apiRouter.HandleFunc("/recipes/delete", func(w http.ResponseWriter, r *http.Request) {
        DeleteRecipe(w, r, db)
    }).Methods("POST")
	apiRouter.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
        SendList(w, r)
    }).Methods("POST")
	r.HandleFunc("/update_recipe_ingredients", func(w http.ResponseWriter, r *http.Request) {
        UpdateRecipeIngredients(w, r, db)
    }).Methods("POST")

    r.HandleFunc("/recipes/add-photo", func(w http.ResponseWriter, r *http.Request) {
        UploadNewPhoto(w, r, db)
    }).Methods("POST")	


    apiRouter.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
        sessTest(w, r, store)
    }).Methods("GET")

}
