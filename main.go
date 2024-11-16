package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	//go:embed img/background.png
	pngBg []byte

	//go:embed img/bluebox.png
	pngBlueBox []byte

	//go:embed img/select.png
	pngSelect []byte

	//go:embed img/zero.png
	pngZero []byte

	//go:embed img/one.png
	pngOne []byte

	//go:embed img/two.png
	pngTwo []byte

	bg,
	blueBox,
	selected,
	zero,
	one,
	two *ebiten.Image
)

const (
	PREVIEW = -1
	EMPTY   = 0
	HUMAN   = 1
	AI      = 2
)

func preloadImage(src []byte) *ebiten.Image {
	decoded, _, err := image.Decode(bytes.NewReader(src))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(decoded)
}

func init() {
	bg = preloadImage(pngBg)
	blueBox = preloadImage(pngBlueBox)
	selected = preloadImage(pngSelect)
	zero = preloadImage(pngZero)
	one = preloadImage(pngOne)
	two = preloadImage(pngTwo)
}

type Game struct {
	board      Board
	cells      []*Cell
	player     int
	isGameOver bool
}

func switchPlayer(player int) (int, error) {
	switch player {
	case HUMAN:
		return AI, nil
	case AI:
		return HUMAN, nil
	default:
		return EMPTY, fmt.Errorf("invalid player %v", player)
	}
}

func (g *Game) Update() error {
	cx, cy := ebiten.CursorPosition()
	var selection *Cell
	for _, cell := range g.cells {
		if g.board[cell.i][cell.j] == PREVIEW {
			g.board[cell.i][cell.j] = EMPTY
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
				g.player, err = switchPlayer(g.player)
				if err != nil {
					return err
				}
			}

		} else {
			var err error
			g.board, err = g.board.Place(selection.j, PREVIEW)
			if err != nil && !errors.Is(err, InvalidMoveError) && !errors.Is(err, InvalidPlayerError) {
				return fmt.Errorf("invalid hover selection: %v", err)
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
		-float64(blueBox.Bounds().Dx())/2,
		-float64(blueBox.Bounds().Dy())/2,
	)

	scale := float64(2) / 3 * ebiten.Monitor().DeviceScaleFactor()
	op.GeoM.Scale(scale, scale)

	op.GeoM.Translate(
		float64(screen.Bounds().Dx())/2,
		float64(screen.Bounds().Dy())/2,
	)

	screen.DrawImage(blueBox, op)
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
	case PREVIEW:
		cell = selected
	case EMPTY:
		cell = zero
	case HUMAN:
		cell = one
	case AI:
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
		HUMAN,
		false,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
