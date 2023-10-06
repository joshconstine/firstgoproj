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



type Todo struct {
    Title string
    Done  bool
}



type TodoPageData struct {
	PageTitle string
    Todos     []Todo
}


type IngredientType struct {
	Ingredient_type_id int
	Name string
}
type Ingredient struct {
	Ingredient_id int
	Name  string
	Ingredient_type_id int
}
type IngredientAndType struct {
	Ingredient_id int
	Name  string
	Ingredient_type_id int
	Ingredient_type_name string

}
type Recipe struct {
	Recipe_id int
	Name  string
	Description string
}





func main() {
	port := 8080
	// Convert the integer port to a string.
	portStr := strconv.Itoa(port)
	
	r := mux.NewRouter()
	
	db, err := sql.Open("mysql", "root:daddy@(db:3306)/food?parseTime=true")
    if err != nil {
		// log.Fatal(err)
		fmt.Print("error connecting to db")
    }
    if err := db.Ping(); err != nil {
		fmt.Printf("Error %d...\n", err)
    }
    // Use the functions from the 'api' package to define routes.
	api.InitRoutes(r, db)
	
 

	
	r.HandleFunc("/update-ingredient", func(w http.ResponseWriter, r *http.Request) {
		
		ingredientID := r.FormValue("id")
		updatedName := r.FormValue("ingredientName")
	
		// Prepare and execute the SQL UPDATE statement
		_, err = db.Exec("UPDATE ingredients SET name = ? WHERE ingredient_id = ?", updatedName, ingredientID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		// Redirect back to the page or provide a response
		fmt.Fprintf(w, `<script>window.location.href = "/";</script>`)
	})
	
	fmt.Printf("Server is listening on port %d...\n", port)

	http.ListenAndServe(":"+portStr, r)
}