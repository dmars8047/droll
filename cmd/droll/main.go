package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dmars8047/droll/pkg/droll"
)

func main() {

	commands := make([]droll.RollTokens, 0)

	args := os.Args[1:]

	if len(args) > 0 {
		arg := strings.ToLower(args[0])

		if arg == "help" || arg == "?" || arg == "/?" || arg == "/help" || arg == "-?" || arg == "-help" {
			fmt.Print(
				"\nName: droll\n\n",
				"Description: droll is a dice rolling simulation program. Without any parameters the program rolls a single d20 (a 20 sided die). ",
				"However, different number/die side combinations can be used when provided as command line arguments.\n\n",
				"Example Usage: `droll 2d6` to roll two six-sided dice or `droll 2d6+1d4` to roll two six-sided dice and one four-sided die.\n\n")

			return
		}

		parseRes, parseErr := droll.ParseRollTokens(args[0])

		if parseErr != nil {
			detParseErr, ok := parseErr.(droll.DRollTokenParsingError)

			if !ok {
				fmt.Print("\nAn unknown error occurred during roll parsing.\n\n")
				return
			}

			fmt.Printf("\n%s\n\n", detParseErr.Details)
			return
		}

		commands = parseRes
	} else {
		commands = append(commands, droll.RollTokens{
			Num:   1,
			Sides: 20,
		})
	}

	droll.Roll(commands, consoleWriter{})
}

// A writer that writes to the console.
type consoleWriter struct{}

func (cw consoleWriter) Write(p []byte) (n int, err error) {
	return fmt.Print(string(p))
}
