package api

import (
    "net/http"
    "database/sql"
	"html/template"
    _ "github.com/go-sql-driver/mysql"
)
type CreateListPageData struct {
	PageTitle string
    Recipes []Recipe
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

