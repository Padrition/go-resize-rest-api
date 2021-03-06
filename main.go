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
	"strconv"

	"github.com/nfnt/resize"
)

type imageData struct {
	imageFile multipart.File
	header    *multipart.FileHeader
	imageType string
	imageName string
}

func ifErrNil(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
func index(rw http.ResponseWriter, r *http.Request) {
	http.ServeFile(rw, r, "resources/html/index.html")
}

func autoResize(rw http.ResponseWriter, r *http.Request) {
	var data imageData
	data, width, height := upload(rw, r)

	data.imageName = data.header.Filename
	resizeAnImage(rw, data, uint(width), uint(height))
}

func uploadAnImage(rw http.ResponseWriter, r *http.Request) {
	var data imageData
	data, _, _ = upload(rw, r)

	out, err := os.Create("resources/images/" + data.header.Filename)
	ifErrNil(err)

	defer out.Close()
	if data.imageType != "image/gif" {
		img, _, err := image.Decode(data.imageFile)
		ifErrNil(err)

		switch data.imageType {
		case "image/jpeg":
			jpeg.Encode(out, img, nil)
			break
		case "image/png":
			png.Encode(out, img)
			break
		}
	} else {
		imgGif, err := gif.DecodeAll(data.imageFile)
		ifErrNil(err)

		gif.EncodeAll(out, imgGif)
	}
}

func upload(rw http.ResponseWriter, r *http.Request) (imageData, uint64, uint64) {
	r.ParseMultipartForm(32 * 1024 * 1024)

	imageFile, header, err := r.FormFile("imageFile")
	ifErrNil(err)

	width, err := strconv.ParseUint((r.FormValue("width")), 10, 32)
	ifErrNil(err)
	height, err := strconv.ParseUint((r.FormValue("height")), 10, 32)
	ifErrNil(err)

	defer imageFile.Close()

	imageType := header.Header.Get("Content-Type")
	if imageType != "image/png" && imageType != "image/jpeg" && imageType != "image/gif" {
		fmt.Println(errors.New("\nEror.A file should be either png, jpeg or gif"))
		http.Error(rw, "Inavalid file format", http.StatusBadRequest)
	}

	data := imageData{imageFile: imageFile, header: header, imageType: imageType}

	return data, width, height
}

func resizeAnImage(rw http.ResponseWriter, data imageData, width uint, height uint) {
	if data.imageType != "image/gif" {
		img, _, err := image.Decode(data.imageFile)
		ifErrNil(err)
		resizedImages := resize.Resize(width, height, img, resize.Lanczos2)
		switch data.imageType {
		case "image/jpeg":
			rw.Header().Set("Content-Type", "image/jpeg")
			jpeg.Encode(rw, resizedImages, nil)
			break

		case "image/png":
			rw.Header().Set("Content-Type", "image/png")
			png.Encode(rw, resizedImages)
			break
		}
	} else {
		newGifImg := gif.GIF{}
		gifImg, err := gif.DecodeAll(data.imageFile)
		ifErrNil(err)

		for _, img := range gifImg.Image {
			resizedGifImg := resize.Resize(width, height, img, resize.Lanczos2)
			palettedImg := image.NewPaletted(resizedGifImg.Bounds(), img.Palette)
			draw.FloydSteinberg.Draw(palettedImg, resizedGifImg.Bounds(), resizedGifImg, image.ZP)

			newGifImg.Image = append(newGifImg.Image, palettedImg)
			newGifImg.Delay = append(newGifImg.Delay, 25)
		}
		rw.Header().Set("Content-Type", "image/gif")
		gif.EncodeAll(rw, &newGifImg)

	}
}

func setupRouts() {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", uploadAnImage)
	http.HandleFunc("/resize", autoResize)
	http.ListenAndServe(":8080", nil)
}

func main() {
	setupRouts()
}
