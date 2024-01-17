package controllers

import (
	"fmt"
	"metube/internal/models"
	"net/http"
	"os"
)

// Channels to synchronize the execution
func UploadTrancoderHandler(w http.ResponseWriter, r *http.Request) {

	uploadDone := make(chan string)
	transcodeDone := make(chan bool)
	storageDone := make(chan bool)

	go models.UploadToDir(w, r, uploadDone)
	newFileName := <-uploadDone // Wait for uploadToDir to finish
	// fmt.Println("New file name:", newFileName)
	
	go models.Transcode(newFileName, transcodeDone)
	<-transcodeDone // Wait for transcode to finish

	// k := os.Stdout
	bucketName := os.Getenv("BUCKET_NAME")
    // objectName := newFileName
	quality := "360"
	go models.UploadToStorage(bucketName, "/tmp/output/"+"output_"+newFileName+quality+".mp4", newFileName+quality+".mp4", storageDone)
	<-storageDone // Wait for uploadToStorage to finish

	fmt.Println("All tasks completed.")
}