package game

import (
	"github.com/sirupsen/logrus"
)

type Game struct {
	log *logrus.Logger
}

func New(log *logrus.Logger) *Game {
	log.Infof("New Game")
	return &Game{
		log: log,
	}
}

func (g *Game) String() string {
	return `Asteboids Game`
}
