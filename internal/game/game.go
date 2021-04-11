package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/sirupsen/logrus"
)

const (
	defaultScreenWidth  = 640
	defaultScreenHeight = 480
)

type Game struct {
	log             *logrus.Logger
	ScreenWidth     int
	ScreenHeight    int
	backgroundColor color.RGBA
}

func New(log *logrus.Logger) *Game {
	log.Infof("New Game")
	return &Game{
		log:             log,
		ScreenWidth:     defaultScreenWidth,
		ScreenHeight:    defaultScreenHeight,
		backgroundColor: color.RGBA{0x00, 0x00, 0x00, 0xff},
	}
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update(image *ebiten.Image) error {
	// Write your game's logical update.
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Write your game's rendering.
	// Draw the ground image.
	err := screen.Fill(g.backgroundColor)
	if err != nil {
		g.log.Errorf("error while drawing background: %s", err.Error())
	}

	// Draw the message.
	usage := "s: take a screenshot\nCmd + q: exit"
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\n%s", ebiten.CurrentTPS(), ebiten.CurrentFPS(), usage)
	ebitenutil.DebugPrint(screen, msg)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.ScreenWidth, g.ScreenHeight
}

func (g *Game) String() string {
	return fmt.Sprintf(`Asteboids
	- screen size: %d x %d
`, g.ScreenWidth, g.ScreenHeight)
}
