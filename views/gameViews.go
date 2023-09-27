package views

import (
	"fmt"
	"time"

	"main/models"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

func DrawBall(win *pixelgl.Window, ballPic *models.AnimatedSprite) {
	// Dibuja el hitbox
	imd := imdraw.New(nil)
	imd.Push(ballPos)
	imd.Circle(40, 0)
	imd.Draw(win)

	// Dibuja el sprite
	ballPic.Pictures[ballPic.CurrentFrame].Draw(win, pixel.IM.Scaled(ballPos, 1).Moved(ballPos))
	ballPic.ElapsedTime += time.Millisecond * 16
	if ballPic.ElapsedTime >= ballPic.Delay[ballPic.CurrentFrame] {
		ballPic.ElapsedTime = 0
		ballPic.CurrentFrame = (ballPic.CurrentFrame + 1) % len(ballPic.Pictures)
	}
}

func DrawPlayer(win *pixelgl.Window, playerPic *models.AnimatedSprite) {
	playerPic.Pictures[0].Draw(win, pixel.IM.Scaled(playerPos, 1).Moved(playerPos))
}

func DrawBackground(win *pixelgl.Window, backgroundPic *models.AnimatedSprite) {
	backgroundPic.Pictures[0].Draw(win, pixel.IM.Scaled(pixel.V(400, 300), 1).Moved(pixel.V(400, 300)))
}

func DrawUI(win *pixelgl.Window) {
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

func DrawButtons(win *pixelgl.Window) {
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
