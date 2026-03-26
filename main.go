package main

import (
	"os"
	"slices"
	"strings"

	"hartool/cmd"
)

var knownCmds = []string{
	"list",
	"req-headers",
	"resp-headers",
	"headers",
	"req-body",
	"resp-body",
	"replay",
	"search",
	"exchange",
	"mcp",
	"help",
	"completion",
}

func main() {
	// Count how many args are known command words.
	// If exactly one is found, move it to front so cobra can dispatch normally.
	// If two or more are found (e.g. a file named after a command), leave args
	// as-is and let cobra handle it - the command word must already be first.
	found := false
	matches := 0
	matchIdx := 0
	for i, arg := range os.Args[1:] {
		if slices.Contains(knownCmds, arg) {
			matches++
			matchIdx = i
		}
	}
	if matches == 1 {
		rest := slices.Delete(slices.Clone(os.Args[1:]), matchIdx, matchIdx+1)
		os.Args = append([]string{os.Args[0], os.Args[1:][matchIdx]}, rest...)
		found = true
	} else if matches >= 2 {
		found = true
	}
	// No command found - count non-flag args to pick a default.
	// Exactly 1 -> list, exactly 2 -> exchange; anything else lets cobra error.
	if !found && len(os.Args) > 1 {
		positional := 0
		for _, arg := range os.Args[1:] {
			if !strings.HasPrefix(arg, "-") {
				positional++
			}
		}
		switch positional {
		case 1:
			os.Args = append([]string{os.Args[0], "list"}, os.Args[1:]...)
		case 2:
			os.Args = append([]string{os.Args[0], "exchange"}, os.Args[1:]...)
		}
	}
	cmd.Execute()
}
