package asteboids

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jtbonhomme/asteboids/internal/config"
	"github.com/jtbonhomme/asteboids/internal/game"
	"github.com/sirupsen/logrus"
)

func Run(log *logrus.Logger, conf *config.Config) error {
	os.Setenv("EBITEN_SCREENSHOT_KEY", "s")
	g := game.New(log, conf)
	log.Infof("Game: %s", g)
	ebiten.SetWindowSize(int(conf.ScreenWidth), int(conf.ScreenHeight))
	ebiten.SetWindowTitle("Asteboids")

	g.StartGame()

	if conf.Optim {
		ebiten.SetVsyncEnabled(false)
		ebiten.SetInitFocused(false)
	}
	if conf.MaxTPS != 0 {
		ebiten.SetMaxTPS(conf.MaxTPS)
	}

	// Call ebiten.RunGame to start your game loop.
	err := ebiten.RunGame(g)
	if err != nil {
		return err
	}
	return nil
}
