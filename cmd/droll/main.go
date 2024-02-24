package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strings"

	"github.com/dmars8047/droll/pkg/droll"
)

func main() {

	commands := make([]droll.RollCommand, 0)

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

		parseRes, parseErr := droll.ParseRollCommand(args[0])

		if parseErr != nil {
			detParseErr, ok := parseErr.(droll.TokenParsingError)

			if !ok {
				fmt.Print("\nAn unknown error occurred during roll parsing.\n\n")
				return
			}

			fmt.Printf("\n%s\n\n", detParseErr.Details)
			return
		}

		commands = parseRes
	} else {
		commands = append(commands, droll.RollCommand{
			Num:   1,
			Sides: 20,
		})
	}

	total := 0

	for _, command := range commands {
		for range command.Num {
			fmt.Printf("\nRolling a d%d... ", command.Sides)
			res := rollDie(uint8(command.Sides))
			total += res
			fmt.Printf("%d", res)
		}
		fmt.Print("\n")
	}

	fmt.Printf("\nTotal: %d\n\n", total)
}

func rollDie(numSides uint8) int {
	return rand.IntN(int(numSides)) + 1
}
