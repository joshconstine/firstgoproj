package api

import (
	"database/sql"
	"log"	
    "encoding/json"
    "net/http"
)
type Tag struct {
	Tag_id int
	Name  string
}

	func getAllTags(db *sql.DB) []Tag {
	rows, err := db.Query(`SELECT * FROM tags`)
        if err != nil {
			log.Fatal(err)
        }
        defer rows.Close()
		
        var tags []Tag
        for rows.Next() {
			var t Tag
			
            err := rows.Scan(&t.Tag_id, &t.Name)
            if err != nil {
				log.Fatal(err)
            }
            tags = append(tags, t)
        }
        if err := rows.Err(); err != nil {
			log.Fatal(err)
        }
		return tags
}	
func getTagsforRecipeId(db *sql.DB, recipeId string) []Tag {
	rows, err := db.Query(`SELECT t.tag_id, t.name
	FROM tags t
	INNER JOIN recipe_tags rt ON t.tag_id = rt.tag_id
	WHERE rt.recipe_id = ?;
	`, recipeId)
        if err != nil {
			log.Fatal(err)
        }
        defer rows.Close()
		
        var tags []Tag
        for rows.Next() {
			var t Tag
			
            err := rows.Scan(&t.Tag_id, &t.Name)
            if err != nil {
				log.Fatal(err)
            }
            tags = append(tags, t)
        }
        if err := rows.Err(); err != nil {
			log.Fatal(err)
        }
		return tags
}
func GetTagsJSON(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    tags := getAllTags(db)
    
    tagsJSON, err := json.Marshal(tags)
    if err != nil {
        log.Fatal(err)
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(tagsJSON)

}