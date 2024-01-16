package models

import (
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func UploadToStorage(bucketName, localFilePath, destFileName string, ch chan bool) error {
    fmt.Println("UploadToStorage: Uploading file to storage...")
    fmt.Println("Bucket name:", bucketName)
    fmt.Println("Local file path:", localFilePath)
    fmt.Println("Destination file name:", destFileName)

    ctx := context.Background()

    // Use the service account key file
    client, err := storage.NewClient(ctx, option.WithCredentialsFile("service-account-key.json"))
    if err != nil {
        return err
    }
    defer client.Close()

    f, err := os.Open(localFilePath)
    if err != nil {
        return err
    }
    defer f.Close()

    wc := client.Bucket(bucketName).Object(destFileName).NewWriter(ctx)
    if _, err = io.Copy(wc, f); err != nil {
        return err
    }
    if err := wc.Close(); err != nil {
        return err
    }
    ch <- true
    return nil
}
