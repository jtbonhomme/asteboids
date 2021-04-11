package asteboids

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jtbonhomme/asteboids/internal/game"
	"github.com/jtbonhomme/asteboids/internal/player"
	"github.com/sirupsen/logrus"
)

func Run(log *logrus.Logger) error {
	os.Setenv("EBITEN_SCREENSHOT_KEY", "s")
	g := game.New(log)
	log.Infof("Game: %s", g)
	ebiten.SetWindowSize(g.ScreenWidth, g.ScreenHeight)
	ebiten.SetWindowTitle("Asteboids")

	p := player.NewStarship(log, g.ScreenWidth/2, g.ScreenHeight/2)
	log.Infof("added created player: %+v", p)
	g.Register(p)
	// Call ebiten.RunGame to start your game loop.
	err := ebiten.RunGame(g)
	if err != nil {
		return err
	}
	return nil
}
