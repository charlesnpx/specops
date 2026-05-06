package output

import "errors"

const (
	ExitOperational = 1
	ExitUsage       = 2
	ExitContract    = 3
)

type CLIError struct {
	Code    int
	Message string
}

func (e *CLIError) Error() string {
	return e.Message
}

func (e *CLIError) ExitCode() int {
	return e.Code
}

func OperationalError(message string) error {
	return &CLIError{Code: ExitOperational, Message: message}
}

func UsageError(message string) error {
	return &CLIError{Code: ExitUsage, Message: message}
}

func ContractError(message string) error {
	return &CLIError{Code: ExitContract, Message: message}
}

func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	var cliErr interface{ ExitCode() int }
	if errors.As(err, &cliErr) {
		return cliErr.ExitCode()
	}
	return ExitOperational
}
