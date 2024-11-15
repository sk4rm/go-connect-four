package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var bg *ebiten.Image

func init() {
	var err error
	bg, _, err = ebitenutil.NewImageFromFile("img/background.png")
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw background image
	op := &colorm.DrawImageOptions{}
	size := bg.Bounds().Size()
	op.GeoM.Translate(-float64(size.X/2), -float64(size.Y/2))
	op.GeoM.Scale(0.667, 0.667)
	c := colorm.ColorM{}
	c.ChangeHSV(0, 1, 24)
	colorm.DrawImage(screen, bg, c, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("%v", ebiten.Monitor().DeviceScaleFactor()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	s := ebiten.Monitor().DeviceScaleFactor()
	return int(320 * s), int(240 * s)
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Connect Four")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
