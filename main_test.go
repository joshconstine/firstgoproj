package main

import (
	"testing"
)

func TestDatabaseConnection(t *testing.T) {
	// Attempt to establish a database connection
	db, err := establishdbConnection("root", "daddy", "0.0.0.0", "3306", "food")
	if err != nil {
		t.Errorf("Error connecting to the database: %v", err)
		return
	}
	defer db.Close()
}
func TestPingDatabase(t *testing.T) {

	db, err := establishdbConnection("root", "daddy", "0.0.0.0", "3306", "food")
	if err != nil {
		t.Errorf("Error connecting to the database: %v", err)
		return
	}
	defer db.Close()

	// Attempt to ping the database
	err = db.Ping()
	if err != nil {
		t.Errorf("Error pinging the database: %v", err)
		return
	}
	
}


