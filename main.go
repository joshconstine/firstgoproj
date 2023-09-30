package main

import (
	"fmt"
	"strconv"
	"net/http"
)


func main() {
	port := 8080
	// Convert the integer port to a string.
	portStr := strconv.Itoa(port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Set the Content-Type header to indicate that we are sending HTML.
		w.Header().Set("Content-Type", "text/html")

		// Write your HTML content to the response writer.
		fmt.Fprintf(w, "<html><body><h1>Hello, Go HTML Server!</h1></body></html>")
	})
	
	fmt.Printf("Server is listening on port %d...\n", port)
	http.HandleFunc("/hello-world", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello World"))
	})
	http.ListenAndServe(":"+portStr, nil)
}