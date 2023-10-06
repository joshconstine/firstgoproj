package api

import (
	"database/sql"
	"log"
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
type IngredientAndType struct {
	Ingredient_id int
	Name  string
	Ingredient_type_id int
	Ingredient_type_name string

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