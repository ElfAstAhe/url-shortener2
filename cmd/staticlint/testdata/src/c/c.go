package main

import "os"

func main() {
	// This is fine, the call is elsewhere
	callExit()
}

func callExit() {
	// os.Exit is fine here because it's not strictly *inside* func main() in the AST traversal
	os.Exit(1) // No error expected based on your current AST inspection logic
}
