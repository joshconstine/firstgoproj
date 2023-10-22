
package api

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"net/http"	
)
// User struct represents a user's data
type User struct {
    ID       int
    Username string
}



// HashPassword hashes a plain text password
func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}



func HandleInsertUser(w http.ResponseWriter, r *http.Request, db *sql.DB)  {
	// Retrieve the form data
	username := r.FormValue("username")
	password := r.FormValue("password")
	
	// Hash the password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Perform the SQL INSERT query to add the user to the database
	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	container := "<div  id=\"successContainer\" data-hx-target=\"ingredientList\" class=\"block w-full rounded-lg p-3 flex h-full justify-center max-h-full flex-col items-center \" >"
			container += `<h1 class="text-m"> User added to database </h1> <h1 class="text-m">` + username + `</h1>`
		container += "</div>"
    w.Header().Set("Content-Type", "text/html") // Set the content type to HTML
    w.Write([]byte(container)) // Write the HTML structure to the response
}