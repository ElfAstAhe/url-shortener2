package main

import "os"

func main() {
	// Use standard double quotes for literal string matching in the test pattern:
	os.Exit(1) // want "Using call os.Exit() into main function."
}

func otherFunc() {
	os.Exit(0) // This is now ignored by the improved analyzer logic
}
