package command

import (
	yup "github.com/gloo-foo/framework"
)

type command yup.Inputs[yup.File, flags]

func Rev(parameters ...any) yup.Command {
	return command(yup.Initialize[yup.File, flags](parameters...))
}

func (p command) Executor() yup.CommandExecutor {
	return yup.LineTransform(func(line string) (string, bool) {
		// Reverse the string
		runes := []rune(line)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes), true
	}).Executor()
}
