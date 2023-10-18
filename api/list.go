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

type IngredientQuantityData struct {
	IngredientName string
    Quantity float32
    QuantityTypeName string
    IngredientId int
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
	
		fmt.Println("", recipeName)
	ingredientData := GetIngredientQuantityDataFromRecipeIds(recipeIds, db)

	for _, data := range ingredientData {
		ingredients = append(ingredients, Ingredient{
			Name: data,
		})
	}
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
// Handle the HTMX GET request for updating ingredients
func UpdateIngredientsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	recipeName := r.FormValue("recipeName")
	recipeIds := r.Form["recipes"]
	fmt.Println("", recipeName)

    // Implement logic to retrieve updated ingredients based on the selected recipe
    // You should generate an HTML list structure here
	// ingredients:= getIngredientsForRecipe(db, recipeID)

	ingredientData := GetIngredientQuantityDataFromRecipeIds(recipeIds, db)


    // Generate an HTML ul and li structure
    ul := "<ul>"
    for _, ingredient := range ingredientData {
        ul += "<li>" + ingredient + "</li>"
    }
    ul += "</ul"

    // Send the updated HTML ingredient list as a response
    w.Header().Set("Content-Type", "text/html") // Set the content type to HTML
    w.Write([]byte(ul)) // Write the HTML structure to the response
}


func GetIngredientQuantityDataFromRecipeIds( recipeIds []string, db *sql.DB) []string {
	var response []string
var ingredientQuantityData []IngredientQuantityData
		// Retrieve the selected ingredients

// Log the selected ingredient IDs
		// Iterate through the selected recipe IDs
		for _, recipeID := range recipeIds {
			// Convert the recipeID string to an integer
			recipeIDInt, err := strconv.Atoi(recipeID)
			if err != nil {
				// http.Error(w, err.Error(), http.StatusBadRequest)
				return response
			}
	
	rows, err := db.Query(`
    SELECT i.name, ri.quantity, qt.name, i.ingredient_id
    FROM ingredients i
    INNER JOIN recipe_ingredients ri ON i.ingredient_id = ri.ingredient_id
    INNER JOIN quantity_type qt ON ri.quantity_type_id = qt.quantity_type_id
    WHERE ri.recipe_id = ?
`, recipeIDInt)

if err != nil {
    // http.Error(w, err.Error(), http.StatusInternalServerError)
    return response
}
defer rows.Close()



// Loop through the rows of ingredients and append them to the list
for rows.Next() {
    var ingredientName string
    var quantity float32
    var quantityTypeName string
    var ingredientId int
    err := rows.Scan(&ingredientName, &quantity, &quantityTypeName, &ingredientId)
    if err != nil {
        // http.Error(w, err.Error(), http.StatusInternalServerError)
        return response
    }
// Create a flag to track if the ingredient with the same ingredientID exists
ingredientExists := false

// Iterate through ingredientQuantityData to find a match
for i, data := range ingredientQuantityData {
    if data.IngredientId == ingredientId {
        // Ingredient with the same ID exists, update the quantity

        ingredientQuantityData[i].Quantity =  quantity + ingredientQuantityData[i].Quantity
        ingredientExists = true
    }
}
if ingredientExists == false {
    ingredientQuantityData = append(ingredientQuantityData, IngredientQuantityData{
        IngredientName:    ingredientName,
        Quantity:          quantity,
        QuantityTypeName:  quantityTypeName,
        IngredientId:      ingredientId,
    })
}
}

}
for _, data := range ingredientQuantityData {
    // Convert the float32 to a string with a specific format
    stringValue := strconv.FormatFloat(float64(data.Quantity), 'f', -1, 32)

    // Create the formatted ingredient name
    ingredientDetails := data.IngredientName + " " + stringValue + " " + data.QuantityTypeName

    // Append the ingredient to the ingredients list
    response = append(response, ingredientDetails)
}
return response
}