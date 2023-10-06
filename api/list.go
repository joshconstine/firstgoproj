package api

import (
	"fmt"
	"strconv"
    "net/http"
    "database/sql"
	"html/template"
    _ "github.com/go-sql-driver/mysql"
)
type CreateListPageData struct {
	PageTitle string
    Recipes []Recipe
}
type ListPageData struct {
	PageTitle string
    Ingredients []Ingredient
}

//HTML TEMPLATES

func GetListTemplate(w http.ResponseWriter, r *http.Request, db *sql.DB) {
  	tmpl := template.Must(template.ParseFiles("public/makeList.html"))
	
		recipes := getAllRecipes(db)
	
		data := CreateListPageData{
			PageTitle: "Make a List",
            Recipes: recipes,
        }

        tmpl.Execute(w, data)
}



func GetGenerateListTemplate(w http.ResponseWriter, r *http.Request, db *sql.DB) {
  // REcipe ids only reads if there is another form value? am i an idot
			recipeName := r.FormValue("recipeName")
			recipeIds := r.Form["recipes"]
		// Define a slice to hold all ingredients
		var ingredients []Ingredient
		// Retrieve the selected ingredients

// Log the selected ingredient IDs
		fmt.Println("", recipeName)

		// Iterate through the selected recipe IDs
		for _, recipeID := range recipeIds {
			// Convert the recipeID string to an integer
			recipeIDInt, err := strconv.Atoi(recipeID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
	
			// Query the ingredients for the current recipe
			rows, err := db.Query("SELECT i.name FROM ingredients i INNER JOIN recipe_ingredients ri ON i.ingredient_id = ri.ingredient_id WHERE ri.recipe_id = ?", recipeIDInt)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()
	
			// Loop through the rows of ingredients and append them to the list
			for rows.Next() {
				var ingredientName string
				err := rows.Scan(&ingredientName)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				ingredients = append(ingredients, Ingredient{Name: ingredientName})
			}
	
			// Check for errors during rows iteration
			if err := rows.Err(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	
		// Now, the 'ingredients' slice contains all ingredients from the selected recipes
		// You can use 'ingredients' as needed, such as displaying them in the response
		// or performing further processing.
	
		tmpl := template.Must(template.ParseFiles("public/list.html"))
		
		data := ListPageData{
			PageTitle: "Your List",
            Ingredients: ingredients,
        }

        tmpl.Execute(w, data)
	}