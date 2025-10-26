package io

import (
	"bufio"
	"fmt"
	"os"
)

var scanner = bufio.NewScanner(os.Stdin)

func PrintHistory(history []string) {
	for _, h := range history {
		fmt.Println(h)
	}
}

func ReadCommand() string {
	fmt.Print("-> ")
	if !scanner.Scan() {
		return "quit"
	}
	return scanner.Text()
}

func PrintError(err error) {
	fmt.Println("Error:", err)
}

func PrintResult(res string) {
	fmt.Println(res)
}
