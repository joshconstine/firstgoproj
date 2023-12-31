package api

import (
    "net/http"
    "fmt"
    "database/sql"
	"strings"
	"errors"
	"html/template"
    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux" 
	 "github.com/aws/aws-sdk-go/aws"
	"strconv"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"io/ioutil"
	"path/filepath"
	"bytes"
	"github.com/google/uuid"
	"context"
	"github.com/srinathgs/mysqlstore"
    "mime/multipart"
)

type Recipe struct {
	Recipe_id int
	Name  string
	Description string
}
type CreateRecipePageData struct {
	PageTitle string
    Ingredients []Ingredient
	MappedIngredients map[string][]IngredientAndType
	Tags []Tag
	User User
	IngredientTypes []IngredientType
}
type RecipeWithIngredients struct {
	Recipe_id int
	Name  string
	Description string
	Ingredients []Ingredient
}
type RecipeWithIngredientsAndPhotos struct {
	Recipe_id int
	Name  string
	Description string
	Ingredients []IngredientWithQuantity
	Photos []string
}
type RecipeWithIngredientsAndPhotosAndTags struct {
	Recipe_id int
	Name  string
	Description string
	Ingredients []IngredientWithQuantity
	Photos []string
	Tags []Tag
}
type RecipeWithPhotos struct {
	Recipe_id int
	Name  string
	Description string
	Photos []string
}
type RecipeWithPhotosAndTags struct {
	Recipe_id int
	Name  string
	Description string
	Photos []string
	Tags []Tag
	TagString string
	Ingredients []IngredientWithQuantity
}


type SingleRecipePageData struct {
	PageTitle string
    Recipe RecipeWithIngredientsAndPhotosAndTags
    QuantityTypes []QuantityType
	User User
}

//HTML TEMPLATES
func GetRecipeById(w http.ResponseWriter, r *http.Request, db *sql.DB, store *mysqlstore.MySQLStore) {
		vars := mux.Vars(r)
        id := vars["id"]
	
		tmpl := template.Must(template.ParseFiles("public/singleRecipe.html"))

		quantitiy_types := getAllQuantityTypes(db)
		recipe, err := getSingleRecipeWithIngredientsAndPhotosAndTags(db, id)
		if err != nil {
			http.Error(w, "Unable to read from db", http.StatusInternalServerError)
		}		
		user, err := GetUserFromRequest(w, r, db, store)
		if err != nil {
			fmt.Println(err)
		}

		data := SingleRecipePageData{
			PageTitle: recipe.Name,
            Recipe: recipe,
			QuantityTypes: quantitiy_types,
			User: user,
            
        }

        tmpl.Execute(w, data)
}


func GetCreateRecipeTemplate(w http.ResponseWriter, r *http.Request, db *sql.DB, store *mysqlstore.MySQLStore) {
 		tmpl := template.Must(template.ParseFiles("public/createRecipe.html"))
	
        ingredients := getAllIngredients(db)
		ingredientTypeMap := getAllIngredientsWithTypes(db)
        tags := getAllTags(db)
		user, err := GetUserFromRequest(w, r, db, store)
		if err != nil {
			fmt.Println(err)
		}
		ingredientTypes := getAllIngredientTypes(db)
		data := CreateRecipePageData{
			PageTitle: "Create Recipe",
            Ingredients: ingredients,
			MappedIngredients: ingredientTypeMap,
			Tags: tags,
			User: user,
			IngredientTypes: ingredientTypes,
        }

        tmpl.Execute(w, data)
}
//S3 transaction
//turns type multipart.File into a byte array
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
	newFilePath := filepath.Join("listify/recipes", newFilename)
	// newFile, err := os.Create(newFilePath)
	// if err != nil {
	// 	http.Error(w, "Failed to create a new file", http.StatusInternalServerError)
	// 	return "nil", err
	// }
	// defer newFile.Close()

	// // Reset the file pointer to the beginning before copying
	// // Copy the uploaded file to the new file
	// _, err = io.Copy(newFile, file)
	// if err != nil {
	// 	http.Error(w, "Failed to copy the file", http.StatusInternalServerError)
	// 	return "nil", err
	// }
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
func getSingleRecipeWithIngredientsAndPhotosAndTags(db *sql.DB, id string) (RecipeWithIngredientsAndPhotosAndTags, error) {
	// Define a variable to hold the result
	var result RecipeWithIngredientsAndPhotosAndTags

	
	tags := getTagsforRecipeId(db, id)
	result.Tags = tags
	// Query the recipe information based on the provided id
	err := db.QueryRow("SELECT name, description, recipe_id FROM recipes WHERE recipe_id = ?", id).
		Scan(&result.Name, &result.Description, &result.Recipe_id)
	if err != nil {
		return result, err
	}

	// Query the associated ingredients with quantity for the recipe
	rows, err := db.Query("SELECT i.ingredient_id, i.name, ri.quantity, ri.quantity_type_id, qt.name FROM ingredients i INNER JOIN recipe_ingredients ri ON i.ingredient_id = ri.ingredient_id INNER JOIN quantity_type qt ON ri.quantity_type_id = qt.quantity_type_id WHERE ri.recipe_id = ?", id)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	// Loop through the rows of ingredients and add them to the result
	for rows.Next() {
		var ingredientWithQuantity IngredientWithQuantity
		err := rows.Scan(&ingredientWithQuantity.Ingredient_id, &ingredientWithQuantity.Name, &ingredientWithQuantity.Quantity, &ingredientWithQuantity.Quantity_type_id, &ingredientWithQuantity.Quantity_type)
		if err != nil {
			return result, err
		}
		// In your Go code, set the Quantity_type value based on Quantity
if ingredientWithQuantity.Quantity > 1 {
    ingredientWithQuantity.Quantity_type = ingredientWithQuantity.Quantity_type + "s"
}
		result.Ingredients = append(result.Ingredients, ingredientWithQuantity)
	}

// Check for errors during rows iteration
if err := rows.Err(); err != nil {
	return result, err
}

	// Query the associated photos for the recipe
	rows, err = db.Query("SELECT photo_url FROM recipe_photos  WHERE recipe_id = ?", id)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	// Loop through the rows of ingredients and add them to the result
	for rows.Next() {
		var photoUrl string
		err := rows.Scan(&photoUrl)
		if err != nil {
			return result, err
		}
		result.Photos = append(result.Photos, photoUrl)
	}
	
	// Check for errors during rows iteration
	if err := rows.Err(); err != nil {
		return result, err
	}
	return result, nil
}
func getAllRecipesWithPhotos(db *sql.DB) ([]RecipeWithPhotos, error) {
	// Define a variable to hold the result
	 recipes := getAllRecipes(db)
	var result []RecipeWithPhotos

for _, recipe := range recipes {
    // Create a RecipeWithPhotos instance for the current recipe
    recipeWithPhotos := RecipeWithPhotos{
        Recipe_id: recipe.Recipe_id,
	Name:recipe.Name,
	Description:recipe.Description,
    }

    // Query the associated photos for the current recipe
    rows, err := db.Query("SELECT photo_url FROM recipe_photos WHERE recipe_id = ?", recipe.Recipe_id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Loop through the rows of photos and add them to the result for the current recipe
    for rows.Next() {
        var photoUrl string
        if err := rows.Scan(&photoUrl); err != nil {
            return nil, err
        }
        recipeWithPhotos.Photos = append(recipeWithPhotos.Photos, photoUrl)
    }

    // Check for errors during rows iteration
    if err := rows.Err(); err != nil {
        return nil, err
    }

    // Append the current recipe with its photos to the final result
    result = append(result, recipeWithPhotos)
}
	return result, nil
}
func getJoinedTags(tags []Tag) string {
    var tagStrings []string
    for _, tag := range tags {
        tagStrings = append(tagStrings, tag.Name)
    }
    return strings.Join(tagStrings, ", ")
}


func getAllRecipesWithPhotosAndTags(db *sql.DB) ([]RecipeWithPhotosAndTags, error) {
	// Define a variable to hold the result
	 recipes := getAllRecipes(db)
	var result []RecipeWithPhotosAndTags

for _, recipe := range recipes {
    // Create a RecipeWithPhotos instance for the current recipe
    recipeWithPhotos := RecipeWithPhotosAndTags{
        Recipe_id: recipe.Recipe_id,
	Name:recipe.Name,
	Description:recipe.Description,
    }


recipeID := strconv.Itoa(recipe.Recipe_id) // Convert int to string

tags := getTagsforRecipeId(db,recipeID) // Now, you can pass the string value
recipeWithPhotos.Tags = tags
tagString := getJoinedTags(tags)
recipeWithPhotos.TagString =  tagString

ingredientWithQuantity := getIngredientsForRecipe(db, recipeID)

recipeWithPhotos.Ingredients = ingredientWithQuantity



    // Query the associated photos for the current recipe
    rows, err := db.Query("SELECT photo_url FROM recipe_photos WHERE recipe_id = ?", recipe.Recipe_id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Loop through the rows of photos and add them to the result for the current recipe
    for rows.Next() {
        var photoUrl string
        if err := rows.Scan(&photoUrl); err != nil {
            return nil, err
        }
        recipeWithPhotos.Photos = append(recipeWithPhotos.Photos, photoUrl)
    }

    // Check for errors during rows iteration
    if err := rows.Err(); err != nil {
        return nil, err
    }

    // Append the current recipe with its photos to the final result
    result = append(result, recipeWithPhotos)
}
	return result, nil
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
			tagIDs := r.Form["tags"]
			quantityValue := 1
			quantityTypeValue := 1

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
	_, err = db.Exec("INSERT INTO recipe_ingredients (recipe_id, ingredient_id, quantity, quantity_type_id) VALUES (?, ?, ?, ?)", recipeID, ingredientID, quantityValue, quantityTypeValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
 for _, tagID := range tagIDs {
	_, err = db.Exec("INSERT INTO recipe_tags (recipe_id, tag_id) VALUES (?, ?)", recipeID, tagID)
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
func UploadNewPhoto(w http.ResponseWriter, r *http.Request, db *sql.DB) {
			recipeIDStr := r.FormValue("id")

    uploader := NewUploader()
	newPhotoLocation, err := UploadHandler(w, r, uploader)
	if err == nil {
		_, err = db.Exec("INSERT INTO recipe_photos (recipe_id, photo_url) VALUES (?, ?)", recipeIDStr, newPhotoLocation)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
		recipeID, err := strconv.Atoi(recipeIDStr) // Convert the string to an integer
		if err != nil {
			// Handle the error if the conversion fails
			http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
			return
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

		fmt.Fprintf(w, `<script>window.location.href = "/list";</script>`)
}

func UpdateRecipeIngredients(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    if r.Method == http.MethodPost {
        err := r.ParseForm()
        if err != nil {
            http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
            return
        }

        recipeID := r.FormValue("recipe_id")
        if recipeID == "" {
            http.Error(w, "Missing recipe_id in form data", http.StatusBadRequest)
            return
        }

        for key, values := range r.Form {
            // Check if the form data key represents an ingredient type or quantity
            if strings.HasSuffix(key, "_type") {
                ingredientID := strings.TrimSuffix(key, "_type")
                ingredientType := values[0]

                // Convert ingredientID and recipeID to integers
                ingredientIDInt, err := strconv.Atoi(ingredientID)
                if err != nil {
                    http.Error(w, "Invalid ingredient ID", http.StatusBadRequest)
                    return
                }

                recipeIDInt, err := strconv.Atoi(recipeID)
                if err != nil {
                    http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
                    return
                }

                // Update the database with the new ingredient type
                // Example SQL query using the database/sql package:
                // Assuming you have a db variable of type *sql.DB
                _, err = db.Exec("UPDATE recipe_ingredients SET quantity_type_id = ? WHERE ingredient_id = ? AND recipe_id = ?", ingredientType, ingredientIDInt, recipeIDInt)
                if err != nil {
                    http.Error(w, "Failed to update ingredient type", http.StatusInternalServerError)
                    return
                }
            } else if strings.HasSuffix(key, "_quantity") {
                ingredientID := strings.TrimSuffix(key, "_quantity")
                ingredientQuantity := values[0]

                // Convert ingredientID and recipeID to integers
                ingredientIDInt, err := strconv.Atoi(ingredientID)
                if err != nil {
                    http.Error(w, "Invalid ingredient ID", http.StatusBadRequest)
                    return
                }

                recipeIDInt, err := strconv.Atoi(recipeID)
                if err != nil {
                    http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
                    return
                }

                // Update the database with the new ingredient quantity
                // Example SQL query using the database/sql package:
                // Assuming you have a db variable of type *sql.DB
                _, err = db.Exec("UPDATE recipe_ingredients SET quantity = ? WHERE ingredient_id = ? AND recipe_id = ?", ingredientQuantity, ingredientIDInt, recipeIDInt)
                if err != nil {
                    http.Error(w, "Failed to update ingredient quantity", http.StatusInternalServerError)
                    return
                }
            }
        }

        
		redirectURL := fmt.Sprintf("/recipes/%s", recipeID)

		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
    }
}
func UpdateDescription(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	
			r.ParseForm()
			
			// Get the recipe ID from the form data
			recipeID := r.FormValue("id")
			
			// Get the new description from the form data
			newDescription := r.FormValue("description")
			
			// Update the description in the database using SQL or an ORM
			// Example SQL query: "UPDATE recipes SET description = ? WHERE recipe_id = ?"
			
			_, err := db.Exec("UPDATE recipes SET description = ? WHERE recipe_id = ?", newDescription, recipeID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}


		redirectURL := fmt.Sprintf("/recipes/%s", recipeID)
			// Redirect to the updated recipe's page or display a success message
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	
		http.ListenAndServe(":8080", nil)
	}
	func ToggleUserFavoriteRecipe(w http.ResponseWriter, r *http.Request, db *sql.DB, store *mysqlstore.MySQLStore) {
		// Get the recipe ID from the form data
		recipeID := r.FormValue("id")
		// Get the user ID from the session
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
            fmt.Printf("session_token not found")
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

    if sessionToken == "" {
        http.Error(w, "Unauthorized, please sign in to view this page", http.StatusUnauthorized)
        return
    }

	// We then get the session from our session map
	session, err := store.Get(r, sessionToken)
    if err != nil {
        switch {
        case errors.Is(err, http.ErrNoCookie):
            http.Error(w, "cookie not found", http.StatusBadRequest)
        default:
            log.Println(err)
            http.Error(w, "server error", http.StatusInternalServerError)
        }
        return
    }



		
		userID := session.Values["user_id"]
		// Check if the user has already favorited the recipe
// Perform a query to check if the user has favorited the recipe
row := db.QueryRow("SELECT user_id FROM user_favorite_recipes WHERE user_id = ? AND recipe_id = ?", userID, recipeID)
var userExists int
err = row.Scan(&userExists)
if err == sql.ErrNoRows {
	// User has not favorited the recipe yet, so add it to the database
	_, err = db.Exec("INSERT INTO user_favorite_recipes (user_id, recipe_id) VALUES (?, ?)", userID, recipeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
} else if err != nil {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
} else {
	// User has already favorited the recipe, so remove it from the database
	_, err = db.Exec("DELETE FROM user_favorite_recipes WHERE user_id = ? AND recipe_id = ?", userID, recipeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

	container := "<div  id=\"FavoritedContainer\" >"
			container += `<h1 class="text-m"> Favorited </h1>`
			container += "<button class='btn btn-ghost'> <a href='/profile'> Show My Favorites </a> </button>"
		container += "</div>"

    // Send the updated HTML ingredient list as a response
    w.Header().Set("Content-Type", "text/html") // Set the content type to HTML
    w.Write([]byte(container)) // Write the HTML structure to the response
	}