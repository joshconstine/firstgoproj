package api

import (
	"database/sql"
	"log"	
	"fmt"	
    "net/http"
)
type Ingredient struct {

	Ingredient_id int
	Name  string
	Ingredient_type_id int
}
type IngredientType struct {
	Ingredient_type_id int
	Name string
}
type QuantityType struct {
	Quantity_type_id int
	Name string
}
type IngredientAndType struct {
	Ingredient_id int
	Name  string
	Ingredient_type_id int
	Ingredient_type_name string

}
type IngredientWithQuantity struct {
	Ingredient_id int
	Name  string
	Quantity float32
	Quantity_type string
	Quantity_type_id int
}
type IngredientWithQuantityAndType struct {
	Ingredient_id int
	Name  string
	Quantity float32
	Quantity_type string
	Quantity_type_id int
	Ingredient_Type_id int	
	Ingredient_Type_Name string
}



type IngredientPageData struct {
	PageTitle string
    Ingredients []Ingredient
	IngredientTypes []IngredientType
	MappedIngredients map[string][]IngredientAndType
	User User
}






//DB Transactions

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
func getAllQuantityTypes(db *sql.DB) []QuantityType {
	rows, err := db.Query(`SELECT * FROM quantity_type`)
        if err != nil {
			log.Fatal(err)
        }
        defer rows.Close()
		
	var quantitiy_types []QuantityType
        for rows.Next() {
			var r QuantityType
			
            err := rows.Scan(&r.Quantity_type_id, &r.Name)
            if err != nil {
				log.Fatal(err)
            }
            quantitiy_types = append(quantitiy_types, r)
        }
        if err := rows.Err(); err != nil {
			log.Fatal(err)
        }
		return quantitiy_types
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

func CreateIngredient(w http.ResponseWriter, r *http.Request, db *sql.DB) {
			// Retrieve the form data
			ingredient := r.FormValue("ingredient")
			ingredientType := r.FormValue("ingredient_type")
		
		
		// Perform the SQL INSERT query to add the ingredient to the database
		_, err := db.Exec("INSERT INTO ingredients (name, ingredient_type_id) VALUES (?, ?)", ingredient, ingredientType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		// Redirect back to the home page
		fmt.Fprintf(w, `<script>window.location.href = "/create-recipe";</script>`)
}
func UpdateIngredient(w http.ResponseWriter, r *http.Request, db *sql.DB) {
		ingredientID := r.FormValue("id")
		updatedName := r.FormValue("ingredientName")
	
		// Prepare and execute the SQL UPDATE statement
		_, err := db.Exec("UPDATE ingredients SET name = ? WHERE ingredient_id = ?", updatedName, ingredientID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		return
}
func DeleteIngredient(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
}
func getIngredientsForRecipe( db *sql.DB, recipeId string) []IngredientWithQuantity {
	rows, err := db.Query(`
		SELECT i.ingredient_id, i.name, ri.quantity, qt.name AS quantity_type, qt.quantity_type_id
		FROM ingredients i
		JOIN recipe_ingredients ri ON ri.ingredient_id = i.ingredient_id
		JOIN quantity_type qt ON qt.quantity_type_id = ri.quantity_type_id
		WHERE ri.recipe_id = ?
	`, recipeId)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var ingredients []IngredientWithQuantity

	for rows.Next() {
		var ingredient IngredientWithQuantity

		err := rows.Scan(&ingredient.Ingredient_id, &ingredient.Name, &ingredient.Quantity, &ingredient.Quantity_type, &ingredient.Quantity_type_id)
		if err != nil {
			panic(err.Error())
		}

		ingredients = append(ingredients, ingredient)
	}
	return ingredients
}
func getIngredientsForRecipeWithType( db *sql.DB, recipeId string) []IngredientWithQuantityAndType {
	rows, err := db.Query(`
		SELECT i.ingredient_id, i.name, ri.quantity, qt.name AS quantity_type, qt.quantity_type_id, i.ingredient_type_id, it.name AS ingredient_type_name
		FROM ingredients i
		JOIN recipe_ingredients ri ON ri.ingredient_id = i.ingredient_id
		JOIN quantity_type qt ON qt.quantity_type_id = ri.quantity_type_id
		JOIN ingredient_type it ON it.ingredient_type_id = i.ingredient_type_id
		WHERE ri.recipe_id = ?
	`, recipeId)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var ingredients []IngredientWithQuantityAndType

	for rows.Next() {
		var ingredient IngredientWithQuantityAndType

		err := rows.Scan(&ingredient.Ingredient_id, &ingredient.Name, &ingredient.Quantity, &ingredient.Quantity_type, &ingredient.Quantity_type_id, &ingredient.Ingredient_Type_id, &ingredient.Ingredient_Type_Name)
		if err != nil {
			panic(err.Error())
		}

		ingredients = append(ingredients, ingredient)
	}
	return ingredients
}