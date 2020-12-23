package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

func index(wr http.ResponseWriter, r *http.Request) {
	http.ServeFile(wr, r, "resources/html/index.html")
}

func errorMessage(wr http.ResponseWriter, r *http.Request) {
	http.ServeFile(wr, r, "resources/html/error.html")
}
func uploadAnImage(wr http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 * 1024 * 1024)

	file, handler, err := r.FormFile("pic")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Printf("\nFile Name: %+v", handler.Filename)
	fmt.Printf("\nFile Name: %+v", handler.Size)
	fmt.Printf("\nMIME Name: %+v", handler.Header)

	buff := make([]byte, 512)

	if _, err = file.Read(buff); err != nil {
		fmt.Println(err)
		return
	}

	if imageType := http.DetectContentType(buff); imageType != "image/png" && imageType != "image/jpeg" && imageType != "image/gif" {
		fmt.Println(errors.New("\nEror.A file should be either png, jpeg or gif"))
		http.Error(wr, "Inavalid file format", http.StatusBadRequest)
	} else {
		http.ServeContent(wr, r, handler.Filename, time.Now(), file)
	}

}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", uploadAnImage)
	http.ListenAndServe(":8080", nil)
}
