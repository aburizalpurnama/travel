package dberror

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// SQLState adalah tipe kustom untuk SQLSTATE error code.
type SQLState string

// Daftar SQLSTATE yang umum ditangani di aplikasi.
// Referensi: https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	// Class 23 — Integrity Constraint Violation
	UniqueViolation     SQLState = "23505"
	ForeignKeyViolation SQLState = "23503"
	NotNullViolation    SQLState = "23502"
	CheckViolation      SQLState = "23514"

	// Class 08 — Connection Exception
	ConnectionException SQLState = "08000"
)

// GetError mel-ekstrak *pgconn.PgError dari sebuah error.
// Mengembalikan nil jika error tersebut bukan PgError.
func GetError(err error) *pgconn.PgError {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr
	}
	return nil
}

// GetSQLState mel-ekstrak SQLState code dari sebuah error.
// Mengembalikan string kosong jika error tersebut bukan PgError.
func GetSQLState(err error) SQLState {
	if pgErr := GetError(err); pgErr != nil {
		return SQLState(pgErr.Code)
	}
	return ""
}
