package controllers

import (
	"metube/internal/models"
	"net/http"
)

// Channels to synchronize the execution
func UploadTrancoderHandler(w http.ResponseWriter, r *http.Request) {

	uploadDone := make(chan string)
	// transcodeDone := make(chan bool)
	// storageDone := make(chan bool)

	go models.UploadToDir(w, r, uploadDone)
	newFileName := <-uploadDone // Wait for uploadToDir to finish
	// fmt.Println("New file name:", newFileName)
	
	go models.Transcode(newFileName)
	// <-transcodeDone // Wait for transcode to finish

	// go uploadToStorage(storageDone)
	// <-storageDone // Wait for uploadToStorage to finish

	// fmt.Println("All tasks completed.")
}