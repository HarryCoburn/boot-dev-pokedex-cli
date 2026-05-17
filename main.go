package main

import (
	"bufio"
	"os"
)

func main() {
	setupREPL(bufio.NewScanner(os.Stdin))
}
