package controllers

import (
	"fmt"
	"metube/internal/models"
)

// Channels to synchronize the execution
func UploadTrancoderHandler() {

	uploadDone := make(chan bool)
	// transcodeDone := make(chan bool)
	// storageDone := make(chan bool)

	go models.UploadToDir(uploadDone)
	<-uploadDone // Wait for uploadToDir to finish

	// go transcode(transcodeDone)
	// <-transcodeDone // Wait for transcode to finish

	// go uploadToStorage(storageDone)
	// <-storageDone // Wait for uploadToStorage to finish

	fmt.Println("All tasks completed.")
}