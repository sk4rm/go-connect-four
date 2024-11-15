package main

import (
	"errors"
)

type Board [][]int

func NewBoard(height int, width int) Board {
	board := make(Board, height)
	for i := range board {
		board[i] = make([]int, width)
	}
	return board
}

var InvalidMoveError = errors.New("invalid move")
var InvalidPlayerError = errors.New("invalid player")

func (board Board) Place(column int, player int) error {
	if player == 0 {
		return InvalidPlayerError
	}

	// Start from the bottom row of board
	row := len(board) - 1
	for row >= 0 {
		if board[row][column] == 0 {
			board[row][column] = player
			return nil
		}
		row--
	}
	return InvalidMoveError
}

func allEqual(values ...int) bool {
	if len(values) == 0 {
		return true
	}
	for _, v := range values[1:] {
		if v != values[0] {
			return false
		}
	}
	return true
}

func (board Board) CheckWinner() int {
	for i, row := range board {
		for j, player := range row {
			// Skip if empty
			if player == 0 {
				continue
			}

			height, width := len(board), len(board[0])

			// Check horizontal
			if j+3 < width {
				if allEqual(player, row[j+1], row[j+2], row[j+3]) {
					return player
				}
			}

			// Check vertical
			if i+3 < height {
				if allEqual(player, board[i+1][j], board[i+2][j], board[i+3][j]) {
					return player
				}
			}

			// Check diagonal /
			if i+3 < width && j-3 >= 0 {
				if allEqual(player, board[i+1][j-1], board[i+2][j-2], board[i+3][j-3]) {
					return player
				}
			}

			// Check diagonal \
			if i+3 < width && j+3 <= height {
				if allEqual(player, board[i+1][j+1], board[i+2][j+2], board[i+3][j+3]) {
					return player
				}
			}
		}
	}
	return 0
}
