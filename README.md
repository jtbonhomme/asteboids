# asteboids

[Asteroids](https://en.wikipedia.org/wiki/Asteroids_(video_game)) is a arcade game released in 1979 by [Atari, Inc](https://en.wikipedia.org/wiki/Atari,_Inc.). The player controls a single spaceship in an asteroid field. The object of the game is to shoot and destroy the asteroids, while not colliding with either. The game becomes harder as the number of asteroids increases.

![](asteroids-by-atari.jpg)

This repository is a simple port of the original game play in Go, while replacing asteroids with boids (autonomous agents).

# Software Design

The 2D game engine is [ebiten](https://ebiten.org/): A dead simple 2D game library for Go.

# Run

```sh
$ make run
```

## Keyboard

* `s`: takes a screenshot (file is stored as `screenshot_<date><time>.png`)
