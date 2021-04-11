package asteboids

import (
	"github.com/jtbonhomme/asteboids/game"
	"github.com/sirupsen/logrus"
)

func Run(log *logrus.Logger) {
	g := game.New(log)
	log.Infof("Game: %s", g)
}
