package reversi

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Turn int

const (
	Draw Turn = iota
	Player1Turn
	Player2Turn
	Player1Won
	Player2Won
)

type PlayerID int

const (
	Player1 PlayerID = 1
	Player2          = 2
)

type GameState struct {
	board         *Board
	turn          Turn
	moveListeners []chan<- Move
}

func (gameState *GameState) GetTurn() Turn {
	return gameState.turn
}

func (gameState *GameState) RegisterMoveListener(moveListener chan<- Move) {
	gameState.moveListeners = append(gameState.moveListeners, moveListener)
}

func (gameState *GameState) IsValidMove(move Move) bool {
	if gameState.IsGameOver() {
		return false
	}

	if !move.IsValid() {
		return false
	}

	if gameState.board[move.y][move.x] != EmptyPiece {
		return false
	}

	var currentPlayerPiece, opponentPiece Piece
	switch gameState.GetTurn() {
	case Player1Turn:
		currentPlayerPiece = Player1Piece
		opponentPiece = Player2Piece
	case Player2Turn:
		currentPlayerPiece = Player2Piece
		opponentPiece = Player1Piece
	default:
		return false
	}

	for yAdjustment := -1; yAdjustment <= 1; yAdjustment++ {
		for xAdjustment := -1; xAdjustment <= 1; xAdjustment++ {
			if xAdjustment == 0 && yAdjustment == 0 {
				continue
			}

			x := move.x + xAdjustment
			y := move.y + yAdjustment

			if x+xAdjustment < 0 || x+xAdjustment >= BoardWidth || y+yAdjustment < 0 || y+yAdjustment >= BoardHeight {
				continue
			}

			if gameState.board[y][x] != opponentPiece {
				continue
			}

			x += xAdjustment
			y += yAdjustment

			loopDone := false
			for !loopDone && x >= 0 && x < BoardWidth && y >= 0 && y < BoardHeight {
				switch gameState.board[y][x] {
				case currentPlayerPiece:
					return true
				case EmptyPiece:
					loopDone = true
				}

				x += xAdjustment
				y += yAdjustment
			}
		}
	}

	return false
}

func (gameState *GameState) GetPossibleMoves() []Move {
	moves := make([]Move, 0, 1)

	if gameState.IsGameOver() {
		return moves
	}

	for y := 0; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			move := Move{x: x, y: y}
			if gameState.IsValidMove(move) {
				moves = append(moves, move)
			}
		}
	}

	return moves
}

func (gameState *GameState) CountPlayerPieces(piece Piece) int {
	return gameState.board.CountPlayerPieces(piece)
}

func (gameState *GameState) String() string {
	output := gameState.board.String()
	switch gameState.turn {
	case Draw:
		output += "GameState Over - Draw!\n"
	case Player1Turn:
		output += "Player 1's turn.\n"
	case Player2Turn:
		output += "Player 2's turn.\n"
	case Player1Won:
		output += fmt.Sprintf("GameState Over - Player 1 Won! (%d vs %d)\n", gameState.CountPlayerPieces(Player1Piece), gameState.CountPlayerPieces(Player2Piece))
	case Player2Won:
		output += fmt.Sprintf("GameState Over - Player 2 Won! (%d vs %d)\n", gameState.CountPlayerPieces(Player1Piece), gameState.CountPlayerPieces(Player2Piece))
	default:
		output += "Invalid Turn!\n"
	}
	return output
}

func (gameState *GameState) Print() {
	print(gameState.String())
}

func (gameState *GameState) IsGameOver() bool {
	return gameState.turn != Player1Turn && gameState.turn != Player2Turn
}

func (gameState *GameState) verifyEndGame() {
	if gameState.turn != Player1Turn && gameState.turn != Player2Turn {
		return
	}

	var haveValidMoves bool

	for y := 0; !haveValidMoves && y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			if gameState.IsValidMove(Move{x: x, y: y}) {
				haveValidMoves = true
				break
			}
		}
	}

	if !haveValidMoves {
		switch gameState.turn {
		case Player1Turn:
			gameState.turn = Player2Turn
		case Player2Turn:
			gameState.turn = Player1Turn
		}

		haveValidMoves = false
		for y := 0; !haveValidMoves && y < BoardHeight; y++ {
			for x := 0; x < BoardWidth; x++ {
				if gameState.IsValidMove(Move{x: x, y: y}) {
					haveValidMoves = true
					break
				}
			}
		}
	}

	if haveValidMoves {
		return
	}

	var player1PieceCount, player2PieceCount int
	for y := 0; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			switch gameState.board[y][x] {
			case Player1Piece:
				player1PieceCount++
			case Player2Piece:
				player2PieceCount++
			}
		}
	}

	if player1PieceCount == player2PieceCount {
		gameState.turn = Draw
	} else if player1PieceCount > player2PieceCount {
		gameState.turn = Player1Won
	} else {
		gameState.turn = Player2Won
	}
}

func (gameState *GameState) MakeMove(move Move) error {
	if !gameState.IsValidMove(move) {
		return errors.New("Invalid Move!")
	}

	var currentPlayerPiece, opponentPiece Piece
	switch gameState.GetTurn() {
	case Player1Turn:
		currentPlayerPiece = Player1Piece
		opponentPiece = Player2Piece
	case Player2Turn:
		currentPlayerPiece = Player2Piece
		opponentPiece = Player1Piece
	}

	for yAdjustment := -1; yAdjustment <= 1; yAdjustment++ {
		for xAdjustment := -1; xAdjustment <= 1; xAdjustment++ {
			if xAdjustment == 0 && yAdjustment == 0 {
				continue
			}

			x := move.x + xAdjustment
			y := move.y + yAdjustment

			if x+xAdjustment < 0 || x+xAdjustment >= BoardWidth || y+yAdjustment < 0 || y+yAdjustment >= BoardHeight {
				continue
			}

			if gameState.board[y][x] != opponentPiece {
				continue
			}

			x += xAdjustment
			y += yAdjustment

			loopDone := false
			for !loopDone && x >= 0 && x < BoardWidth && y >= 0 && y < BoardHeight {
				switch gameState.board[y][x] {
				case currentPlayerPiece:
					x -= xAdjustment
					y -= yAdjustment
					for x != move.x || y != move.y {
						gameState.board[y][x] = currentPlayerPiece
						x -= xAdjustment
						y -= yAdjustment
					}
					loopDone = true
				case EmptyPiece:
					loopDone = true
				}

				if loopDone {
					break
				}

				x += xAdjustment
				y += yAdjustment
			}
		}
	}

	gameState.board[move.y][move.x] = currentPlayerPiece

	if gameState.turn == Player1Turn {
		gameState.turn = Player2Turn
	} else if gameState.turn == Player2Turn {
		gameState.turn = Player1Turn
	}

	gameState.verifyEndGame()

	for _, moveListener := range gameState.moveListeners {
		moveListener <- move
	}

	if gameState.IsGameOver() {
		for _, moveListener := range gameState.moveListeners {
			close(moveListener)
		}
	}

	return nil
}

func NewGame() *GameState {
	return &GameState{NewGameBoard(), Player1Turn, nil}
}

func (gameState *GameState) Clone() *GameState {
	return &GameState{gameState.board.Clone(), gameState.turn, nil}
}

/* GameState File Format:7x6 array of (RY )
 * Turn is determined by count of R vs Y
 */

func (gameState *GameState) Save(filename string) error {
	f, err := os.Create(filename)

	if err != nil {
		return err
	}

	defer f.Close()

	f.WriteString(gameState.String())
	return nil
}

func ParseGame(gameDescription string) (*GameState, error) {
	gameState := NewGame()

	player1PieceCount := 0
	player2PieceCount := 0

	y := 0
	x := 0
	expectingPiece := false
	for _, c := range gameDescription {
		if !expectingPiece {
			if c == '|' {
				expectingPiece = true
			}
		} else {
			switch c {
			case 'B':
				gameState.board[y][x] = Player1Piece
				player1PieceCount++
				expectingPiece = false
			case 'W':
				gameState.board[y][x] = Player2Piece
				player2PieceCount++
				expectingPiece = false
			case '|':
				gameState.board[y][x] = EmptyPiece
			case '\n':
				expectingPiece = false
				continue
			default:
				continue
			}

			if x < BoardWidth-1 {
				x++
			} else {
				y++
				x = 0
			}

			if y >= BoardHeight {
				break
			}
		}
	}

	if player1PieceCount == player2PieceCount {
		gameState.turn = Player1Turn
	} else if player1PieceCount == player2PieceCount+1 {
		gameState.turn = Player2Turn
	} else {
		gameState.Print()
		return nil, errors.New(fmt.Sprintf("invalid gameState description: (%d red pieces, %d yellow pieces)", player1PieceCount, player2PieceCount))
	}

	gameState.verifyEndGame()

	return gameState, nil
}

func LoadGame(filename string) (*GameState, error) {
	gameDescriptionBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return ParseGame(string(gameDescriptionBytes))
}
