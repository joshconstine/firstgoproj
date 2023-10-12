package api

import (
	"os"
	"fmt"
	"strconv"
    "net/http"
    "database/sql"
	"html/template"
    _ "github.com/go-sql-driver/mysql"
	"github.com/twilio/twilio-go"
	twapi "github.com/twilio/twilio-go/rest/api/v2010"
)
type CreateListPageData struct {
	PageTitle string
    Recipes []RecipeWithPhotos
}
type ListPageData struct {
	PageTitle string
    Ingredients []Ingredient
}

//HTML TEMPLATES

func GetListTemplate(w http.ResponseWriter, r *http.Request, db *sql.DB) {
  	tmpl := template.Must(template.ParseFiles("public/makeList.html"))
	
		recipes,_ := getAllRecipesWithPhotos(db)
	
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
			rows, err := db.Query("SELECT i.name, ri.quantity FROM ingredients i INNER JOIN recipe_ingredients ri ON i.ingredient_id = ri.ingredient_id WHERE ri.recipe_id = ?", recipeIDInt)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()
	
			// Loop through the rows of ingredients and append them to the list
			for rows.Next() {
				var ingredientName string
				var quantity float32
				err := rows.Scan(&ingredientName, &quantity)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				 // Convert the float32 to a string with a specific format
				 stringValue := strconv.FormatFloat(float64(quantity), 'f', -1, 32)

				ingredients = append(ingredients, Ingredient{Name: ingredientName + " " + stringValue})
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
	

	//DB transactions

func SendList(w http.ResponseWriter, r *http.Request) {
			// Retrieve the form data
			phoneNumber := r.FormValue("phone")
			list := r.FormValue("list")
			accountSid :=os.Getenv("TWILIO_ACCOUNT_SID")	
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
}