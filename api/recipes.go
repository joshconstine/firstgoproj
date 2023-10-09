package api

import (
    "net/http"
    "fmt"
    "database/sql"
	"html/template"
    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux" 
	 "github.com/aws/aws-sdk-go/aws"

    "github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"io"
	"log"
	"io/ioutil"
	"path/filepath"
	"bytes"
	"github.com/google/uuid"
	"context"
    "mime/multipart"
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
type RecipeWithIngredients struct {
	Recipe_id int
	Name  string
	Description string
	Ingredients []Ingredient
}
type SingleRecipePageData struct {
	PageTitle string
    Recipe RecipeWithIngredients
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
//S3 transaction
func fileToBytes(file multipart.File) ([]byte, error) {
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return fileBytes, nil
}



func UploadHandler(w http.ResponseWriter, r *http.Request, uploader *s3manager.Uploader) (string, error) {
	err := r.ParseMultipartForm(10 * 1024 * 1024) // 10 MB limit
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusInternalServerError)
		return "", err
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusInternalServerError)
		return "", err
	}
	defer file.Close()

	// Generate a unique filename using a UUID
	fileExt := filepath.Ext(header.Filename)
	newFilename := uuid.New().String() + fileExt

	// Create a new file in the "public/static" directory with the unique filename
	newFilePath := filepath.Join("public/static/images", newFilename)
	newFile, err := os.Create(newFilePath)
	if err != nil {
		http.Error(w, "Failed to create a new file", http.StatusInternalServerError)
		return "nil", err
	}
	defer newFile.Close()

	// Reset the file pointer to the beginning before copying
	// Copy the uploaded file to the new file
	_, err = io.Copy(newFile, file)
	if err != nil {
		http.Error(w, "Failed to copy the file", http.StatusInternalServerError)
		return "nil", err
	}
	log.Println("uploading so S3")

	// file, err := ioutil.ReadFile(newFilePath)
	// if err != nil {
		// 	log.Fatal(err)
		// }
		
		
		BUCKET_NAME := "foodly-bucket"
		// BUCKET_URL := "https://foodly-bucket.s3.us-west-1.amazonaws.com"
		// NEXT_PUBLIC_BUCKET_URL := "https://foodly-bucket.s3.us-west-1.amazonaws.com"

		
		file.Seek(0, 0)
fileBytes, err := fileToBytes(file)
if err != nil {
    // Handle the error
		http.Error(w, "Failed to read photo to bytes", http.StatusInternalServerError)
		return "" , err
}



	upInput := &s3manager.UploadInput{
		Bucket:      aws.String(BUCKET_NAME), // bucket's name
		Key:         aws.String(newFilePath),        // files destination location
		Body:        bytes.NewReader(fileBytes),                   // content of the file
		ContentType: aws.String(fileExt),                 // content type
	}
	res, err := uploader.UploadWithContext(context.Background(), upInput)
	log.Printf("res %+v\n", res)
	log.Printf("err %+v\n", err)
	 
	 createdFileLocation := res.Location

	log.Printf("Create file location %+v\n", createdFileLocation)
	// Respond with the unique filename or other relevant information
	// fmt.Fprintf(w, "File uploaded successfully with filename: %s", newFilename)
	return createdFileLocation, nil
}
// DB Transactions
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

    uploader := NewUploader()
	newPhotoLocation, err := UploadHandler(w, r, uploader)
	if err == nil {
		_, err = db.Exec("INSERT INTO recipe_photos (recipe_id, photo_url) VALUES (?, ?)", recipeID, newPhotoLocation)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}



	log.Printf("new location in create recipe file %+v\n", newPhotoLocation)
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

func GetRecipeById(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
}

