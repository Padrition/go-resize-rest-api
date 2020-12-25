package main

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/nfnt/resize"
)

func index(wr http.ResponseWriter, r *http.Request) {
	http.ServeFile(wr, r, "resources/html/index.html")
}

func upload(rw http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 * 1024 * 1024)

	imageFile, header, err := r.FormFile("imageFile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer imageFile.Close()

	imageType := header.Header.Get("Content-Type")
	if imageType != "image/png" && imageType != "image/jpeg" && imageType != "image/gif" {
		fmt.Println(errors.New("\nEror.A file should be either png, jpeg or gif"))
		http.Error(rw, "Inavalid file format", http.StatusBadRequest)
		return
	}

	resizeAnImage(500, imageFile, imageType, header)
}

func resizeAnImage(width uint, imageFile multipart.File, imageType string, header *multipart.FileHeader) {
	switch imageType {
	case "image/jpeg":

		img, err := jpeg.Decode(imageFile)
		if err != nil {
			fmt.Println(err)
		}

		jpegImg := resize.Resize(width, 0, img, resize.Lanczos2)

		out, err := os.Create("resources/images/" + header.Filename)
		if err != nil {
			fmt.Println(err)
		}

		jpeg.Encode(out, jpegImg, nil)
		break

	case "image/png":
		img, err := png.Decode(imageFile)
		if err != nil {
			fmt.Println(err)
		}

		pngImg := resize.Resize(width, 0, img, resize.Lanczos2)

		out, err := os.Create("resources/images/" + header.Filename)
		if err != nil {
			fmt.Println(err)
		}

		png.Encode(out, pngImg)
		break

	case "image/gif":
		newGifImg := gif.GIF{}
		gifImg, err := gif.DecodeAll(imageFile)
		if err != nil {
			fmt.Println(err)
		}

		out, err := os.Create("resources/images/" + header.Filename)
		if err != nil {
			fmt.Println(err)
		}

		for _, img := range gifImg.Image {
			resizedGifImg := resize.Resize(width, 0, img, resize.Lanczos2)
			palettedImg := image.NewPaletted(resizedGifImg.Bounds(), img.Palette)
			draw.FloydSteinberg.Draw(palettedImg, resizedGifImg.Bounds(), resizedGifImg, image.ZP)

			newGifImg.Image = append(newGifImg.Image, palettedImg)
			newGifImg.Delay = append(newGifImg.Delay, 25)
		}

		gif.EncodeAll(out, &newGifImg)

		break
	}

}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", upload)
	http.ListenAndServe(":8080", nil)
}
