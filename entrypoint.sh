#!/bin/bash

# Function to wait for the MySQL container to become healthy
wait_for_mysql() {
  until mysqladmin ping -h"db" -uroot -pdaddy; do
    echo "Waiting for MySQL server to become healthy..."
    sleep 1
  done
  echo "MySQL server is now healthy."
}

# Wait for the MySQL container to become healthy
wait_for_mysql
# Run your Go script
 go run ./scripts/initdb.go

# Start your Go application (replace with the appropriate command)
 go run ./main.go
