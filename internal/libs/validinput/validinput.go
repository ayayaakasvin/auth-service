package validinput

import (
	"regexp"
	"strings"
)

type ValidationError string

func (v ValidationError) Error() string {
	return string(v)
}

const (
	ErrorEmptyTitle               ValidationError = "Event title cannot be empty"
	ErrorEmptyTickets             ValidationError = "Event must have at least one ticket"
	ErrorEmptyLocation            ValidationError = "Event location cannot be empty"
	ErrorInvalidCapacity          ValidationError = "Event capacity is invalid"
	ErrorInvalidLocationPlacement ValidationError = "Event location coordinates are invalid"
	ErrorInvalidTime              ValidationError = "Event time is invalid"
	ErrorInvalidTicketPrice       ValidationError = "Ticket price is invalid"
	ErrorInvalidCurrencyFormat    ValidationError = "Ticket currency format is invalid"
)

var (
	uppercase = regexp.MustCompile(`[A-Z]`)
	lowercase = regexp.MustCompile(`[a-z]`)
	digit     = regexp.MustCompile(`[0-9]`)

	minLengthPassword = 8
	minLengthUsername = 3
)

func IsValidPassword(password string) bool {
	if len(password) < minLengthPassword {
		return false
	}

	if !uppercase.MatchString(password) {
		return false
	}

	if !lowercase.MatchString(password) {
		return false
	}

	if !digit.MatchString(password) {
		return false
	}

	return true
}

func IsValidUsername(username string) bool {
	if len(username) < minLengthUsername {
		return false
	}

	if !(lowercase.MatchString(username) || uppercase.MatchString(username)) {
		return false
	}

	return true
}

func IsValidFileName(filename string) bool {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return false
	}
	// Disallow common illegal filename characters
	illegal := regexp.MustCompile(`[\\/:\*\?"<>\|]`)
	return !illegal.MatchString(filename)
}