
package api

import (
	"database/sql"
	"log"
	"golang.org/x/crypto/bcrypt"
	"net/http"	
	"github.com/srinathgs/mysqlstore"
	"github.com/google/uuid"
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

func LoginUser(w http.ResponseWriter, r *http.Request, db *sql.DB, store *mysqlstore.MySQLStore) {
	err := r.ParseForm()
	if err != nil {
	   http.Error(w, "Please pass the data as URL form encoded",
  http.StatusBadRequest)
	  return
	}
	username := r.FormValue( "username")
	password := r.FormValue( "password")
	var hashedPassword string
	err = db.QueryRow("SELECT password FROM users WHERE username=?", username).Scan(&hashedPassword)
	if err != nil {
		
	container := "<div  id=\"successContainer\" data-hx-target=\"ingredientList\" class=\"block w-full rounded-lg p-3 flex h-full justify-center max-h-full flex-col items-center \" >"
			container += `<h1 class="text-m"> User not found</h1> <h1 class="text-m">` + username + `</h1>`
		container += "</div>"
	w.Header().Set("Content-Type", "text/html") // Set the content type to HTML
	w.Write([]byte(container)) // Write the HTML structure to the response
	return
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		http.Error(w, "Invalid Credentials",http.StatusUnauthorized)
	}

	sessionToken := uuid.NewString()

	session, err := store.Get(r, sessionToken)
	session.Values["username"] = username
	session.Values["authenticated"] = true
	if err := session.Save(r, w); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Set a cookie with the session ID
    cookie := http.Cookie{
        Name:   "session_token",
        Value:  sessionToken,
        MaxAge: 86400, // Session duration (in seconds)
		HttpOnly: true,
		Path: "/",
        Secure:   true,
        SameSite: http.SameSiteLaxMode,
    }
    http.SetCookie(w, &cookie) 
	
	
	container := "<div  id=\"successContainer\" data-hx-target=\"ingredientList\" class=\"block w-full rounded-lg p-3 flex h-full justify-center max-h-full flex-col items-center \" >"
			container += `<h1 class="text-m"> User logged in </h1> <h1 class="text-m">` + username + `</h1>`
		container += "</div>"
	w.Header().Set("Content-Type", "text/html") // Set the content type to HTML
	w.Write([]byte(container)) // Write the HTML structure to the response
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
func LogoutHandler(w http.ResponseWriter, r *http.Request, store *mysqlstore.MySQLStore) {
	   //fix this session id
    session, _ := store.Get(r, "session.id")
    session.Values["authenticated"] = false
    session.Save(r, w)
log.Println("User logged out")
    w.Write([]byte(""))
}