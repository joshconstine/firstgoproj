package api

import (
    "net/http"
    "fmt"
	"log"	
    "database/sql"
	"html/template"
    _ "github.com/go-sql-driver/mysql"
)

type Recipe struct {
	Recipe_id int
	Name  string
	Description string
}
type RecipesPageData struct {
	PageTitle string
    Recipes []Recipe
    Ingredients []Ingredient
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
// GetRecipes is a sample handler for listing all recipes.
func GetRecipes(w http.ResponseWriter, r *http.Request, db *sql.DB) {
   tmpl := template.Must(template.ParseFiles("public/recipes.html"))
		recipes := getAllRecipes(db)
        ingredients := getAllIngredients(db)
		data := RecipesPageData{
			PageTitle: "Recipes",
            Recipes: recipes,
            Ingredients: ingredients,
        }

        tmpl.Execute(w, data)
}
func CreateRecipe(w http.ResponseWriter, r *http.Request, db *sql.DB) {
  	// Retrieve the form data
			recipeName := r.FormValue("recipeName")
			recipeDescription := r.FormValue("recipeDescription")
			// Retrieve the selected ingredients
			ingredientIDs := r.Form["ingredients"]



		// Perform the SQL INSERT query to add the ingredient to the database
		_, err := db.Exec("INSERT INTO recipes (name, description) VALUES (?, ?)", recipeName, recipeDescription )
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Retrieve the auto-generated recipe_id for the newly inserted recipe
		var recipeID int
		err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&recipeID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
 // Insert the selected ingredients into the recipe_ingredients table
 for _, ingredientID := range ingredientIDs {
	_, err = db.Exec("INSERT INTO recipe_ingredients (recipe_id, ingredient_id) VALUES (?, ?)", recipeID, ingredientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}



		// Redirect back to the home page
		fmt.Fprintf(w, `<script>window.location.href = "/recipes/%d";</script>`, recipeID)
}


// // GetRecipe is a sample handler for retrieving a recipe by ID.
// func GetRecipe(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     recipeID := vars["id"]
//     // Implement logic to retrieve a recipe by ID.
//     // Respond with the requested recipe.
//     // Example:
//     recipe := Recipe{ID: recipeID, Title: "Recipe " + recipeID}
//     fmt.Fprintf(w, "Recipe: %+v\n", recipe)
// }
