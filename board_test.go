package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewBoard(t *testing.T) {
	const (
		HEIGHT = 6
		WIDTH  = 7
	)

	board := NewBoard(HEIGHT, WIDTH)

	if len(board) != HEIGHT {
		t.Errorf("Board height does not match")
	}
	if len(board[0]) != WIDTH {
		t.Errorf("Board width does not match")
	}
}

func TestBoard_At(t *testing.T) {
	x := 1
	board := Board([][]int{
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, x, 0},
	})

	if board.At(0, 5) != x {
		t.Errorf("Board At should return player at given x, y position")
	}
}

func TestBoard_Place(t *testing.T) {
	board := NewBoard(6, 7)
	x := 1

	nextBoard, err := board.Place(5, x)
	if err != nil {
		t.Error(err)
	}

	expected := Board([][]int{
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, x, 0},
	})

	if reflect.DeepEqual(nextBoard, board) {
		t.Errorf("Place should not affect original board")
	}

	if !reflect.DeepEqual(nextBoard, expected) {
		t.Errorf("Board placed wrongly\nExpected: %v\nActual: %v", expected, board)
	}
}

func TestBoard_PlaceInvalidMove(t *testing.T) {
	// Players x and y
	x, y := 1, 2

	board := Board([][]int{
		{0, 0, x, 0, 0, 0, 0},
		{0, 0, x, 0, 0, 0, 0},
		{0, 0, x, 0, 0, 0, 0},
		{0, 0, x, 0, 0, 0, 0},
		{0, 0, x, 0, 0, y, 0},
		{0, 0, x, 0, 0, x, 0},
	})

	_, err := board.Place(2, x)
	if !errors.Is(err, InvalidMoveError) {
		t.Errorf("Placing on a full column should return InvalidMoveError")
	}
}

func TestBoard_PlaceInvalidPlayer(t *testing.T) {
	board := NewBoard(6, 7)

	_, err := board.Place(2, 0)

	if !errors.Is(err, InvalidPlayerError) {
		t.Errorf("Placing as player 0 (reserved as empty) should return InvalidPlayerError")
	}

}

func TestBoard_CheckWinner(t *testing.T) {
	var board Board

	x := 1

	board = Board([][]int{
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
	})
	if board.CheckWinner() != 0 {
		t.Errorf("CheckWinner should return 0 when no winner is detected")
	}

	board = Board([][]int{
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, x, x, x, x, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
	})
	if board.CheckWinner() != x {
		t.Errorf("Player should win on 4 consecutive horizontal pieces")
	}

	board = Board([][]int{
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, x, 0, 0, 0},
		{0, 0, 0, x, 0, 0, 0},
		{0, 0, 0, x, 0, 0, 0},
		{0, 0, 0, x, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
	})
	if board.CheckWinner() != x {
		t.Errorf("Player should win on 4 consecutive vertical pieces")
	}

	board = Board([][]int{
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, x, 0, 0, 0},
		{0, 0, x, 0, 0, 0, 0},
		{0, x, 0, 0, 0, 0, 0},
		{x, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
	})
	if board.CheckWinner() != x {
		t.Errorf("Player should win on 4 consecutive diagonal (/) pieces")
	}

	board = Board([][]int{
		{x, 0, 0, 0, 0, 0, 0},
		{0, x, 0, 0, 0, 0, 0},
		{0, 0, x, 0, 0, 0, 0},
		{0, 0, 0, x, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
	})
	if board.CheckWinner() != x {
		t.Errorf("Player should win on 4 consecutive diagonal (\\) pieces")
	}
}
