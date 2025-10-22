package dberror

import (
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

// Pre-compile regex untuk efisiensi.
// Regex ini akan menangkap:
// Grup 1: Daftar kolom (e.g., "name, is_active")
// Grup 2: Daftar nilai (e.g., "dolor, t")
var uniqueDetailRegex = regexp.MustCompile(`^Key \((.*?)\)=\((.*?)\) already exists\.$`)

// convertPsqlValue adalah helper kecil untuk mengubah nilai string 't'/'f' dari Postgres
func convertPsqlValue(val string) any {
	if val == "t" {
		return true
	}
	if val == "f" {
		return false
	}
	// Di sini Anda bisa menambahkan konversi lain jika perlu (misal: string angka ke int)
	return val
}

// ParseUniqueConstraintError mengekstrak detail dari pgErr 23505.
// Ia mengembalikan pesan error yang ramah dan map[string]interface{} dari field yang duplikat.
func ParseUniqueConstraintError(pgErr *pgconn.PgError) (string, map[string]any) {
	// Fallback message default
	defaultMessage := "An entry with this data already exists."

	if pgErr.Code != UniqueViolation || pgErr.Detail == "" {
		return defaultMessage, nil
	}

	matches := uniqueDetailRegex.FindStringSubmatch(pgErr.Detail)
	if len(matches) != 3 {
		// Gagal mem-parsing, kembalikan pesan generik berdasarkan nama constraint
		return "Data already exists: " + pgErr.ConstraintName, nil
	}

	keys := strings.Split(matches[1], ", ")
	values := strings.Split(matches[2], ", ")

	if len(keys) != len(values) {
		// Seharusnya tidak terjadi, tapi untuk jaga-jaga
		return defaultMessage, nil
	}

	details := make(map[string]any)
	for i, key := range keys {
		details[key] = convertPsqlValue(values[i])
	}

	// Buat pesan error yang lebih spesifik
	msg := "An entry with this " + strings.Join(keys, " and ") + " already exists."

	return msg, details
}
