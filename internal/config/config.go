package config

import "github.com/jtbonhomme/conf"

const (
	defaultAsteroids        int     = 4
	defaultBoids            int     = 60
	defaultScreenWidth      float64 = 1080
	defaultScreenHeight     float64 = 720
	defaultScoreTimeUnit    float64 = 5
	defaultAsteroidsRespawn float64 = 10
	defaultMaxTPS           int     = 60
	defaultVisionRadius     float64 = 75
	defaultMute             bool    = true
)

type Config struct {
	Mute             bool    `conf:"mute" help:"Mute sound (default is true)."`
	Debug            bool    `conf:"debug" help:"Debug log level activated (default is false)."`
	Optim            bool    `conf:"optim" help:"Optimized mode activated (default is false)."`
	CPUProfile       string  `conf:"cpuprofile" help:"Write CPU profile to file (default is empty)."`
	Asteroids        int     `conf:"asteroids" help:"Number of asteroids at the start of the game (default is 4)."`
	Boids            int     `conf:"boids" help:"Number of boids at the start of the game (default is 60)."`
	ScreenWidth      float64 `conf:"screenWidth" help:"Screen width (in pixels, default is 1080)."`
	ScreenHeight     float64 `conf:"screenHeight" help:"Screen height (in pixels, default is 720)."`
	ScoreTimeUnit    float64 `conf:"scoreTimeUnit" help:"Time delay (in second) to win one point (default is 5)."`
	AsteroidsRespawn float64 `conf:"asteroidsRespawn" help:"Time delay (in second) before a new asteroids spawn (default is 10)."`
	MaxTPS           int     `conf:"maxTPS" help:"Maximum ticks per second  (default is 60)."`
	VisionRadius     float64 `conf:"visionRadius" help:"Radius (in pixels) of the agents vision (default is 150)."`
}

func New() *Config {
	config := &Config{
		Mute:             defaultMute,
		Asteroids:        defaultAsteroids,
		Boids:            defaultBoids,
		ScreenWidth:      defaultScreenWidth,
		ScreenHeight:     defaultScreenHeight,
		ScoreTimeUnit:    defaultScoreTimeUnit,
		AsteroidsRespawn: defaultAsteroidsRespawn,
		MaxTPS:           defaultMaxTPS,
		VisionRadius:     defaultVisionRadius,
	}
	conf.Load(config)
	return config
}
