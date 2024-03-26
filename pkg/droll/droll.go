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

var ErrInvalidRollCommandToken = errors.New("invalid roll command token")

type DRollTokenParsingError struct {
	Details string
}

func (parseErr DRollTokenParsingError) Error() string {
	return ErrInvalidRollCommandToken.Error()
}

func ParseRollTokens(rollCommandString string) ([]RollTokens, error) {
	commands := make([]RollTokens, 0)

	if len(rollCommandString) == 0 {
		return nil, DRollTokenParsingError{
			Details: "No command provided.",
		}
	}

	valid_die_side_nums := []uint64{4, 6, 8, 10, 12, 20}

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
			return nil, DRollTokenParsingError{
				Details: fmt.Sprintf("'%s' is not a valid nor allowable number of sides for dice. Accaptable values include: %v.", tokenSplit[1], valid_die_side_nums),
			}
		}

		commands = append(commands, RollTokens{
			Num:   uint8(n),
			Sides: uint8(s),
		})
	}

	return commands, nil
}

type RollTokens struct {
	Num   uint8
	Sides uint8
}

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

	_, err := io.WriteString(writer, fmt.Sprintf("\nTotal: %d\n\n", total))

	return err
}
