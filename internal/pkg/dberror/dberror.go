package dberror

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// Daftar SQLSTATE yang umum ditangani di aplikasi.
// Referensi: https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	// Class 23 — Integrity Constraint Violation
	UniqueViolation     string = "23505"
	ForeignKeyViolation string = "23503"
	NotNullViolation    string = "23502"
	CheckViolation      string = "23514"

	// Class 08 — Connection Exception
	ConnectionException string = "08000"
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
func GetSQLState(err error) string {
	if pgErr := GetError(err); pgErr != nil {
		return pgErr.Code
	}
	return ""
}
