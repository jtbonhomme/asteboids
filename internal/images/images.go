package images

import (
	"bytes"
	// import embed to load image files
	_ "embed"
	"errors"
	"image"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed asteroid0.png
var asteroid0PNG []byte

//go:embed asteroid1.png
var asteroid1PNG []byte

//go:embed asteroid2.png
var asteroid2PNG []byte

//go:embed asteroid3.png
var asteroid3PNG []byte

//go:embed asteroid4.png
var asteroid4PNG []byte

//go:embed rubble0.png
var rubble0PNG []byte

//go:embed rubble1.png
var rubble1PNG []byte

//go:embed rubble2.png
var rubble2PNG []byte

//go:embed rubble3.png
var rubble3PNG []byte

//go:embed rubble4.png
var rubble4PNG []byte

//go:embed ship.png
var shipPNG []byte

//go:embed bullet.png
var bulletPNG []byte

//go:embed boid.png
var boidPNG []byte

// loadImageFromSlice loads an image into an ebiten image.
func loadImageFromImage(screenWidth, screenHeight int, rawImage image.Image) (*ebiten.Image, error) {
	newImage := ebiten.NewImage(screenWidth, screenHeight)
	newImageFromImage := ebiten.NewImageFromImage(rawImage)
	if newImageFromImage == nil {
		return newImage, errors.New("error when creating image from raw")
	}
	newImage.DrawImage(newImageFromImage, nil)
	return newImage, nil

}

// LoadImageFromSlice loads an image into an ebiten image from a slice of bytes.
func LoadImageFromSlice(screenWidth, screenHeight int, name string) (*ebiten.Image, error) {
	var data []byte
	switch name {
	case "boid.png":
		data = boidPNG
	case "ship.png":
		data = shipPNG
	case "bullet.png":
		data = bulletPNG
	case "asteroid0.png":
		data = asteroid0PNG
	case "asteroid1.png":
		data = asteroid1PNG
	case "asteroid2.png":
		data = asteroid2PNG
	case "asteroid3.png":
		data = asteroid3PNG
	case "asteroid4.png":
		data = asteroid4PNG
	case "rubble0.png":
		data = rubble0PNG
	case "rubble1.png":
		data = rubble1PNG
	case "rubble2.png":
		data = rubble2PNG
	case "rubble3.png":
		data = rubble3PNG
	case "rubble4.png":
		data = rubble4PNG
	default:
		return nil, errors.New("error, unknown rawdata slice")
	}
	buf := bytes.NewBuffer(data)

	rawImage, _, err := image.Decode(buf)
	if err != nil {
		return nil, errors.New("error when decoding image from file " + err.Error())
	}
	return loadImageFromImage(screenWidth, screenHeight, rawImage)
}

// LoadImageFromFile loads an image into an ebiten image from a file.
func LoadImageFromFile(screenWidth, screenHeight int, file string) (*ebiten.Image, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.New("error when opening file " + err.Error())
	}

	defer f.Close()
	rawImage, _, err := image.Decode(f)
	if err != nil {
		return nil, errors.New("error when decoding image from file " + err.Error())
	}

	return loadImageFromImage(screenWidth, screenHeight, rawImage)
}
