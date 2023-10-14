package api

import (
	"database/sql"
	"log"	
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