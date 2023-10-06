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
	"log"
	"html/template"

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

type IngredientPageData struct {
	PageTitle string
    Ingredients []Ingredient
	IngredientTypes []IngredientType
	MappedIngredients map[string][]IngredientAndType
}







func getAllIngredients(db *sql.DB)  []Ingredient {
	rows, err := db.Query(`SELECT * FROM ingredients`)
        if err != nil {
			log.Fatal(err)
        }
        defer rows.Close()
		
        var ingredients []Ingredient
        for rows.Next() {
			var i Ingredient
			
            err := rows.Scan(&i.Ingredient_id, &i.Name, &i.Ingredient_type_id)
            if err != nil {
				log.Fatal(err)
            }
            ingredients = append(ingredients, i)
        }
        if err := rows.Err(); err != nil {
			log.Fatal(err)
        }
		return ingredients
}
func getAllIngredientsWithTypes(db *sql.DB)  map[string][]IngredientAndType {
	  rows, err := db.Query(`
		  SELECT i.ingredient_id, i.name, i.ingredient_type_id, t.name AS ingredient_type_name
		  FROM ingredients i
		  JOIN ingredient_type t ON i.ingredient_type_id = t.ingredient_type_id
	  `)
	  if err != nil {
		  panic(err.Error())
	  }
	  defer rows.Close()
  
	  // Map to store ingredients grouped by ingredient type
	  ingredientTypeMap := make(map[string][]IngredientAndType)
  
	  for rows.Next() {
		  var ingredient IngredientAndType
  
		  err := rows.Scan(&ingredient.Ingredient_id, &ingredient.Name, &ingredient.Ingredient_type_id, &ingredient.Ingredient_type_name)
		  if err != nil {
			  panic(err.Error())
		  }
  
		  // Append the ingredient to the corresponding ingredient type
		  ingredientTypeMap[ingredient.Ingredient_type_name] = append(ingredientTypeMap[ingredient.Ingredient_type_name], ingredient)
	  }
	  return ingredientTypeMap
}
func getAllRecipes(db *sql.DB) []Recipe {
	rows, err := db.Query(`SELECT * FROM recipes`)
        if err != nil {
			log.Fatal(err)
        }
        defer rows.Close()
		
        var recipes []Recipe
        for rows.Next() {
			var r Recipe
			
            err := rows.Scan(&r.Recipe_id, &r.Name,&r.Description)
            if err != nil {
				log.Fatal(err)
            }
            recipes = append(recipes, r)
        }
        if err := rows.Err(); err != nil {
			log.Fatal(err)
        }
		return recipes
}
func getAllIngredientTypes(db *sql.DB) []IngredientType {
	rows, err := db.Query(`SELECT * FROM ingredient_type`)
        if err != nil {
			log.Fatal(err)
        }
        defer rows.Close()
		
        var ingredient_types []IngredientType
        for rows.Next() {
			var r IngredientType
			
            err := rows.Scan(&r.Ingredient_type_id, &r.Name)
            if err != nil {
				log.Fatal(err)
            }
            ingredient_types = append(ingredient_types, r)
        }
        if err := rows.Err(); err != nil {
			log.Fatal(err)
        }
		return ingredient_types
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
	
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("public/layout.html"))
	
        ingredients := getAllIngredients(db)
        ingredientTypes := getAllIngredientTypes(db)
		ingredientTypeMap := getAllIngredientsWithTypes(db)
  
		
		data := IngredientPageData{
			PageTitle: "Ingredients list",
            Ingredients: ingredients,
			IngredientTypes: ingredientTypes,
			MappedIngredients: ingredientTypeMap,
        }

        tmpl.Execute(w, data)
    })
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

 


//List


	
	

	fmt.Printf("Server is listening on port %d...\n", port)
	
	http.ListenAndServe(":"+portStr, r)
}