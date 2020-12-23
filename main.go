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

	contentType := http.DetectContentType(buff)

	if contentType != "image/gif" && contentType != "image/png" && contentType != "image/jpeg" {
		fmt.Println(errors.New("\nError. The file is neither a gif, png or jpeg"))
		return
	}

	http.ServeContent(wr, r, handler.Filename, time.Now(), file)

}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", uploadAnImage)
	http.ListenAndServe(":8080", nil)
}
