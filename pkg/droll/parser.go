package droll

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRollCommandToken = errors.New("invalid roll command token")

type TokenParsingError struct {
	Token   string
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
			// fmt.Printf("\nCommand token not recognized as a valid dice roll: %s.\n\nFor help try `droll help`.\n\n", token)
			return nil, TokenParsingError{
				Details: fmt.Sprintf("Command token not recognized as a valid dice roll: %s.", token),
			}
		}

		var n, s uint64

		n, err := strconv.ParseUint(tokenSplit[0], 10, 8)

		if err != nil || n > 255 {
			// fmt.Printf("\n'%s' is not a valid number of dice. Acceptable values are greater than 0 and less than 256.\n\n", tokenSplit[0])
			return nil, TokenParsingError{
				Details: fmt.Sprintf("'%s' is not a valid number of dice. Acceptable values are greater than 0 and less than 256.", tokenSplit[0]),
			}
		}

		s, err = strconv.ParseUint(tokenSplit[1], 10, 8)

		if err != nil || s > 20 {
			// fmt.Printf("\n'%s' is not a valid nor allowable number of sides for dice. Accaptable values are greater than 1 and less than 20.\n\n", tokenSplit[1])
			return nil, TokenParsingError{
				Details: fmt.Sprintf("'%s' is not a valid nor allowable number of sides for dice. Accaptable values are greater than 1 and less than 20.", tokenSplit[1]),
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