package validation

import (
	"net/mail"
	"strings"
)

// IsValidEmail performs a practical email validation.
// It is stricter than just checking for an "@".
func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}

	// Use the standard library parser first
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	// mail.ParseAddress accepts "Name <email@domain.com>",
	// so we ensure it is a pure email
	if addr.Address != email {
		return false
	}

	// Extra practical checks
	parts := strings.Split(addr.Address, "@")
	if len(parts) != 2 {
		return false
	}

	local, domain := parts[0], parts[1]
	if local == "" || domain == "" {
		return false
	}

	if strings.Contains(domain, ".") == false {
		return false
	}

	return true
}
