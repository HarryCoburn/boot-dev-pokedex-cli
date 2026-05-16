package main

import (
	"bufio"
	"os"
)

func main() {
	runREPL(bufio.NewScanner(os.Stdin), commandMap)
}
