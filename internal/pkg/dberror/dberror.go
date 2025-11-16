package dberror

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// Common SQLSTATE codes handled by the application.
// Reference: https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	// Class 23 — Integrity Constraint Violation
	UniqueViolation     = "23505"
	ForeignKeyViolation = "23503"
	NotNullViolation    = "23502"
	CheckViolation      = "23514"

	// Class 08 — Connection Exception
	ConnectionException = "08000"
)

// GetError extracts *pgconn.PgError from an error chain.
// Returns nil if the error does not contain a PgError.
func GetError(err error) *pgconn.PgError {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr
	}
	return nil
}

// GetSQLState extracts the SQLState code from an error.
// Returns an empty string if the error is not a PgError.
func GetSQLState(err error) string {
	if pgErr := GetError(err); pgErr != nil {
		return pgErr.Code
	}
	return ""
}
