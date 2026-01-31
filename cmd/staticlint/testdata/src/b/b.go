package utility

import "os"

// This file is in a non-main package, so os.Exit is fine here.
func UtilityExit() {
	os.Exit(0) // No error expected
}
