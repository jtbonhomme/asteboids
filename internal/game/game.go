package game

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jtbonhomme/asteboids/internal/agents"
	"github.com/jtbonhomme/asteboids/internal/ai"
	"github.com/jtbonhomme/asteboids/internal/physics"
	"github.com/sirupsen/logrus"
)

const (
	defaultScreenWidth         float64 = 1080
	defaultScreenHeight        float64 = 720
	scoreTimeUnit              float64 = 5
	autoGenerateAsteroidsRatio float64 = 10
)

type Game struct {
	log             *logrus.Logger
	gameOver        bool
	gameWon         bool
	nAsteroids      int
	nBoids          int
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
	boids           map[string]physics.Physic
	starshipImage   *ebiten.Image
	bulletImage     *ebiten.Image
	boidImage       *ebiten.Image
	asteroidImages  []*ebiten.Image
	rubbleImages    []*ebiten.Image
}

func New(log *logrus.Logger,
	nAsteroids, nBoids int,
	debug bool) *Game {
	log.Infof("New Game")
	rand.Seed(time.Now().UnixNano())
	g := &Game{
		log:             log,
		gameOver:        false,
		gameWon:         false,
		nAsteroids:      nAsteroids,
		nBoids:          nBoids,
		startTime:       time.Now(),
		gameDuration:    0,
		kills:           0,
		highScore:       0,
		highestDuration: 0,
		debug:           debug,
		ScreenWidth:     defaultScreenWidth,
		ScreenHeight:    defaultScreenHeight,
		backgroundColor: 	BgColor  = color.Black,
		starships:       make(map[string]physics.Physic),
		asteroids:       make(map[string]physics.Physic),
		bullets:         make(map[string]physics.Physic),
		boids:           make(map[string]physics.Physic),
		asteroidImages:  make([]*ebiten.Image, 5),
		rubbleImages:    make([]*ebiten.Image, 5),
	}

	for i := 0; i < 5; i++ {
		aFilename := fmt.Sprintf("./resources/images/asteroid%d.png", i)
		asteroidImage, err := g.LoadImage(aFilename)
		if err != nil {
			log.Errorf("error when loading image from file: %s", err.Error())
		}
		g.asteroidImages[i] = asteroidImage
		rFilename := fmt.Sprintf("./resources/images/rubble%d.png", i)
		rubbleImage, err := g.LoadImage(rFilename)
		if err != nil {
			log.Errorf("error when loading image from file: %s", err.Error())
		}
		g.rubbleImages[i] = rubbleImage
	}
	starshipImage, err := g.LoadImage("./resources/images/ship.png")
	if err != nil {
		log.Errorf("error when loading image from file: %s", err.Error())
	}
	g.starshipImage = starshipImage
	bulletImage, err := g.LoadImage("./resources/images/bullet.png")
	if err != nil {
		log.Errorf("error when loading image from file: %s", err.Error())
	}
	g.bulletImage = bulletImage
	boidImage, err := g.LoadImage("./resources/images/boid.png")
	if err != nil {
		log.Errorf("error when loading image from file: %s", err.Error())
	}
	g.boidImage = boidImage
	return g
}

// LoadImage loads a picture into an ebiten image.
func (g *Game) LoadImage(file string) (*ebiten.Image, error) {
	newImage := ebiten.NewImage(int(g.ScreenWidth), int(g.ScreenHeight))

	f, err := os.Open(file)
	if err != nil {
		newImage.Fill(color.White)
		return newImage, errors.New("error when opening file " + err.Error())
	}

	defer f.Close()
	rawImage, _, err := image.Decode(f)
	if err != nil {
		newImage.Fill(color.White)
		return newImage, errors.New("error when decoding image from file " + err.Error())
	}

	newImageFromImage := ebiten.NewImageFromImage(rawImage)
	if newImageFromImage == nil {
		return newImage, errors.New("error when creating image from raw " + err.Error())
	}
	newImage.DrawImage(newImageFromImage, nil)
	return newImage, nil
}

// StartGame initializes a new game.
func (g *Game) StartGame() {
	// add starship
	p := agents.NewStarship(
		g.log,
		g.ScreenWidth/2,
		g.ScreenHeight/2,
		g.ScreenWidth,
		g.ScreenHeight,
		g.Register,
		g.Unregister,
		g.starshipImage,
		g.bulletImage,
		g.debug)
	g.log.Debugf("added starship: %s", p.ID())
	g.Register(p)

	// add asteroids
	for i := 0; i < g.nAsteroids; i++ {
		g.AddAsteroid(g.asteroidImages[rand.Intn(5)])
	}

	// add boids
	for i := 0; i < g.nBoids; i++ {
		g.AddBoid()
	}

	g.startTime = time.Now()
	g.gameDuration = 0
	g.gameOver = false
	g.gameWon = false
	g.kills = 0
}

// AddAsteroid insert a new asteroid in the game.
func (g *Game) AddAsteroid(asteroidImage *ebiten.Image) {
	a := agents.NewAsteroid(g.log,
		float64(rand.Intn(int(g.ScreenWidth))),
		float64(rand.Intn(int(g.ScreenHeight/4))),
		g.ScreenWidth, g.ScreenHeight,
		g.Register, g.Unregister,
		asteroidImage,
		g.rubbleImages,
		g.debug)
	g.Register(a)
}

// AddAsteroid insert a new asteroid in the game.
func (g *Game) AddBoid() {
	b := ai.NewBoid(g.log,
		float64(rand.Intn(int(g.ScreenWidth))),
		float64(rand.Intn(int(g.ScreenHeight/4))),
		g.ScreenWidth, g.ScreenHeight,
		g.boidImage,
		g.Vision,
		g.debug)
	g.Register(b)
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
	for k := range g.boids {
		delete(g.boids, k)
	}

	g.StartGame()
}

// Vision returns all agents located in a radius from (x,y)
func (g *Game) Vision(x, y, radius float64) []physics.Physic {
	nearestAgents := []physics.Physic{}

	for _, v := range g.starships {
		if (v.Dimension().X-x)*(v.Dimension().X-x)+(v.Dimension().Y-y)*(v.Dimension().Y-y) < radius*radius {
			nearestAgents = append(nearestAgents, v)
		}
	}
	for _, v := range g.asteroids {
		if (v.Dimension().X-x)*(v.Dimension().X-x)+(v.Dimension().Y-y)*(v.Dimension().Y-y) < radius*radius {
			nearestAgents = append(nearestAgents, v)
		}
	}
	for _, v := range g.bullets {
		if (v.Dimension().X-x)*(v.Dimension().X-x)+(v.Dimension().Y-y)*(v.Dimension().Y-y) < radius*radius {
			nearestAgents = append(nearestAgents, v)
		}
	}
	for _, v := range g.boids {
		if (v.Dimension().X-x)*(v.Dimension().X-x)+(v.Dimension().Y-y)*(v.Dimension().Y-y) < radius*radius {
			nearestAgents = append(nearestAgents, v)
		}
	}

	return nearestAgents
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
	case physics.BoidAgent:
		g.boids[agent.ID()] = agent
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
