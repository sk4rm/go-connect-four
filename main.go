package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed img/background.png
var _bg []byte

//go:embed img/bluebox.png
var _bluebox []byte

//go:embed img/select.png
var _selected []byte

//go:embed img/zero.png
var _zero []byte

//go:embed img/one.png
var _one []byte

//go:embed img/two.png
var _two []byte

var bg, bluebox, selected, zero, one, two *ebiten.Image

func init() {
	bgDecoded, _, err := image.Decode(bytes.NewReader(_bg))
	if err != nil {
		log.Fatal(err)
	}
	bg = ebiten.NewImageFromImage(bgDecoded)

	blueboxDecoded, _, err := image.Decode(bytes.NewReader(_bluebox))
	if err != nil {
		log.Fatal(err)
	}
	bluebox = ebiten.NewImageFromImage(blueboxDecoded)

	selectedDecoded, _, err := image.Decode(bytes.NewReader(_selected))
	if err != nil {
		log.Fatal(err)
	}
	selected = ebiten.NewImageFromImage(selectedDecoded)

	zeroDecoded, _, err := image.Decode(bytes.NewReader(_zero))
	if err != nil {
		log.Fatal(err)
	}
	zero = ebiten.NewImageFromImage(zeroDecoded)

	oneDecoded, _, err := image.Decode(bytes.NewReader(_one))
	if err != nil {
		log.Fatal(err)
	}
	one = ebiten.NewImageFromImage(oneDecoded)

	twoDecoded, _, err := image.Decode(bytes.NewReader(_two))
	if err != nil {
		log.Fatal(err)
	}
	two = ebiten.NewImageFromImage(twoDecoded)
}

type Game struct {
	board      Board
	cells      []*Cell
	player     int
	isGameOver bool
}

func (g *Game) Update() error {
	cx, cy := ebiten.CursorPosition()
	var selection *Cell
	for _, cell := range g.cells {
		if g.board[cell.i][cell.j] == -1 {
			g.board[cell.i][cell.j] = 0
		}
		if cell.At(float64(cx), float64(cy)) {
			selection = cell
		}
	}

	if selection != nil && !g.isGameOver && g.board.IsColumnPlaceable(selection.j) {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			nextBoard, err := g.board.Place(selection.j, g.player)
			if err != nil {
				return err
			}
			g.board = nextBoard

			if g.board.CheckWinner() == g.player {
				g.isGameOver = true
			} else {
				g.player %= 2
				g.player += 1
			}
		} else {

			var err error
			g.board, err = g.board.Place(selection.j, -1)
			if err != nil && !errors.Is(err, InvalidMoveError) && !errors.Is(err, InvalidPlayerError) {
				log.Fatal("Something went wrong with hover selection: ", err)
			}
		}
	}

	return nil
}

func drawBackground(screen *ebiten.Image) {
	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(
		-float64(bg.Bounds().Dx()/2),
		-float64(bg.Bounds().Dy()/2),
	)
	op.GeoM.Scale(0.667, 0.667)
	op.GeoM.Translate(
		float64(screen.Bounds().Dx()/2),
		float64(screen.Bounds().Dy()/2),
	)
	c := colorm.ColorM{}
	c.ChangeHSV(0, 1, 24)
	colorm.DrawImage(screen, bg, c, op)
}

func drawBluebox(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		-float64(bluebox.Bounds().Dx())/2,
		-float64(bluebox.Bounds().Dy())/2,
	)

	scale := float64(2) / 3 * ebiten.Monitor().DeviceScaleFactor()
	op.GeoM.Scale(scale, scale)

	op.GeoM.Translate(
		float64(screen.Bounds().Dx())/2,
		float64(screen.Bounds().Dy())/2,
	)

	screen.DrawImage(bluebox, op)
}

type Cell struct {
	img        *ebiten.Image
	x, y, w, h float64
	i, j       int
}

func (c *Cell) At(x, y float64) bool {
	return x >= c.x && x < c.x+c.w && y >= c.y && y < c.y+c.h
}

func drawCell(screen *ebiten.Image, i, j, player int, game *Game) {
	var cell *ebiten.Image
	switch player {
	case -1:
		cell = selected
	case 0:
		cell = zero
	case 1:
		cell = one
	case 2:
		cell = two
	default:
		log.Printf("no player %v", player)
		cell = zero
	}

	// Set center to image center
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		-float64(cell.Bounds().Dx())/2,
		-float64(cell.Bounds().Dy())/2,
	)

	// Hi-DPI scaling
	scale := ebiten.Monitor().DeviceScaleFactor() * 0.8
	op.GeoM.Scale(scale, scale)

	w := float64(cell.Bounds().Dx()) * scale
	h := float64(cell.Bounds().Dy()) * scale

	// Move to center of window + manual positioning
	sw := screen.Bounds().Dx()
	sh := screen.Bounds().Dy()
	op.GeoM.Translate(float64(sw)/2+50, float64(sh)/2+15)

	// Offset by row, col
	offsetX := float64(j-3) * (w + 1)
	offsetY := float64(i-3) * (h + 1)
	op.GeoM.Translate(offsetX, offsetY)

	c := &Cell{}
	c.img = cell
	c.x, c.y = op.GeoM.Apply(1, 1)
	c.w, c.h = w, h
	c.i, c.j = i, j

	game.cells = append(game.cells, c)

	screen.DrawImage(cell, op)
}

func drawCells(screen *ebiten.Image, game *Game) {
	game.cells = []*Cell{}
	for i, row := range game.board {
		for j, player := range row {
			drawCell(screen, i, j, player, game)
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawBackground(screen)
	drawBluebox(screen)
	drawCells(screen, g)

	if g.isGameOver {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Player %v wins!\nRestart to play again", g.player))
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	s := ebiten.Monitor().DeviceScaleFactor()
	return int(320 * s), int(240 * s)
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Connect Four")

	// Test
	game := &Game{
		NewBoard(6, 7),
		[]*Cell{},
		1,
		false,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
