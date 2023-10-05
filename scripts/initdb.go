package main

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, err := sql.Open("mysql", "root:daddy@(db:3306)/?parseTime=true")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // Create the food database
    _, err = db.Exec("CREATE DATABASE IF NOT EXISTS food;")
    if err != nil {
        panic(err)
    }

    // Connect to the food database
    _, err = db.Exec("USE food;")
    if err != nil {
        panic(err)
    }

    // Create the ingredients table
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS ingredients (
            ingredient_id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(255) NOT NULL
        );
    `)
    if err != nil {
        panic(err)
    }

    // Create the recipes table
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS recipes (
            recipe_id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            description TEXT
        );
    `)
    if err != nil {
        panic(err)
    }

    // Create the recipe_ingredients table
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS recipe_ingredients (
            recipe_id INT,
            ingredient_id INT,
            PRIMARY KEY (recipe_id, ingredient_id),
            FOREIGN KEY (recipe_id) REFERENCES recipes(recipe_id),
            FOREIGN KEY (ingredient_id) REFERENCES ingredients(ingredient_id)
        );
    `)
    if err != nil {
        panic(err)
    }

    fmt.Println("Database initialization completed.")
}
