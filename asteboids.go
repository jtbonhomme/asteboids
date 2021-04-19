package asteboids

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jtbonhomme/asteboids/internal/game"
	"github.com/sirupsen/logrus"
)

const (
	nAsteroids int = 4
	nBoids     int = 40
)

func Run(log *logrus.Logger, optim, debug bool) error {
	os.Setenv("EBITEN_SCREENSHOT_KEY", "s")
	g := game.New(log, nAsteroids, nBoids, debug)
	log.Infof("Game: %s", g)
	ebiten.SetWindowSize(int(g.ScreenWidth), int(g.ScreenHeight))
	ebiten.SetWindowTitle("Asteboids")

	g.StartGame()

	if optim {
		ebiten.SetVsyncEnabled(false)
		ebiten.SetInitFocused(false)
	}
	// Call ebiten.RunGame to start your game loop.
	err := ebiten.RunGame(g)
	if err != nil {
		return err
	}
	return nil
}
