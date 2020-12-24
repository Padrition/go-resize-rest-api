package main

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"

	"github.com/nfnt/resize"
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
	fmt.Printf("\nMIME Name: %+v\n", handler.Header)

	if imageType := handler.Header.Get("Content-Type"); imageType != "image/png" && imageType != "image/jpeg" && imageType != "image/gif" {
		fmt.Println(errors.New("\nEror.A file should be either png, jpeg or gif"))
		http.Error(wr, "Inavalid file format", http.StatusBadRequest)
		return
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println(err)
	}

	automaticResize(img)
}
func automaticResize(img image.Image) {
	m := resize.Resize(500, 0, img, resize.Lanczos2)
	out, err := os.Create("resources/images/test_resize.jpg")
	if err != nil {
		fmt.Println(err)
	}

	jpeg.Encode(out, m, nil)
}
func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", uploadAnImage)
	http.ListenAndServe(":8080", nil)
}
