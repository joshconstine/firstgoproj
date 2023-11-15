package api

import (
	"database/sql"
	"log"	
	"encoding/json"
    "net/http"
	"github.com/gorilla/mux"
)
type MobileUser struct {
	Clerk_id string
	Photo_url string
	Username string 
}


	func GetMobileUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
		id := mux.Vars(r)["id"]
		
		rows, err := db.Query(`SELECT * FROM mobile_users WHERE clerk_id = ?`, id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var mobileUser MobileUser

		for rows.Next() {
			err := rows.Scan(&mobileUser.Clerk_id, &mobileUser.Photo_url, &mobileUser.Username)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Convert the recipe to JSON
		json, err := json.Marshal(mobileUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}