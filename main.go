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
type IngredientPageData struct {
	PageTitle string
    Ingredients []Ingredient
}
type RecipesPageData struct {
	PageTitle string
    Recipes []Recipe
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
		
		data := IngredientPageData{
			PageTitle: "My ingredients list",
            Ingredients: ingredients,
        }

        tmpl.Execute(w, data)
    })
	// r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 	w.Header().Set("Content-Type", "text/html")
		// 	htmlFile, err := os.Open("public/index.html")
		// 	if err != nil {
			// 		http.Error(w, "Unable to read HTML file", http.StatusInternalServerError)
			// 		return
			// 	}
			// 	defer htmlFile.Close()
			
			// 	// Copy the HTML content to the response writer.
			// 	_, err = io.Copy(w, htmlFile)
			// 	if err != nil {
				// 		http.Error(w, "Unable to copy HTML content to response", http.StatusInternalServerError)
	// 		return
	// 	}
	// })
	
	r.HandleFunc("/recipes", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("public/recipes.html"))
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
		
		data := RecipesPageData{
			PageTitle: "Recipes",
            Recipes: recipes,
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
	r.HandleFunc("/add-recipe", func(w http.ResponseWriter, r *http.Request) {
			// Retrieve the form data
			recipeName := r.FormValue("recipeName")
			recipeDescription := r.FormValue("recipeDescription")
		
		
		// Perform the SQL INSERT query to add the ingredient to the database
		_, err = db.Exec("INSERT INTO recipes (name, description) VALUES (?, ?)", recipeName, recipeDescription )
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
		fmt.Fprintf(w, "You've requested the recipe: %s\n", id)

		// w.Header().Set("Content-Type", "text/html")
		// htmlFile, err := os.Open("public/recipes.html")
		// if err != nil {
		// 	http.Error(w, "Unable to read HTML file", http.StatusInternalServerError)
		// 	return
		// }
		// defer htmlFile.Close()

		// // Copy the HTML content to the response writer.
		// _, err = io.Copy(w, htmlFile)
		// if err != nil {
		// 	http.Error(w, "Unable to copy HTML content to response", http.StatusInternalServerError)
		// 	return
		// }
	})

	

	fmt.Printf("Server is listening on port %d...\n", port)
	
	http.ListenAndServe(":"+portStr, r)
}