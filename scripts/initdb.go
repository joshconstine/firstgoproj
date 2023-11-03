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

    var dbName string
err = db.QueryRow("SHOW DATABASES LIKE 'food'").Scan(&dbName)
if err == nil {
    // The "food" database already exists
    fmt.Println("The 'food' database already exists. Skipping initilization.")
} else if err == sql.ErrNoRows {
    // The "food" database doesn't exist, so create it
    _, err := db.Exec("CREATE DATABASE IF NOT EXISTS food;")
    
    if err != nil {
        panic(err)
    }
    _, err = db.Exec("USE food;")
    if err != nil {
     panic(err)
    }
// Create the users TYPE table
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INT AUTO_INCREMENT PRIMARY KEY,
            username VARCHAR(255) UNIQUE NOT NULL,
            password VARCHAR(255) NOT NULL
        );
    `)
    if err != nil {
        panic(err)
    } 
       _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS user_info (
            id INT AUTO_INCREMENT PRIMARY KEY,
            user_id INT NOT NULL,
            phone_number VARCHAR(20),
            UNIQUE (user_id),
            FOREIGN KEY (user_id) REFERENCES users(id)
        );
    `)
    if err != nil {
        panic(err)
    }

// Create the user_favorite_recipes TYPE table
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS user_favorite_recipes (
            user_id INT,
            recipe_id INT,
            PRIMARY KEY (user_id, recipe_id),
            FOREIGN KEY (user_id) REFERENCES users(id),
            FOREIGN KEY (recipe_id) REFERENCES recipes(recipe_id)
        );
    `)
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
    // Create the ingredient Quantity type table
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS quantity_type (
          quantity_type_id INT AUTO_INCREMENT PRIMARY KEY,
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
            quantity_type_id INT,
            quantity FLOAT,
            PRIMARY KEY (recipe_id, ingredient_id),
            FOREIGN KEY (recipe_id) REFERENCES recipes(recipe_id),
            FOREIGN KEY (ingredient_id) REFERENCES ingredients(ingredient_id),
            FOREIGN KEY (quantity_type_id) REFERENCES quantity_type(quantity_type_id)
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
         // Create the tags table
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS tags (
            tag_id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(255) NOT NULL UNIQUE
        );
    `)
    if err != nil {
        panic(err)
    }

// Create the recipe_tags table
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS recipe_tags (
            recipe_id INT,
            tag_id INT,
            PRIMARY KEY (recipe_id, tag_id),
            FOREIGN KEY (recipe_id) REFERENCES recipes(recipe_id),
            FOREIGN KEY (tag_id) REFERENCES tags(tag_id)
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
        ('Alcohol'),
        ('Grain'),
        ('Snack'),
        ('Baking');
    `)
    if err != nil {
        panic(err)
    } 
    // Seed the quantitiy types
    _, err = db.Exec(`
    INSERT INTO quantity_type (name)
    VALUES
        (''),
        ('Cup'),
        ('Ounce'),
        ('Tablespoon'),
        ('Teaspoon'),
        ('Pound'),
        ('Gram');
    `)
    if err != nil {
        panic(err)
    } 
  // Seed the tags
    _, err = db.Exec(`
    INSERT INTO tags (name)
    VALUES
        ('Recipes With Video'),
        ('Munchies'),
        ('Soup'),
        ('Pasta'),
        ('Beverages'),
        ('Salad'),
        ('Desserts'),
        ('Fish And Seafood'),
        ('Cool Stuff'),
        ('Appitizers'),
        ('Main Dishes');
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
        ('Bananna', '2'),
        ('Grape', '2'),
        ('Strawberry', '2'),
        ('Blueberry', '2'),
        ('Grape Jelly', '2'),
        ('Carrot', '1'),
        ('Broccoli', '1'),
        ('Red Pepper', '1'),
        ('Green Pepper', '1'),
        ('Jalapeno Pepper', '1'),
        ('Yellow Onion', '1'),
        ('Mushroom', '1'),
        ('Onion', '1'),
        ('Salt', '3'),
        ('Black Pepper', '3'),
        ('Cinnamon', '3'),
        ('Cumin', '3'),
        ('Red Pepper Flakes', '3'),
        ('Onion Powder', '3'),
        ('Garlic Powder', '3'),
        ('Nutmeg', '3'),
        ('Milk', '4'),
        ('Butter', '4'),
        ('Yogurt', '4'),
        ('Egg', '4'),
        ('Chedar Cheese', '4'),
        ('Parmessean Cheese', '4'),
        ('Heavy Whipping Creme', '4'),
        ('Chicken', '5'),
        ('Beef', '5'),
        ('Salmon', '5'),
        ('Pork', '5'),
        ('Tuna', '5'),
        ('Bison', '5'),
        ('Bacon', '5'),
        ('Rice', '6'),
        ('Bread', '6'),
        ('Peanut Butter', '6'),
        ('Wheat', '6'),
        ('Potato Chips', '7'),
        ('Chocolate Bar', '7'),
        ('Flour', '8'),
        ('Sugar', '8'),
        ('Baking Soda', '8'),
        ('Vanilla Extract', '8'),
        ('Baking Powder', '8'),
        ('Brown Sugar', '8');
    `)
    if err != nil {
        panic(err)
    }
    fmt.Println("seeded the 'food' database.")
    fmt.Println("Database initialization completed.")


} else {
    panic(err)
}

// Connect to the "food" database
_, err = db.Exec("USE food;")
if err != nil {
    panic(err)
}
}



