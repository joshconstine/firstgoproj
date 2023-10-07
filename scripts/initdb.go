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
    // Create the ingredient TYPE table
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS ingredient_type (
          ingredient_type_id INT AUTO_INCREMENT PRIMARY KEY,
          name VARCHAR(255) NOT NULL
        );
    `)
    if err != nil {
        panic(err)
    }


    // Create the ingredients table
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS ingredients (
            ingredient_id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            ingredient_type_id INT,
            FOREIGN KEY (ingredient_type_id) REFERENCES ingredient_type(ingredient_type_id)
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
    // Create the recipe_photos table
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS recipe_photos (
            recipe_id INT,
            photo_url VARCHAR(255),
            PRIMARY KEY (recipe_id, photo_url),
            FOREIGN KEY (recipe_id) REFERENCES recipes(recipe_id)
        );
    `)
    if err != nil {
        panic(err)
    }



    // Seed the ingredient types
    _, err = db.Exec(`
    INSERT INTO ingredient_type (name)
    VALUES
        ('Vegetable'),
        ('Fruit'),
        ('Spice'),
        ('Dairy'),
        ('Meet'),
        ('Grain'),
        ('Snack'),
        ('Baking');
    `)
    if err != nil {
        panic(err)
    } 
      // Seed the ingredient types
    _, err = db.Exec(`
    INSERT INTO ingredients (name, ingredient_type_id)
    VALUES
        ('Apple', '2'),
        ('Pear', '2'),
        ('Carrot', '1'),
        ('Broccoli', '1'),
        ('Cinnamon', '3'),
        ('Nutmeg', '3'),
        ('Milk', '4'),
        ('Cheese', '4'),
        ('Chicken', '5'),
        ('Beef', '5'),
        ('Rice', '6'),
        ('Wheat', '6'),
        ('Potato Chips', '7'),
        ('Chocolate Bar', '7'),
        ('Flour', '8'),
        ('Sugar', '8');
    `)
    if err != nil {
        panic(err)
    }
    fmt.Println("Database initialization completed.")
}
