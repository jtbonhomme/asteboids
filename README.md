# asteboids

Asteboids is an 'asteroids'-like game.

![](screen1.png)
![](screen2.png)

# A little bit of history

[Asteroids](https://en.wikipedia.org/wiki/Asteroids_(video_game)) is a arcade game released in 1979 by [Atari, Inc](https://en.wikipedia.org/wiki/Atari,_Inc.). The player controls a single spaceship in an asteroid field. The object of the game is to shoot and destroy the asteroids, while not colliding with either. The game becomes harder as the number of asteroids increases.

![](asteroids-by-atari.jpg)

This repository is a simple port of the original game play in Go, while replacing asteroids with boids (autonomous agents).

# Software Design

The 2D game engine is [ebiten](https://ebiten.org/): A dead simple 2D game library for Go.

# Run

```sh
$ make run
```

# Run with debug information

```sh
$ make debug
```

![](screen3.png)

## Keyboard control

* `key up`: startship move forward
* `key left`: startship rotate counter clockwise
* `key right`: startship rotate clockwise
* `space`: startship shot
* `enter`: game restart
* `s`: takes a screenshot (file is stored as `screenshot_<date><time>.png`)
* `cmd+q`: exit

# Resources

## Fonts

Some fonts in this repository are copyright (c) Jakob Fischer at www.pizzadude.dk,  all rights reserved. 
Do not distribute without the author's permission.
Use these font for non-commercial use only! If you plan to use them for commercial purposes, contact the author before doing so!
For more original fonts take a look at www.pizzadude.dk
