package utils

import (
	"net/mail"
	"strings"
)

const (
	ASCII_CHARS = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	SYMBOLS = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	DIGITS = "1234567890"
)

func IsValidEmail(email string) bool {
	if _, err := mail.ParseAddress(email); err != nil {
		return false
	}
	return true
}

func IsGoodPassword(pass string) (string, bool) {
	password := strings.TrimSpace(pass)

	if password == "" {
		return "Password must be at least 8 characters", false
	}

	if len(password) < 8{
		return "Password must be at least 8 characters", false
	}

	if !strings.ContainsAny(password, ASCII_CHARS){
		return "Password must contain at least one upper case letter", false
	}

	if !strings.ContainsAny(password, SYMBOLS){
		return "Password must contain at least one special symbol", false
	}

	if !strings.ContainsAny(password, DIGITS){
		return "Password must contain at least one number", false
	}

	return pass, true
}

func IsGoodText(text string) (bool) {
	val := strings.TrimSpace(text)

	if val == ""  {
		return false
	}
	return true
}
