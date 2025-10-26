package main

import (
	"Calculator/io"
	"Calculator/logic"
	"Calculator/storage"
)

func main() {
	history := storage.LoadHistory()
	vars := storage.LoadVars()

	io.PrintHistory(history)

	for {
		cmd := io.ReadCommand()
		if cmd == "quit" {
			break
		}

		history = append(history, cmd)

		output, err := logic.Interpret(cmd, vars)
		if err != nil {
			io.PrintError(err)
		} else {
			io.PrintResult(output)
		}

		if len(history) > 10 {
			history = history[len(history)-10:]
		}
	}

	storage.SaveHistory(history)
	storage.SaveVars(vars)
}
