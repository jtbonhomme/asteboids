package game

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jtbonhomme/asteboids/internal/physics"
	"github.com/jtbonhomme/asteboids/internal/sounds"
)

// UpdateAgents loops over all game agents to update them
func (g *Game) UpdateAgents() {
	for _, b := range g.bullets {
		b.Update()
	}
	for _, a := range g.asteroids {
		a.Update()
	}
	for _, s := range g.starships {
		s.Update()
	}
	for _, b := range g.boids {
		b.Update()
	}
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Update the agents
	g.UpdateAgents()

	// detect starship collision with asteroids
	for _, starship := range g.starships {
		_, ok := starship.IntersectMultiple(g.asteroids)
		if ok {
			starship.Explode()
		}
	}

	// detect asteroid collision with bullet
	for _, asteroid := range g.asteroids {
		bID, ok := asteroid.IntersectMultiple(g.bullets)
		if ok {
			asteroid.Explode()
			asteroidType := asteroid.Type()
			delete(g.bullets, bID)
			g.kills++
			// Only add a new asteroids if the destroyed agent is also an asteroid (not a rubble)
			if asteroidType == physics.AsteroidAgent {
				g.AddAsteroid(g.asteroidImages[rand.Intn(5)])
			}
		}
	}

	// game ends when there is no starship left
	if len(g.starships) == 0 {
		g.gameOver = true
		g.gameWon = false
	}

	// update time until game ends
	if !g.gameOver {
		g.gameDuration = time.Since(g.startTime).Round(time.Second)
	}

	// periodically add new asteroids
	if g.conf.AsteroidsRespawn > 0 && int(g.gameDuration.Seconds()/g.conf.AsteroidsRespawn) > len(g.asteroids) {
		g.AddAsteroid(g.asteroidImages[rand.Intn(5)])
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		err := g.Dump()
		if err != nil {
			g.log.Errorf("can't dump: %s", err.Error())
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyM) {
		if g.mute {
			sounds.Unmute()
			g.mute = false
		} else {
			sounds.Mute()
			g.mute = true
		}
	}

	return nil
}

// Dump saves internal game state in a file.
func (g *Game) Dump() error {
	var err error
	const datetimeFormat = "20060102030405000"

	now := time.Now()
	name := fmt.Sprintf("asteboids_%s.dump", now.Format(datetimeFormat))
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, a := range g.starships {
		err = a.Dump(f)
		if err != nil {
			return err
		}
	}
	for _, a := range g.asteroids {
		err = a.Dump(f)
		if err != nil {
			return err
		}
	}
	for _, a := range g.bullets {
		err = a.Dump(f)
		if err != nil {
			return err
		}
	}
	for _, a := range g.boids {
		err = a.Dump(f)
		if err != nil {
			return err
		}
	}
	g.log.Infof("Saved dump: %s", name)
	return err
}
