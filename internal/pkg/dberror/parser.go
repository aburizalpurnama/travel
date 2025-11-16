package dberror

import (
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

// uniqueDetailRegex captures the column names and values from a PostgreSQL unique violation detail string.
// Example Detail: "Key (email, username)=(test@example.com, user1) already exists."
// Group 1 captures: "email, username"
// Group 2 captures: "test@example.com, user1"
var uniqueDetailRegex = regexp.MustCompile(`^Key \((.*?)\)=\((.*?)\) already exists\.$`)

// ParseUniqueConstraintError extracts details from a UniqueViolation (23505) error.
// It returns a user-friendly error message and a map of the conflicting fields and their values.
func ParseUniqueConstraintError(pgErr *pgconn.PgError) (string, map[string]any) {
	defaultMessage := "An entry with this data already exists."

	if pgErr.Code != UniqueViolation || pgErr.Detail == "" {
		return defaultMessage, nil
	}

	// Attempt to parse the detail string using regex
	matches := uniqueDetailRegex.FindStringSubmatch(pgErr.Detail)
	if len(matches) != 3 {
		// If parsing fails, return a generic message with the constraint name
		return "Data already exists: " + pgErr.ConstraintName, nil
	}

	keys := strings.Split(matches[1], ", ")
	values := strings.Split(matches[2], ", ")

	if len(keys) != len(values) {
		return defaultMessage, nil
	}

	details := make(map[string]any)
	for i, key := range keys {
		details[key] = convertPsqlValue(values[i])
	}

	// Construct a readable error message listing the conflicting fields
	msg := "An entry with this " + strings.Join(keys, " and ") + " already exists."

	return msg, details
}

// convertPsqlValue is a helper to convert PostgreSQL specific string representations (like 't'/'f')
// into their native Go types.
func convertPsqlValue(val string) any {
	switch val {
	case "t":
		return true
	case "f":
		return false
	}
	// Additional conversions (e.g., numbers) can be added here if needed.
	return val
}
