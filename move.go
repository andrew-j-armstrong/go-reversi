package reversi

import (
	"fmt"
	"regexp"
	"strconv"
)

type Move struct {
	x int
	y int
}

var moveParseRegex = regexp.MustCompile(`\s*\(?\s*(\d+)\s*,\s*(\d+)\s*\)?\s*`)

func ParseMove(moveDescription string) (Move, error) {
	match := moveParseRegex.FindStringSubmatch(moveDescription)

	if match == nil {
		return Move{}, fmt.Errorf("could not parse move: %s", moveDescription)
	}

	x, err := strconv.Atoi(match[1])
	if err != nil {
		return Move{}, fmt.Errorf("unable to parse move: %s (%s)", moveDescription, err)
	}

	y, err := strconv.Atoi(match[2])
	if err != nil {
		return Move{}, fmt.Errorf("unable to parse move: %s (%s)", moveDescription, err)
	}

	move := Move{x: x, y: y}
	if !move.IsValid() {
		return Move{}, fmt.Errorf("invalid move: %s", move.String())
	}

	return move, nil
}

func (move Move) String() string {
	return fmt.Sprintf("(%d, %d)", move.x, move.y)
}

func (move Move) IsValid() bool {
	return move.x >= 0 && move.x < BoardWidth && move.y >= 0 && move.y < BoardHeight
}
