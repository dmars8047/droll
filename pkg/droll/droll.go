/*
Package droll provides a dice rolling simulation for role-playing games.

It includes functionality to parse dice roll commands (like "2d6+1d4"), roll dice with varying numbers of sides, and write the results of dice rolls to an io.Writer.

Example usage:

	commands, err := droll.ParseRollTokens("2d6+1d4")
	if err != nil {
	    // handle error
	}

	err = droll.Roll(commands, os.Stdout)
	if err != nil {
	    // handle error
	}

This would parse the command "2d6+1d4", roll two six-sided dice and one four-sided die, and write the results to standard output.
*/
package droll

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// validDieSideNums is a map of valid die sides that can be used in a dice roll.
var validDieSideNums = map[uint64]struct{}{4: {}, 6: {}, 8: {}, 10: {}, 12: {}, 20: {}}

// ErrInvalidRollToken is an error that is returned when a roll token is invalid.
var ErrInvalidRollToken = errors.New("invalid roll token")

// DRollTokenParsingError is an error that is returned when a roll token is invalid.
type DRollTokenParsingError struct {
	Details string
}

// Error returns the error message for the DRollTokenParsingError.
func (parseErr DRollTokenParsingError) Error() string {
	return ErrInvalidRollToken.Error()
}

// ParseRollTokens parses a string of roll commands into a slice of RollTokens.
func ParseRollTokens(rollCommandString string) ([]RollTokens, error) {
	commands := make([]RollTokens, 0)

	if len(rollCommandString) == 0 {
		return nil, DRollTokenParsingError{
			Details: "No command provided.",
		}
	}

	for _, token := range strings.Split(rollCommandString, "+") {
		tokenSplit := strings.Split(token, "d")

		if len(tokenSplit) != 2 {
			return nil, DRollTokenParsingError{
				Details: fmt.Sprintf("Command token not recognized as a valid dice roll: %s.", token),
			}
		}

		var n, s uint64

		n, err := strconv.ParseUint(tokenSplit[0], 10, 8)

		if err != nil || n > 255 {
			return nil, DRollTokenParsingError{
				Details: fmt.Sprintf("'%s' is not a valid number of dice. Acceptable values are greater than 0 and less than 256.", tokenSplit[0]),
			}
		}

		s, err = strconv.ParseUint(tokenSplit[1], 10, 8)

		if err != nil || s > 20 {
			return nil, DRollTokenParsingError{
				Details: fmt.Sprintf("'%s' is not a valid nor allowable number of sides for dice. Accaptable values include: %v.", tokenSplit[1], validDieSideNums),
			}
		}

		_, allowable := validDieSideNums[s]

		if !allowable {
			return nil, DRollTokenParsingError{
				Details: fmt.Sprintf("'%s' is not a valid nor allowable number of sides for dice. Accaptable values include: %v.", tokenSplit[1], validDieSideNums),
			}
		}

		commands = append(commands, RollTokens{
			Num:   uint8(n),
			Sides: uint8(s),
		})
	}

	return commands, nil
}

// RollTokens is a struct that represents a dice roll command.
type RollTokens struct {
	Num   uint8
	Sides uint8
}

// Roll rolls a set of dice and writes the results to a writer.
func Roll(tokens []RollTokens, writer io.Writer) error {
	total := 0

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for _, token := range tokens {
		for i := uint8(0); i < token.Num; i++ {
			res := r.Intn(int(token.Sides)) + 1
			total += res
			_, err := io.WriteString(writer, fmt.Sprintf("\nRolling a d%d... %d", token.Sides, res))

			if err != nil {
				return err
			}
		}

		_, err := io.WriteString(writer, "\n")

		if err != nil {
			return err
		}
	}

	_, err := io.WriteString(writer, fmt.Sprintf("\nTotal: %d", total))

	return err
}
