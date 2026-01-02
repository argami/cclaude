package utils

import (
	"fmt"
	"os"
)

type ExitCode int

const (
	ExitSuccess ExitCode = iota
	ExitProviderNotFound
	ExitAPIKeyMissing
	ExitClaudeNotFound
	ExitConfigError
	ExitValidationError
)

func HandleError(err error, code ExitCode) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(int(code))
	}
}

func HandleErrorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "❌ Error: %s\n", fmt.Sprintf(format, args...))
	os.Exit(1)
}