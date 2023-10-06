package main

import (
	"fmt"
	"os"
	"strconv"
	"net/http"
    "github.com/gorilla/mux"
    "firstgoprog/api" // Replace "firstgoprog" with your actual module name.
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/joho/godotenv/autoload"
	"log"
	"html/template"
	"github.com/twilio/twilio-go"
	twapi "github.com/twilio/twilio-go/rest/api/v2010"
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
type RecipeWithIngredients struct {
	Recipe_id int
	Name  string
	Description string
	Ingredients []Ingredient
}
type IngredientPageData struct {
	PageTitle string
    Ingredients []Ingredient
	IngredientTypes []IngredientType
	MappedIngredients map[string][]IngredientAndType
}
type RecipesPageData struct {
	PageTitle string
    Recipes []Recipe
    Ingredients []Ingredient
}

type SingleRecipePageData struct {
	PageTitle string
    Recipe RecipeWithIngredients
}
type CreateListPageData struct {
	PageTitle string
    Recipes []Recipe
}
type ListPageData struct {
	PageTitle string
    Ingredients []Ingredient
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
func getSingleRecipeWithIngredients(db *sql.DB, id string) (RecipeWithIngredients, error) {
	 // Define a variable to hold the result
	 var result RecipeWithIngredients

	 // Query the recipe information based on the provided id
	 err := db.QueryRow("SELECT name, description, recipe_id FROM recipes WHERE recipe_id = ?", id).
		 Scan(&result.Name, &result.Description, &result.Recipe_id)
	 if err != nil {
		 return result, err
	 }
 
	 // Query the associated ingredients for the recipe
	 rows, err := db.Query("SELECT i.name FROM ingredients i INNER JOIN recipe_ingredients ri ON i.ingredient_id = ri.ingredient_id WHERE ri.recipe_id = ?", id)
	 if err != nil {
		 return result, err
	 }
	 defer rows.Close()
 
	 // Loop through the rows of ingredients and add them to the result
	 for rows.Next() {
		 var ingredientName string
		 err := rows.Scan(&ingredientName)
		 if err != nil {
			 return result, err
		 }
		 result.Ingredients = append(result.Ingredients, Ingredient{Name: ingredientName})
	 }
 
	 // Check for errors during rows iteration
	 if err := rows.Err(); err != nil {
		 return result, err
	 }
	 return result, nil
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

 

	r.HandleFunc("/recipes/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
        id := vars["id"]
			tmpl := template.Must(template.ParseFiles("public/singleRecipe.html"))
	
		recipe, err := getSingleRecipeWithIngredients(db, id)
		if err != nil {
			http.Error(w, "Unable to read from db", http.StatusInternalServerError)
		}		


		data := SingleRecipePageData{
			PageTitle: recipe.Name,
            Recipe: recipe,
            
        }

        tmpl.Execute(w, data)
	})
//List

r.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("public/makeList.html"))
	
		recipes := getAllRecipes(db)
	
		data := CreateListPageData{
			PageTitle: "Make a List",
            Recipes: recipes,
        }

        tmpl.Execute(w, data)
	})	
	
	r.HandleFunc("/generate-list", func(w http.ResponseWriter, r *http.Request) {
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
	})
	r.HandleFunc("/send-list", func(w http.ResponseWriter, r *http.Request) {
			// Retrieve the form data
			phoneNumber := r.FormValue("phone")
			list := r.FormValue("list")
			accountSid :=	os.Getenv("TWILIO_ACCOUNT_SID")
			authToken := os.Getenv("TWILIO_AUTH_TOKEN")
			fullPhoneNumber := "+1" + phoneNumber
			client := twilio.NewRestClientWithParams(twilio.ClientParams{
				Username: accountSid,
				Password: authToken,
			})
			

			params := &twapi.CreateMessageParams{}
			params.SetFrom("+18888415616")
			params.SetBody(list)
			params.SetTo(fullPhoneNumber)
		
			resp, err := client.Api.CreateMessage(params)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				if resp.Sid != nil {
					// fmt.Println(*resp.Sid)
				} else {
					// fmt.Println(resp.Sid)
				}
			}


		// Redirect back to the home page
    fmt.Fprintf(w, "List send to : %+v\n", phoneNumber)
	})

	fmt.Printf("Server is listening on port %d...\n", port)
	
	http.ListenAndServe(":"+portStr, r)
}