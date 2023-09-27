package scenes

import (
	"fmt"
	"image"
	"math"
	"os"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

type GameState int

const (
	GameStateMenu GameState = iota
	GameStatePlaying
	GameStatePaused
	GameStateWon
	GameStateLost
)

var ballPos, playerPos, ballDirection pixel.Vec
var counter int
var gameState GameState
var ballSpeed float64
var isPaused bool
var buttons map[string]pixel.Rect
var playerPic, backgroundPic, ballPic *AnimatedSprite

type AnimatedSprite struct {
	pictures     []*pixel.Sprite
	delay        []time.Duration
	currentFrame int
	elapsedTime  time.Duration
}

func Run() {
	winCfg := pixelgl.WindowConfig{
		Title:  "Juego de Esquivar",
		Bounds: pixel.R(0, 0, 800, 600),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(winCfg)
	if err != nil {
		panic(err)
	}

	loadPictures()
	initGame()

	go playerRoutine(win)
	go ballRoutine(win)
	go counterRoutine()

	for !win.Closed() {
		handleInput(win)

		drawBackground(win)
		drawPlayer(win)
		drawBall(win)
		drawUI(win)
		drawButtons(win)
		win.Update()
		time.Sleep(time.Millisecond * 16)
	}
}

func initGame() {
	gameState = GameStateMenu
	ballPos = pixel.V(100, 300)
	playerPos = pixel.V(400, 300)
	counter = 30
	ballSpeed = 6
	ballDirection = pixel.V(1, 1)
	isPaused = false
	buttons = map[string]pixel.Rect{
		"start":   pixel.R(320, 480, 480, 520),
		"restart": pixel.R(320, 480, 480, 520),
		"exit":    pixel.R(320, 400, 480, 440),
	}
}

func playerRoutine(win *pixelgl.Window) {
	speed := 4.0
	for {
		if !isPaused && gameState == GameStatePlaying {
			if win.Pressed(pixelgl.KeyA) && playerPos.X > 50 {
				playerPos.X -= speed
			}
			if win.Pressed(pixelgl.KeyD) && playerPos.X < win.Bounds().W()-50 {
				playerPos.X += speed
			}
			if win.Pressed(pixelgl.KeyW) && playerPos.Y < win.Bounds().H()-50 {
				playerPos.Y += speed
			}
			if win.Pressed(pixelgl.KeyS) && playerPos.Y > 50 {
				playerPos.Y -= speed
			}
		}
		time.Sleep(time.Millisecond)
	}
}

func ballRoutine(win *pixelgl.Window) {
	trackingFactor := 0.10

	for {
		if !isPaused && gameState == GameStatePlaying {
			ballSpeed += 0.053
			if counter <= 10 {
				ballSpeed = math.Min(ballSpeed, 30)
			}

			adjustmentDirection := playerPos.Sub(ballPos).Unit().Scaled(trackingFactor)
			overallDirection := ballDirection.Add(adjustmentDirection).Unit()

			ballPos = ballPos.Add(overallDirection.Scaled(ballSpeed))

			if ballPos.X <= 0 || ballPos.X >= win.Bounds().W() {
				ballDirection.X = -ballDirection.X
			}
			if ballPos.Y <= 0 || ballPos.Y >= win.Bounds().H() {
				ballDirection.Y = -ballDirection.Y
			}

			playerRect := pixel.Rect{Min: playerPos.Sub(pixel.V(50, 50)), Max: playerPos.Add(pixel.V(50, 50))}
			if playerRect.Contains(ballPos) {
				gameState = GameStateLost
			}
		}
		time.Sleep(time.Millisecond * 16)
	}
}

func counterRoutine() {
	for {
		time.Sleep(time.Second)
		if !isPaused && gameState == GameStatePlaying {
			counter--
		}
		if counter <= 0 {
			gameState = GameStateWon
			break
		}
	}
}

func handleInput(win *pixelgl.Window) {
	mousePos := win.MousePosition()

	if win.JustPressed(pixelgl.MouseButtonLeft) {
		for btnName, btnRect := range buttons {
			if btnRect.Contains(mousePos) {
				switch btnName {
				case "start":
					if gameState == GameStateMenu {
						gameState = GameStatePlaying
						return
					}
				case "restart":
					if gameState == GameStateLost || gameState == GameStateWon {
						initGame()
					}
				case "exit":
					win.SetClosed(true)
				}
			}
		}
	}

	if gameState == GameStatePlaying && win.JustPressed(pixelgl.KeyP) {
		isPaused = !isPaused
	}
}

func loadPictures() {
	ballPic = loadPicture("./assets/SpikeBallM.png")
	playerPic = loadPicture("./assets/personajeM.png")
	backgroundPic = loadPicture("./assets/HEscenarioM.png")
}

func loadPicture(path string) *AnimatedSprite {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	bounds := img.Bounds()
	pixelBounds := pixel.R(float64(bounds.Min.X), float64(bounds.Min.Y), float64(bounds.Max.X), float64(bounds.Max.Y))
	sprite := pixel.NewSprite(pixel.PictureDataFromImage(img), pixelBounds)

	// Devuelve un AnimatedSprite con solo un frame para imágenes estáticas
	return &AnimatedSprite{
		pictures: []*pixel.Sprite{sprite},
		delay:    []time.Duration{0},
	}
}

func drawBall(win *pixelgl.Window) {
	// 1. Dibuja el hitbox
	imd := imdraw.New(nil)
	//imd.Color = colornames.Black o cualquier otro color que desees -> ver hitbox
	imd.Push(ballPos)
	imd.Circle(40, 0) // Asume que el radio es 50, ajústalo como necesites
	imd.Draw(win)

	// 2. Dibuja el sprite
	ballPic.pictures[ballPic.currentFrame].Draw(win, pixel.IM.Scaled(ballPos, 1).Moved(ballPos))
	ballPic.elapsedTime += time.Millisecond * 16
	if ballPic.elapsedTime >= ballPic.delay[ballPic.currentFrame] {
		ballPic.elapsedTime = 0
		ballPic.currentFrame = (ballPic.currentFrame + 1) % len(ballPic.pictures)
	}
}

func drawPlayer(win *pixelgl.Window) {
	playerPic.pictures[0].Draw(win, pixel.IM.Scaled(playerPos, 1).Moved(playerPos))
}

func drawBackground(win *pixelgl.Window) {
	backgroundPic.pictures[0].Draw(win, pixel.IM.Scaled(pixel.V(400, 300), 1).Moved(pixel.V(400, 300)))
}

func drawUI(win *pixelgl.Window) {
	if gameState == GameStatePlaying {
		counterText := fmt.Sprintf("Tiempo restante: %d", counter)
		txt := text.New(pixel.V(340, win.Bounds().H()-50), text.Atlas7x13)
		txt.Color = colornames.White
		txt.Clear()
		txt.WriteString(counterText)
		txt.Draw(win, pixel.IM)
	}

	if isPaused {
		pauseText := "Juego Pausado. Presiona [P] para continuar."
		txt := text.New(pixel.V(250, win.Bounds().H()-150), text.Atlas7x13)
		txt.Color = colornames.White
		txt.Clear()
		txt.WriteString(pauseText)
		txt.Draw(win, pixel.IM)
	}
}

func drawButtons(win *pixelgl.Window) {
	imd := imdraw.New(nil)
	txt := text.New(pixel.V(0, 0), text.Atlas7x13)
	txt.Color = colornames.White

	if gameState == GameStateMenu {
		imd.Color = colornames.Black
		btnRect := buttons["start"]
		imd.Push(btnRect.Min, btnRect.Max)
		imd.Rectangle(0)

		btnRect = buttons["exit"]
		imd.Push(btnRect.Min, btnRect.Max)
		imd.Rectangle(0)
	}

	if gameState == GameStateLost || gameState == GameStateWon {
		imd.Color = colornames.Black
		btnRect := buttons["restart"]
		imd.Push(btnRect.Min, btnRect.Max)
		imd.Rectangle(0)

		btnRect = buttons["exit"]
		imd.Push(btnRect.Min, btnRect.Max)
		imd.Rectangle(0)
	}

	imd.Draw(win)

	if gameState == GameStateMenu {
		txt.Dot = pixel.V(340, 495)
		txt.WriteString("Iniciar")
		txt.Draw(win, pixel.IM)

		txt.Dot = pixel.V(340, 415)
		txt.WriteString("Salir")
		txt.Draw(win, pixel.IM)
	}

	if gameState == GameStateLost {
		txt.Dot = pixel.V(340, 495)
		txt.WriteString("Reiniciar")
		txt.Draw(win, pixel.IM)

		txt.Dot = pixel.V(340, 415)
		txt.WriteString("Salir")
		txt.Draw(win, pixel.IM)

		txt.Dot = pixel.V(360, win.Bounds().H()-50)
		txt.WriteString("Has perdido!")
		txt.Draw(win, pixel.IM)
	}

	if gameState == GameStateWon {
		txt.Dot = pixel.V(340, 495)
		txt.WriteString("Reiniciar")
		txt.Draw(win, pixel.IM)

		txt.Dot = pixel.V(340, 415)
		txt.WriteString("Salir")
		txt.Draw(win, pixel.IM)

		txt.Dot = pixel.V(360, win.Bounds().H()-50)
		txt.WriteString("Has ganado!")
		txt.Draw(win, pixel.IM)
	}
}
