package main

import (
	"log"
	"math"
)

func possibleBoardMoves(board Board) []int {
	var moves []int
	for j := range board[0] {
		if board.IsColumnPlaceable(j) {
			moves = append(moves, j)
		}
	}
	return moves
}

func minimax(board Board, depth int, isPlayerTurn bool) (_score float64, _column int) {
	const (
		winningScore = 1.0
		noMove       = -1
	)

	if depth == 0 {
		return 0.0, noMove
	}
	if board.CheckWinner() != 0 {
		return winningScore, noMove
	}

	var (
		moves  = possibleBoardMoves(board)
		score  = -1.0
		column = -1
	)

	if isPlayerTurn {
		score = math.Inf(1)
		for column := range moves {
			// Execute the move as player
			nextBoard, err := board.Place(column, HUMAN)
			if err != nil {
				log.Fatal(err)
			}

			// Simulate from opponent POV and determine outcome desirability
			// Higher score is desirable
			nextScore, nextColumn := minimax(nextBoard, depth-1, false)
			if nextScore > score {
				score = nextScore
				column = nextColumn
			}
		}

	} else { // Opponent turn
		score = math.Inf(-1)
		for column := range moves {
			// Execute the move as opponent
			nextBoard, err := board.Place(column, AI)
			if err != nil {
				log.Fatal(err)
			}

			// Simulate from player POV and determine outcome desirability
			// Lower score is more desirable
			nextScore, nextColumn := minimax(nextBoard, depth-1, true)
			if nextScore < score {
				score = nextScore
				column = nextColumn
			}
		}
	}

	return score, column
}
