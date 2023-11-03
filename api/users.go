
package api

import (
	"database/sql"
	"log"
	"fmt"
	"html/template"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"	
	"github.com/srinathgs/mysqlstore"
	"github.com/google/uuid"
)
// User struct represents a user's data
type User struct {
    ID       int
    Username string
	PhoneNumber string
}

type ProfilePageData struct {
	PageTitle string
	Username interface{}
	FavoritedRecipes []RecipeWithIngredientsAndPhotosAndTags
	PhoneNumber string
}
func getFavoritedRecipes(db *sql.DB, userID int) []RecipeWithIngredientsAndPhotosAndTags {
	var recipes []RecipeWithIngredientsAndPhotosAndTags
	rows, err := db.Query("SELECT recipe_id FROM user_favorite_recipes WHERE user_id = ? ", userID)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var recipeId int
		err := rows.Scan(&recipeId)
		if err != nil {
			log.Println(err)
		}
		recipeIdString := fmt.Sprintf("%d", recipeId)
		
		recipe, err := getSingleRecipeWithIngredientsAndPhotosAndTags(db,recipeIdString)
		if err != nil {
			log.Println(err)
		}


		recipes = append(recipes, recipe)

	}
	return recipes

}

func GetUserFromRequest(w http.ResponseWriter,r *http.Request, db *sql.DB, store *mysqlstore.MySQLStore) (User, error) {
	c, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Printf("session_token not found")
				return User{}, err
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return User{}, err
		}
		sessionToken := c.Value
	
		if sessionToken == "" {
			// http.Error(w, "Unauthorized, please sign in to view this page", http.StatusUnauthorized)
			return User{}, err
		}
	
		userSession, err := store.Get(r, sessionToken)
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				http.Error(w, "cookie not found", http.StatusBadRequest)
			default:
				log.Println(err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
			return User{}, err
		}
		userID := userSession.Values["user_id"]
		username := userSession.Values["username"]

		var phoneNumber string

		err = db.QueryRow("SELECT phone_number FROM user_info WHERE user_id=?", userID).Scan(&phoneNumber)
		if err != nil {
			log.Println(err)
		}



		return User{
			ID: userID.(int),
			Username: username.(string),
			PhoneNumber: phoneNumber,
		}, nil
	}

func ProfileHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, store *mysqlstore.MySQLStore) {
	user, err := GetUserFromRequest(w,r, db, store)	
	if err != nil {
		log.Println(err)
		return
	}
	favoritedRecipes := getFavoritedRecipes(db, user.ID)
	tmpl := template.Must(template.ParseFiles("public/profile.html"))
	data := ProfilePageData{
		PageTitle: "Profile",
        Username: user.Username,
		FavoritedRecipes: favoritedRecipes,
		PhoneNumber: user.PhoneNumber,
    }
    tmpl.Execute(w, data)
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
	var userId int
	err = db.QueryRow("SELECT password, id FROM users WHERE username=?", username).Scan(&hashedPassword, &userId)
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
	session.Values["user_id"] = userId

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
func LogoutHandler(w http.ResponseWriter, r *http.Request, store *mysqlstore.MySQLStore, db *sql.DB) {
	   //fix this session id
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
 
   store.Delete(r, w, session)
	http.SetCookie(w, &http.Cookie{
	   Name:   "session_token",
        Value:  "",
        MaxAge: 86400, // Session duration (in seconds)
		HttpOnly: true,
		Path: "/",
        Secure:   true,
        SameSite: http.SameSiteLaxMode,
	})
	log.Println("reset cookie")


	log.Println("User logged out")
    w.Write([]byte(""))
}
func UpdateUserPhoneNumber(w http.ResponseWriter, r *http.Request, db *sql.DB, store *mysqlstore.MySQLStore) {
	phoneNumber := r.FormValue("phone")
	user, err := GetUserFromRequest(w,r, db, store)
	var userExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users_info WHERE user_id = ?)", user.ID).Scan(&userExists)
	if userExists && err != nil {
		if  _, err = db.Exec("UPDATE user_info SET phone_number = ? WHERE id = ?", phoneNumber, user.ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		} else if !userExists {
		if  _, err = db.Exec("INSERT INTO user_info (user_id, phone_number) VALUES (?, ?)", user.ID, phoneNumber); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}	else {
		log.Println(userExists)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Redirect back to the page or provide a response
	fmt.Fprintf(w, `<script>window.location.href = "/profile";</script>`)
}