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
type CreateRecipePageData struct {
	PageTitle string
    Ingredients []Ingredient
	MappedIngredients map[string][]IngredientAndType
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

//HTML TEMPLATES

func GetRecipeTemplate(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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


func GetCreateRecipeTemplate(w http.ResponseWriter, r *http.Request, db *sql.DB) {
  tmpl := template.Must(template.ParseFiles("public/createRecipe.html"))
	
        ingredients := getAllIngredients(db)
      
		ingredientTypeMap := getAllIngredientsWithTypes(db)
		
		data := CreateRecipePageData{
			PageTitle: "Create Recipe",
            Ingredients: ingredients,
			MappedIngredients: ingredientTypeMap,
        }

        tmpl.Execute(w, data)
}


// DB Transactions
func CreateRecipe(w http.ResponseWriter, r *http.Request, db *sql.DB) {
			recipeName := r.FormValue("recipeName")
			recipeDescription := r.FormValue("recipeDescription")
			ingredientIDs := r.Form["ingredients"]

		_, err := db.Exec("INSERT INTO recipes (name, description) VALUES (?, ?)", recipeName, recipeDescription )
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var recipeID int
		err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&recipeID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
 for _, ingredientID := range ingredientIDs {
	_, err = db.Exec("INSERT INTO recipe_ingredients (recipe_id, ingredient_id) VALUES (?, ?)", recipeID, ingredientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
		fmt.Fprintf(w, `<script>window.location.href = "/recipes/%d";</script>`, recipeID)
}
func DeleteRecipe(w http.ResponseWriter, r *http.Request, db *sql.DB) {
			
		id := r.FormValue("id")

		// Perform the SQL INSERT query to add the ingredient to the database
		stmt, err := db.Prepare("DELETE FROM recipes WHERE recipe_id = ?")
    if err != nil {
        // return err
    }
    defer stmt.Close()

    // Execute the SQL statement
    _, err = stmt.Exec(id)
    if err != nil {
        // return err
    }

		fmt.Fprintf(w, `<script>window.location.href = "/recipes";</script>`)
}

