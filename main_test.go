package main

import (
	"database/sql"
	"testing"
)

func TestDatabaseConnection(t *testing.T) {
	// Attempt to establish a database connection
	db, err := sql.Open("mysql", "root:daddy@(0.0.0.0:3306)/food?parseTime=true")
	if err != nil {
		t.Errorf("Error connecting to the database: %v", err)
		return
	}
	defer db.Close()

	// Check if the database connection is working
	err = db.Ping()
	if err != nil {
		t.Errorf("Error pinging the database: %v", err)
	}
}


