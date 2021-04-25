package config

import "github.com/jtbonhomme/conf"

const (
	defaultAsteroids                  int     = 4
	defaultBoids                      int     = 40
	defaultScreenWidth                float64 = 1080
	defaultScreenHeight               float64 = 720
	defaultScoreTimeUnit              float64 = 5
	defaultAutoGenerateAsteroidsRatio float64 = 10
	defaultMaxTPS                     int     = 60
	defaultVisionRadius               float64 = 75
)

type Config struct {
	Debug                      bool    `conf:"debug" help:"Debug log level activated."`
	Optim                      bool    `conf:"optim" help:"Optimized mode activated."`
	CPUProfile                 string  `conf:"cpuprofile" help:"Write CPU profile to file."`
	Asteroids                  int     `conf:"asteroids" help:"Number of asteroids at the start of the game."`
	Boids                      int     `conf:"boids" help:"Number of boids at the start of the game."`
	ScreenWidth                float64 `conf:"screenWidth" help:"Screen width (in pixels)."`
	ScreenHeight               float64 `conf:"screenHeight" help:"Screen height (in pixels)."`
	ScoreTimeUnit              float64 `conf:"scoreTimeUnit" help:"Time delay (in second) to win one point, just being still alive."`
	AutoGenerateAsteroidsRatio float64 `conf:"autoGenerateAsteroidsRatio" help:"Time delay (in second) before a new asteroids spawn."`
	MaxTPS                     int     `conf:"maxTPS" help:"Maximum ticks per second."`
	VisionRadius               float64 `conf:"visionRadius" help:"Radius of the boids vision."`
}

func New() *Config {
	config := &Config{
		Asteroids:                  defaultAsteroids,
		Boids:                      defaultBoids,
		ScreenWidth:                defaultScreenWidth,
		ScreenHeight:               defaultScreenHeight,
		ScoreTimeUnit:              defaultScoreTimeUnit,
		AutoGenerateAsteroidsRatio: defaultAutoGenerateAsteroidsRatio,
		MaxTPS:                     defaultMaxTPS,
		VisionRadius:               defaultVisionRadius,
	}
	conf.Load(config)
	return config
}
