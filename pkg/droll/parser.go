package droll

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	valid_die_side_nums = []uint64{4, 6, 8, 10, 12, 20}
)

var ErrInvalidRollCommandToken = errors.New("invalid roll command token")

type TokenParsingError struct {
	Details string
}

func (parseErr TokenParsingError) Error() string {
	return ErrInvalidRollCommandToken.Error()
}

func ParseRollCommand(rollCommandString string) ([]RollCommand, error) {
	commands := make([]RollCommand, 0)

	for _, token := range strings.Split(rollCommandString, "+") {
		tokenSplit := strings.Split(token, "d")

		if len(tokenSplit) != 2 {
			return nil, TokenParsingError{
				Details: fmt.Sprintf("Command token not recognized as a valid dice roll: %s.", token),
			}
		}

		var n, s uint64

		n, err := strconv.ParseUint(tokenSplit[0], 10, 8)

		if err != nil || n > 255 {
			return nil, TokenParsingError{
				Details: fmt.Sprintf("'%s' is not a valid number of dice. Acceptable values are greater than 0 and less than 256.", tokenSplit[0]),
			}
		}

		s, err = strconv.ParseUint(tokenSplit[1], 10, 8)

		if err != nil || s > 20 {
			return nil, TokenParsingError{
				Details: fmt.Sprintf("'%s' is not a valid nor allowable number of sides for dice. Accaptable values include: %v.", tokenSplit[1], valid_die_side_nums),
			}
		}

		allowable := false

		for _, allowableDie := range valid_die_side_nums {
			if s == allowableDie {
				allowable = true
			}
		}

		if !allowable {
			return nil, TokenParsingError{
				Details: fmt.Sprintf("'%s' is not a valid nor allowable number of sides for dice. Accaptable values include: %v.", tokenSplit[1], valid_die_side_nums),
			}
		}

		commands = append(commands, RollCommand{
			Num:   uint8(n),
			Sides: uint8(s),
		})
	}

	return commands, nil
}

type RollCommand struct {
	Num   uint8
	Sides uint8
}
