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

//go:embed karmatic-arcade.ttf
var karmaticArcadeFontData []byte

//go:embed arcadeclassic.ttf
var arcadeClassicFontData []byte

var FurturisticRegularFont font.Face
var MonoSansRegularFont font.Face
var KarmaticArcadeFont font.Face
var ArcadeClassicFont font.Face

func init() {
	var err error

	futuristicRegularFont, err := truetype.Parse(furturisticFontData)
	if err != nil {
		log.Fatal(err)
	}
	FurturisticRegularFont = truetype.NewFace(futuristicRegularFont, &truetype.Options{
		Size:    55,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	monoSansRegularFont, err := truetype.Parse(monoSansFontData)
	if err != nil {
		log.Fatal(err)
	}
	MonoSansRegularFont = truetype.NewFace(monoSansRegularFont, &truetype.Options{
		Size:    10,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	karmaticArcadeFont, err := truetype.Parse(karmaticArcadeFontData)
	if err != nil {
		log.Fatal(err)
	}
	KarmaticArcadeFont = truetype.NewFace(karmaticArcadeFont, &truetype.Options{
		Size:    70,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	arcadeClassicFont, err := truetype.Parse(arcadeClassicFontData)
	if err != nil {
		log.Fatal(err)
	}
	ArcadeClassicFont = truetype.NewFace(arcadeClassicFont, &truetype.Options{
		Size:    50,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}
