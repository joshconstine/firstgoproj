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
	
 
	r.HandleFunc("/add-ingredient", func(w http.ResponseWriter, r *http.Request) {
			// Retrieve the form data
			ingredient := r.FormValue("ingredient")
			ingredientType := r.FormValue("ingredient_type")
		
		
		// Perform the SQL INSERT query to add the ingredient to the database
		_, err = db.Exec("INSERT INTO ingredients (name, ingredient_type_id) VALUES (?, ?)", ingredient, ingredientType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		// Redirect back to the home page
		fmt.Fprintf(w, `<script>window.location.href = "/";</script>`)
	})
	r.HandleFunc("/delete-ingredient", func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
	
		// Start a transaction
		tx, err := db.Begin()
		if err != nil {
			// Handle the error
			return
		}
		defer tx.Rollback()
	
		
	
		// Delete the rows where this ingredient is used as a foreign key in other tables
		// You need to execute additional DELETE statements for each table that references ingredients
	
		// For example, if there is a table 'recipe_ingredients' with 'ingredient_id' as a foreign key:
		stmt, err := tx.Prepare("DELETE FROM recipe_ingredients WHERE ingredient_id = ?")
		if err != nil {
			// Handle the error
			return
		}
		defer stmt.Close()
	
		_, err = stmt.Exec(id)
		if err != nil {
			// Handle the error
			return
		}
	// Delete the ingredient from the ingredients table
		stmt, err = tx.Prepare("DELETE FROM ingredients WHERE ingredient_id = ?")
		if err != nil {
			// Handle the error
			return
		}
		defer stmt.Close()
	
		_, err = stmt.Exec(id)
		if err != nil {
			// Handle the error
			return
		}
		// Commit the transaction
		if err := tx.Commit(); err != nil {
			// Handle the error
			return
		}
	
		// Redirect back to the home page
		fmt.Fprintf(w, `<script>window.location.href = "/";</script>`)
	})
	
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