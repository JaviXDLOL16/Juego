package models

import "github.com/faiface/pixel"

type GameState int

const (
	GameStateMenu GameState = iota
	GameStatePlaying
	GameStatePaused
	GameStateWon
	GameStateLost
)

var (
	BallPos   pixel.Vec
	PlayerPos pixel.Vec
	GameState GameState
	Counter   int
	IsPaused  bool
	Buttons   map[string]pixel.Rect
)
