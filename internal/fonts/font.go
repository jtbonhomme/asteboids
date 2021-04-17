package fonts

import (
	// import embed to load truetype font
	_ "embed"
	"log"

	"github.com/golang/freetype/truetype"

	"golang.org/x/image/font"
)

const (
	dpi float64 = 72
)

//go:embed Exan-Regular.ttf
var furturisticFontData []byte

//go:embed NotoMono-Regular.ttf
var monoSansFontData []byte

var FurturisticRegularFontTitle font.Face
var MonoSansRegularFont10 font.Face

func init() {
	var err error

	futuristicRegularFont, err := truetype.Parse(furturisticFontData)
	if err != nil {
		log.Fatal(err)
	}
	FurturisticRegularFontTitle = truetype.NewFace(futuristicRegularFont, &truetype.Options{
		Size:    40,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	monoSansRegularFont, err := truetype.Parse(monoSansFontData)
	if err != nil {
		log.Fatal(err)
	}
	MonoSansRegularFont10 = truetype.NewFace(monoSansRegularFont, &truetype.Options{
		Size:    10,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}
