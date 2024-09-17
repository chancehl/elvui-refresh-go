package main

import (
	"fmt"
	"strings"
	"time"
)

// startSpinner runs a spinner in a separate goroutine.
// It continuously displays a rotating spinner and dynamically updates the message using the provided channel.
// The spinner stops when the "done" signal is received.
func StartSpinner(done chan bool, message chan string) {
	const interval = 100

	chars := []rune{'|', '/', '-', '\\'}
	currentMessage := "Loading"
	maxLength := 0

	for {
		select {
		case <-done:
			clearLine(maxLength + 2) // token + space
			return
		case newMessage := <-message:
			currentMessage = newMessage
			if len(newMessage) > maxLength {
				maxLength = len(newMessage)
			}
		default:
			for _, char := range chars {
				clearLinePartial(maxLength+2, char, currentMessage)
				time.Sleep(interval * time.Millisecond)
			}
		}
	}
}

// clearLinePartial clears the line partially and prints the spinner character followed by the message.
func clearLinePartial(lineLength int, char rune, message string) {
	clearLine := strings.Repeat(" ", lineLength)
	fmt.Printf("\r%s\r%c %s", clearLine, char, message)
}

// clearLine clears the entire line by printing spaces equal to the given line length.
func clearLine(lineLength int) {
	clearLine := strings.Repeat(" ", lineLength)
	fmt.Printf("\r%s\r", clearLine)
}
