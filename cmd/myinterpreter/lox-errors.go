package main

type loxError struct {
	message string
	line    int
	code    int
}

func (e loxError) Error() string {
	return e.message
}

func newParsingError(message string) error {
	return &loxError{
		message: message,
		code:    65,
	}
}

func newRuntimeError(message string) error {
	return loxError{
		message: message,
		code:    70,
	}
}
