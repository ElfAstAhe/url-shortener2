package main

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestOsExitChecker(t *testing.T) {
	// The analysistest package handles running the analyzer against
	// the source files found in the 'testdata' directory.
	// We run it against three different scenarios:
	// 1. "a": Code that should trigger the diagnostic (os.Exit in main).
	// 2. "b": Code that should NOT trigger the diagnostic (os.Exit in a non-main package).
	// 3. "c": Code that should NOT trigger the diagnostic (os.Exit in another function, not main).
	analysistest.Run(t, analysistest.TestData(), OsExitAnalyzer, "b", "c")
}
