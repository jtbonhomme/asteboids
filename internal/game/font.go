package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const dpi = 72

var MplusNormalFont font.Face

func init() {
	var err error

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	MplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		logrus.Fatal(err)
	}
}
