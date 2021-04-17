package fonts

import (
	// import embed to load truetype font
	_ "embed"
	"log"

	"github.com/golang/freetype/truetype"

	"golang.org/x/image/font"
)

const (
	dpi  float64 = 72
	size float64 = 20
)

//go:embed Exan-Regular.ttf
var fontData []byte

var ExanRegularTTF font.Face

func init() {
	var err error

	exanRegularFont, err := truetype.Parse(fontData)
	if err != nil {
		log.Fatal(err)
	}
	ExanRegularTTF = truetype.NewFace(exanRegularFont, &truetype.Options{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}
