package main

import (
	"fmt"
	"math"
)

func possibleBoardMoves(board Board) []int {
	var moves []int
	for j := range len(board[0]) {
		if board.IsColumnPlaceable(j) {
			moves = append(moves, j)
		}
	}
	return moves
}

func minimax(board Board, depth int, isPlayerTurn bool) (_score float64, _column int, _ error) {
	const (
		WINNING_SCORE = 1.0
		NO_SCORE      = 0.0
		NO_MOVE       = -1
	)

	if depth == 0 {
		return NO_SCORE, NO_MOVE, nil
	}
	if winner := board.CheckWinner(); winner != 0 {
		if winner == HUMAN {
			return -WINNING_SCORE - float64(depth), NO_MOVE, nil
		} else {
			return WINNING_SCORE + float64(depth), NO_MOVE, nil
		}
	}

	var (
		moves = possibleBoardMoves(board)
		score = 0.0
		move  = -1
	)

	if isPlayerTurn {
		score = math.Inf(1)
		for _, column := range moves {
			// Execute the move as human player
			nextBoard, err := board.Place(column, HUMAN)
			if err != nil {
				return NO_SCORE, NO_MOVE, fmt.Errorf("minimax player turn: %v", err)
			}

			// Simulate from AI POV and determine outcome desirability
			// Lower score is desirable because humans want AI to lose
			nextScore, nextColumn, err := minimax(nextBoard, depth-1, false)
			if err != nil {
				return NO_SCORE, NO_MOVE, err
			}
			if nextScore < score {
				score = nextScore
				column = nextColumn
			}
		}

	} else { // AI player's turn
		score = math.Inf(-1)
		for _, column := range moves {
			// Execute the move as AI player
			nextBoard, err := board.Place(column, AI)
			if err != nil {
				return NO_SCORE, NO_MOVE, fmt.Errorf("minimax AI turn: %v", err)
			}

			// Simulate from human POV and determine outcome desirability
			// Higher score is more desirable since AI wants to win
			nextScore, nextColumn, err := minimax(nextBoard, depth-1, true)
			if err != nil {
				return NO_SCORE, NO_MOVE, err
			}
			if nextScore > score {
				score = nextScore
				move = nextColumn
			}
		}
	}

	return score, move, nil
}
