package main

import (
	"fmt"
	"os"
	"io"
	"io/ioutil"

	"path/filepath"
	"log"
	"bytes"
	"github.com/google/uuid"
	"strconv"
	"context"

	"net/http"
    "github.com/gorilla/mux"
    "firstgoprog/api" // Replace "firstgoprog" with your actual module name.
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/joho/godotenv/autoload"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
)


func NewUploader() *s3manager.Uploader {
	ACCESS_KEY:= "AKIAX6ZNEPNPAR6OXRRO"
	SECRET_KEY:= "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"
	s3Config := &aws.Config{
		Region:      aws.String("us-west-1"),
		Credentials: credentials.NewStaticCredentials(ACCESS_KEY, SECRET_KEY, ""),
	}

	s3Session := session.New(s3Config)

	uploader := s3manager.NewUploader(s3Session)
	fmt.Printf("Created new S3 Uploder")
	
	return uploader
}

func UploadHandler(w http.ResponseWriter, r *http.Request, uploader *s3manager.Uploader) {
	err := r.ParseMultipartForm(10 * 1024 * 1024) // 10 MB limit
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Generate a unique filename using a UUID
	fileExt := filepath.Ext(header.Filename)
	newFilename := uuid.New().String() + fileExt

	// Create a new file in the "public/static" directory with the unique filename
	newFilePath := filepath.Join("public/static/images", newFilename)
	newFile, err := os.Create(newFilePath)
	if err != nil {
		http.Error(w, "Failed to create a new file", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	// Copy the uploaded file to the new file
	_, err = io.Copy(newFile, file)
	if err != nil {
		http.Error(w, "Failed to copy the file", http.StatusInternalServerError)
		return
	}
	log.Println("uploading so S3")

	// file, err := ioutil.ReadFile(newFilePath)
	// if err != nil {
	// 	log.Fatal(err)
	// }


BUCKET_NAME := "foodly-bucket"
// BUCKET_URL := "https://foodly-bucket.s3.us-west-1.amazonaws.com"
// NEXT_PUBLIC_BUCKET_URL := "https://foodly-bucket.s3.us-west-1.amazonaws.com"

//** try with copy
fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file to bytes", http.StatusInternalServerError)

		return 
	}



	upInput := &s3manager.UploadInput{
		Bucket:      aws.String(BUCKET_NAME), // bucket's name
		Key:         aws.String(newFilePath),        // files destination location
		Body:        bytes.NewReader(fileBytes),                   // content of the file
		ContentType: aws.String(fileExt),                 // content type
	}
	res, err := uploader.UploadWithContext(context.Background(), upInput)
	log.Printf("res %+v\n", res)
	log.Printf("err %+v\n", err)

	// Respond with the unique filename or other relevant information
	fmt.Fprintf(w, "File uploaded successfully with filename: %s", newFilename)
}

func main() {
	port := 8080
	// Convert the integer port to a string.
	portStr := strconv.Itoa(port)
	
	uploader := NewUploader()

	r := mux.NewRouter()
	// r.HandleFunc("/upload", UploadHandler(uploader))
	r.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
        UploadHandler(w, r, uploader)
    }).Methods("POST")	


	db, err := sql.Open("mysql", "root:daddy@(db:3306)/food?parseTime=true")
	
    if err != nil {
		// log.Fatal(err)
		fmt.Print("error connecting to db")
    }
    if err := db.Ping(); err != nil {
		fmt.Printf("Error %d...\n", err)
    }
	staticDir := "/images/"
	r.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("./public/images"))))

	// Your other application routes go here...
    // Use the functions from the 'api' package to define routes.
	api.InitRoutes(r, db)

	fmt.Printf("Server is listening on port %d...\n", port)

	http.ListenAndServe(":"+portStr, r)
}