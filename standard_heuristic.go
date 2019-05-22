package reversi

type StandardHeuristic struct {
	targetPlayer PlayerID
}

func NewStandardHeuristic(targetPlayer PlayerID) *StandardHeuristic {
	return &StandardHeuristic{targetPlayer}
}

func (heuristic *StandardHeuristic) Heuristic(gameState *GameState) float64 {
	if gameState.turn == Draw {
		return 0.0
	} else if gameState.turn == Player1Won {
		if heuristic.targetPlayer == Player1 {
			return 1.0
		} else {
			return -1.0
		}
	} else if gameState.turn == Player2Won {
		if heuristic.targetPlayer == Player1 {
			return -1.0
		} else {
			return 1.0
		}
	}

	player1LockedPiecesBoard := [BoardHeight][BoardWidth]bool{}
	nextMaxWidth := BoardWidth
	for y := 0; y < BoardHeight && nextMaxWidth > 0; y++ {
		for x := 0; x < nextMaxWidth; x++ {
			if gameState.board[y][x] != Player1Piece {
				nextMaxWidth = x
				break
			} else {
				player1LockedPiecesBoard[y][x] = true
			}
		}
	}
	nextMaxWidth = -1
	for y := 0; y < BoardHeight && nextMaxWidth < BoardWidth-1; y++ {
		for x := BoardWidth - 1; x > nextMaxWidth; x-- {
			if gameState.board[y][x] != Player1Piece {
				nextMaxWidth = x
				break
			} else {
				player1LockedPiecesBoard[y][x] = true
			}
		}
	}
	nextMaxWidth = BoardWidth
	for y := BoardHeight - 1; y >= 0 && nextMaxWidth > 0; y-- {
		for x := 0; x < nextMaxWidth; x++ {
			if gameState.board[y][x] != Player1Piece {
				nextMaxWidth = x
				break
			} else {
				player1LockedPiecesBoard[y][x] = true
			}
		}
	}
	nextMaxWidth = -1
	for y := BoardHeight - 1; y >= 0 && nextMaxWidth < BoardWidth-1; y-- {
		for x := BoardWidth - 1; x > nextMaxWidth; x-- {
			if gameState.board[y][x] != Player1Piece {
				nextMaxWidth = x
				break
			} else {
				player1LockedPiecesBoard[y][x] = true
			}
		}
	}

	player2LockedPiecesBoard := [BoardHeight][BoardWidth]bool{}
	nextMaxWidth = BoardWidth
	for y := 0; y < BoardHeight && nextMaxWidth > 0; y++ {
		for x := 0; x < nextMaxWidth; x++ {
			if gameState.board[y][x] != Player2Piece {
				nextMaxWidth = x
				break
			} else {
				player2LockedPiecesBoard[y][x] = true
			}
		}
	}
	nextMaxWidth = -1
	for y := 0; y < BoardHeight && nextMaxWidth < BoardWidth-1; y++ {
		for x := BoardWidth - 1; x > nextMaxWidth; x-- {
			if gameState.board[y][x] != Player2Piece {
				nextMaxWidth = x
				break
			} else {
				player2LockedPiecesBoard[y][x] = true
			}
		}
	}
	nextMaxWidth = BoardWidth
	for y := BoardHeight - 1; y >= 0 && nextMaxWidth > 0; y-- {
		for x := 0; x < nextMaxWidth; x++ {
			if gameState.board[y][x] != Player2Piece {
				nextMaxWidth = x
				break
			} else {
				player2LockedPiecesBoard[y][x] = true
			}
		}
	}
	nextMaxWidth = -1
	for y := BoardHeight - 1; y >= 0 && nextMaxWidth < BoardWidth-1; y-- {
		for x := BoardWidth - 1; x > nextMaxWidth; x-- {
			if gameState.board[y][x] != Player2Piece {
				nextMaxWidth = x
				break
			} else {
				player2LockedPiecesBoard[y][x] = true
			}
		}
	}

	player1LockedPieceCount := 0
	player2LockedPieceCount := 0

	var player1Viability int
	var player2Viability int

	for y := 0; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			if player1LockedPiecesBoard[y][x] {
				player1LockedPieceCount++
			} else if player2LockedPiecesBoard[y][x] {
				player2LockedPieceCount++
			} else { /*
					locationValue := 1
					if y == 0 || y == BoardHeight-1 {
						locationValue *= 8
					}
					if x == 0 || x == BoardWidth-1 {
						locationValue *= 8
					}
					switch gameState.board[y][x] {
					case Player1Piece:
						player1Viability += locationValue
					case Player2Piece:
						player2Viability += locationValue
					}*/
			}
		}
	}

	currentPlayerMoveCount := 0
	for y := 0; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			if gameState.IsValidMove(Move{x: x, y: y}) {
				currentPlayerMoveCount++
			}
		}
	}

	if gameState.GetTurn() == Player1Turn {
		player1Viability += 3 * currentPlayerMoveCount
	} else if gameState.GetTurn() == Player2Turn {
		player2Viability += 3 * currentPlayerMoveCount
	}

	viability := 2.0*(((1.0+float64(player1Viability-player2Viability)/float64(100+player1Viability+player2Viability))*float64(BoardWidth*BoardHeight-player1LockedPieceCount-player2LockedPieceCount)/2.0+float64(player1LockedPieceCount))/float64(BoardWidth*BoardHeight)) - 1.0

	if heuristic.targetPlayer == Player1 {
		return viability
	} else {
		return -viability
	}
}
