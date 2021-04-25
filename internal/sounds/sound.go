package sounds

import (
	// import embed to load truetype font
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	sampleRate = 11025
)

var audioContext *audio.Context

//go:embed fire.wav
var fireWAV []byte
var FirePlayer *audio.Player

//go:embed thrust.wav
var thrustWAV []byte
var ThrustPlayer *audio.Player

//go:embed beat1.wav
var beat1WAV []byte
var Beat1Player *audio.Player

//go:embed beat2.wav
var beat2WAV []byte
var Beat2Player *audio.Player

//go:embed bangSmall.wav
var bangSmallWAV []byte
var BangSmallPlayer *audio.Player

//go:embed bangMedium.wav
var bangMediumWAV []byte
var BangMediumPlayer *audio.Player

//go:embed bangLarge.wav
var bangLargeWAV []byte
var BangLargePlayer *audio.Player

func Init() {
	audioContext = audio.NewContext(sampleRate)
	FirePlayer = audio.NewPlayerFromBytes(audioContext, fireWAV)
	ThrustPlayer = audio.NewPlayerFromBytes(audioContext, thrustWAV)
	Beat1Player = audio.NewPlayerFromBytes(audioContext, beat1WAV)
	Beat2Player = audio.NewPlayerFromBytes(audioContext, beat2WAV)
	BangSmallPlayer = audio.NewPlayerFromBytes(audioContext, bangSmallWAV)
	BangMediumPlayer = audio.NewPlayerFromBytes(audioContext, bangMediumWAV)
	BangLargePlayer = audio.NewPlayerFromBytes(audioContext, bangLargeWAV)
}

func Mute() {
	SetVolume(0)
}

func Unmute() {
	SetVolume(1)
}

func SetVolume(v float64) {
	if v < 0 || v > 1 {
		return
	}
	FirePlayer.SetVolume(v)
	ThrustPlayer.SetVolume(v)
	Beat1Player.SetVolume(v)
	Beat2Player.SetVolume(v)
	BangSmallPlayer.SetVolume(v)
	BangMediumPlayer.SetVolume(v)
	BangLargePlayer.SetVolume(v)
}
