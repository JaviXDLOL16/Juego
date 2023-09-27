package models

import (
	"time"

	"github.com/faiface/pixel"
)

type AnimatedSprite struct {
	Pictures     []*pixel.Sprite
	Delay        []time.Duration
	CurrentFrame int
	ElapsedTime  time.Duration
}
