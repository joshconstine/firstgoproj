package main

import (
	"fmt"
	"strconv"
	"net/http"
	
    "github.com/gorilla/mux"
    "firstgoprog/api" // Replace "firstgoprog" with your actual module name.
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
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

type Ingredient struct {
	Ingredient_id int
	Name  string
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
// func seed()  {	
// 	db, err := sql.Open("mysql", "root:daddy@(127.0.0.1:3306)/food?parseTime=true")
//     if err != nil {
// 		// log.Fatal(err)
// 		fmt.Print("error connecting to db")
//     }
//     if err := db.Ping(); err != nil {
// 		fmt.Printf("Error %d...\n", err)
//     }
// 	query := `CREATE TABLE ingredients (
// 		ingredient_id INT AUTO_INCREMENT PRIMARY KEY,
// 		name VARCHAR(255) NOT NULL
// 	);

// 	CREATE TABLE recipes (
// 		recipe_id INT AUTO_INCREMENT PRIMARY KEY,
// 		name VARCHAR(255) NOT NULL,
// 		description TEXT
// 	);

// 	CREATE TABLE recipe_ingredients (
// 		recipe_id INT,
// 		ingredient_id INT,
// 		PRIMARY KEY (recipe_id, ingredient_id),
// 		FOREIGN KEY (recipe_id) REFERENCES recipes(recipe_id),
// 		FOREIGN KEY (ingredient_id) REFERENCES ingredients(ingredient_id)
// 	);`

// _, erro := db.Exec(query)
// if erro != nil {
// 	return 
// 	fmt.Printf("err %s", erro)
// }

// fmt.Println("Tables created successfully")
// return 
// }


func getAllIngredients(db *sql.DB)  []Ingredient {
	rows, err := db.Query(`SELECT * FROM ingredients`)
        if err != nil {
			log.Fatal(err)
        }
        defer rows.Close()
		
        var ingredients []Ingredient
        for rows.Next() {
			var i Ingredient
			
            err := rows.Scan(&i.Ingredient_id, &i.Name)
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
	
	db, err := sql.Open("mysql", "root:daddy@(127.0.0.1:3306)/food?parseTime=true")
    if err != nil {
		// log.Fatal(err)
		fmt.Print("error connecting to db")
    }
    if err := db.Ping(); err != nil {
		fmt.Printf("Error %d...\n", err)
    }
    // Use the functions from the 'api' package to define routes.
	api.InitRoutes(r)
	
	
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("public/layout.html"))
	
        ingredients := getAllIngredients(db)
      
		
		data := IngredientPageData{
			PageTitle: "My ingredients list",
            Ingredients: ingredients,
        }

        tmpl.Execute(w, data)
    })
	r.HandleFunc("/add-ingredient", func(w http.ResponseWriter, r *http.Request) {
			// Retrieve the form data
			ingredient := r.FormValue("ingredient")
		
		
		// Perform the SQL INSERT query to add the ingredient to the database
		_, err = db.Exec("INSERT INTO ingredients (name) VALUES (?)", ingredient)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		// Redirect back to the home page
		fmt.Fprintf(w, `<script>window.location.href = "/";</script>`)
	})
	r.HandleFunc("/delete-ingredient", func(w http.ResponseWriter, r *http.Request) {
		
		id := r.FormValue("id")
		// Perform the SQL INSERT query to add the ingredient to the database
		stmt, err := db.Prepare("DELETE FROM ingredients WHERE ingredient_id = ?")
    if err != nil {
        // return err
    }
    defer stmt.Close()

    // Execute the SQL statement
    _, err = stmt.Exec(id)
    if err != nil {
        // return err
    }

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

//RECIPES	
r.HandleFunc("/recipes", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("public/recipes.html"))
	
		recipes := getAllRecipes(db)
        ingredients := getAllIngredients(db)
	
		data := RecipesPageData{
			PageTitle: "Recipes",
            Recipes: recipes,
            Ingredients: ingredients,
        }

        tmpl.Execute(w, data)
	})	
	r.HandleFunc("/add-recipe", func(w http.ResponseWriter, r *http.Request) {
			// Retrieve the form data
			recipeName := r.FormValue("recipeName")
			recipeDescription := r.FormValue("recipeDescription")
			// Retrieve the selected ingredients
			ingredientIDs := r.Form["ingredients"]



		// Perform the SQL INSERT query to add the ingredient to the database
		_, err = db.Exec("INSERT INTO recipes (name, description) VALUES (?, ?)", recipeName, recipeDescription )
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
		fmt.Fprintf(w, `<script>window.location.href = "/recipes";</script>`)
	})
r.HandleFunc("/delete-recipe", func(w http.ResponseWriter, r *http.Request) {
		
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
	

	fmt.Printf("Server is listening on port %d...\n", port)
	
	http.ListenAndServe(":"+portStr, r)
}