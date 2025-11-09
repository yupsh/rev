package command

import (
	gloo "github.com/gloo-foo/framework"
)

type command gloo.Inputs[gloo.File, flags]

func Rev(parameters ...any) gloo.Command {
	return command(gloo.Initialize[gloo.File, flags](parameters...))
}

func (p command) Executor() gloo.CommandExecutor {
	return gloo.LineTransform(func(line string) (string, bool) {
		// Reverse the string
		runes := []rune(line)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes), true
	}).Executor()
}
