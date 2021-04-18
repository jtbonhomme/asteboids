package game

import (
	"crypto/rand"
	"fmt"
	"image/color"
	"math/big"
	"time"

	"github.com/jtbonhomme/asteboids/internal/agents"
	"github.com/jtbonhomme/asteboids/internal/physics"
	"github.com/sirupsen/logrus"
)

const (
	defaultScreenWidth         = 1080
	defaultScreenHeight        = 720
	scoreTimeUnit              = 5
	autoGenerateAsteroidsRatio = 10
)

type Game struct {
	log             *logrus.Logger
	gameOver        bool
	gameWon         bool
	nAsteroids      int
	startTime       time.Time
	gameDuration    time.Duration
	highestDuration time.Duration
	highScore       int
	kills           int
	debug           bool
	ScreenWidth     float64
	ScreenHeight    float64
	backgroundColor color.RGBA
	starships       map[string]physics.Physic
	asteroids       map[string]physics.Physic
	bullets         map[string]physics.Physic
}

func New(log *logrus.Logger, nAsteroids int, debug bool) *Game {
	log.Infof("New Game")
	return &Game{
		log:             log,
		gameOver:        false,
		gameWon:         false,
		nAsteroids:      nAsteroids,
		startTime:       time.Now(),
		gameDuration:    0,
		kills:           0,
		highScore:       0,
		highestDuration: 0,
		debug:           debug,
		ScreenWidth:     defaultScreenWidth,
		ScreenHeight:    defaultScreenHeight,
		backgroundColor: color.RGBA{0x00, 0x00, 0x00, 0xff},
		starships:       make(map[string]physics.Physic),
		asteroids:       make(map[string]physics.Physic),
		bullets:         make(map[string]physics.Physic),
	}
}

// StartGame initializes a new game.
func (g *Game) StartGame() {
	// add starship
	p := agents.NewStarship(g.log, g.ScreenWidth/2, g.ScreenHeight/2, g.ScreenWidth, g.ScreenHeight, g.Register, g.Unregister, g.debug)
	g.log.Debugf("added starship: %s", p.ID())
	g.Register(p)

	// add asteroids
	for i := 0; i < g.nAsteroids; i++ {
		g.AddAsteroid()
	}
	g.startTime = time.Now()
	g.gameDuration = 0
	g.gameOver = false
	g.gameWon = false
	g.kills = 0
}

// AddAsteroid insert a new asteroid in the game.
func (g *Game) AddAsteroid() {
	nWidth, err := rand.Int(rand.Reader, big.NewInt(int64(g.ScreenWidth/2)))
	if err != nil {
		g.log.Fatal(err)
	}
	nHeight, err := rand.Int(rand.Reader, big.NewInt(int64(g.ScreenHeight/2)))
	if err != nil {
		g.log.Fatal(err)
	}
	a := agents.NewAsteroid(g.log,
		float64(nWidth.Int64()),
		float64(nHeight.Int64()),
		g.ScreenWidth, g.ScreenHeight,
		g.Register, g.Unregister,
		g.debug)
	g.Register(a)
}

// RestartGame cleans current game and a start a new game.
func (g *Game) RestartGame() {
	for k := range g.starships {
		delete(g.starships, k)
	}
	for k := range g.asteroids {
		delete(g.asteroids, k)
	}
	for k := range g.bullets {
		delete(g.bullets, k)
	}

	g.StartGame()
}

// Register adds a new agent (player or ai) to the game.
func (g *Game) Register(agent physics.Physic) {
	switch agent.Type() {
	case physics.StarshipAgent:
		g.starships[agent.ID()] = agent
	case physics.AsteroidAgent:
		g.asteroids[agent.ID()] = agent
	case physics.RubbleAgent:
		g.asteroids[agent.ID()] = agent
	case physics.BulletAgent:
		g.bullets[agent.ID()] = agent
	default:
	}
}

// Unregister deletes an agent (player or ai) from the game.
func (g *Game) Unregister(id, agentType string) {
	switch agentType {
	case physics.StarshipAgent:
		delete(g.starships, id)
	case physics.AsteroidAgent:
		delete(g.asteroids, id)
	case physics.RubbleAgent:
		delete(g.asteroids, id)
	case physics.BulletAgent:
		delete(g.bullets, id)
	default:
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(g.ScreenWidth), int(g.ScreenHeight)
}

func (g *Game) String() string {
	return fmt.Sprintf(`Asteboids
	- screen size: %d x %d
`, g.ScreenWidth, g.ScreenHeight)
}
