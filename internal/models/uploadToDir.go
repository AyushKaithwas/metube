package models

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const uploadPath = "tmp/videos/"

func UploadToDir(w http.ResponseWriter, r *http.Request, ch chan string) {
    fmt.Println("UploadToDir: Uploading file to directory...")
    
    if err:=r.ParseMultipartForm(10 << 20); err!=nil{
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    } // 10MB

    file, fileHeader, err := r.FormFile("myFile")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        fmt.Println(err)
        return
    }
    defer file.Close()

    allowedMimeTypes := map[string]bool{
        "video/mp4":  true,
        "video/webm": true,
        "video/ogg":  true,
    }

    if _, allowed := allowedMimeTypes[fileHeader.Header.Get("Content-Type")]; !allowed {
        http.Error(w, "The uploaded file type is not allowed.", http.StatusBadRequest)
        return
    }

    newUUID := uuid.New().String()
    ext := filepath.Ext(fileHeader.Filename)

    newFileName := newUUID + ext

    if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
        http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
        return
    }

    dst, err := os.Create(filepath.Join(uploadPath, newFileName))
    if err != nil {
        http.Error(w, "Failed to create the file", http.StatusInternalServerError)
        return
    }
    defer dst.Close()

    if _, err := io.Copy(dst, file); err != nil {
        http.Error(w, "Failed to copy the file", http.StatusInternalServerError)
        return
    }

    ch <- newFileName

    fmt.Fprintf(w, "File uploaded successfully: %s", fileHeader.Filename)
}