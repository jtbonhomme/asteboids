package asteboids

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jtbonhomme/asteboids/internal/agents"
	"github.com/jtbonhomme/asteboids/internal/game"
	"github.com/sirupsen/logrus"
)

const (
	maxAsteroids int = 6
)

func Run(log *logrus.Logger, optim, debug bool) error {
	os.Setenv("EBITEN_SCREENSHOT_KEY", "s")
	g := game.New(log)
	log.Infof("Game: %s", g)
	ebiten.SetWindowSize(g.ScreenWidth, g.ScreenHeight)
	ebiten.SetWindowTitle("Asteboids")

	// add starship
	p := agents.NewStarship(log, g.ScreenWidth/2, g.ScreenHeight/2, g.ScreenWidth, g.ScreenHeight, g.Unregister, debug)
	log.Infof("added starship: %+v", p)
	g.Register(p)
	g.StarshipID = p.ID()

	// add asteroids
	for i := 0; i < maxAsteroids; i++ {
		p := agents.NewAsteroid(log, g.ScreenWidth, g.ScreenHeight, g.Unregister, debug)
		g.Register(p)
	}

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
