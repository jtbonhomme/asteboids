package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	agents          []Agent
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

// Register adds a new agent (player or ai) to the game.
func (g *Game) Register(agent Agent) {
	g.agents = append(g.agents, agent)
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Write your game's logical update.
	// Update the agents
	for _, a := range g.agents {
		a.Update()
	}
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Write your game's rendering.
	// Draw the ground image.
	screen.Fill(g.backgroundColor)

	// Draw the agents
	for _, a := range g.agents {
		a.Draw(screen)
	}
	// Draw the message.
	usage := "s: take a screenshot\nCmd + q: exit"
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\n", ebiten.CurrentTPS(), ebiten.CurrentFPS())
	msg += fmt.Sprintf("direction: %0.2f\nspeed: %0.2f\n", g.agents[0].Direction(), g.agents[0].Speed())
	msg += fmt.Sprintf("x, y: %v\n", g.agents[0].Position())
	msg += fmt.Sprintf("%s\n", usage)
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
