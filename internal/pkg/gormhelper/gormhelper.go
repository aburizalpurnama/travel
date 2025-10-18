package gormhelper

import (
	"errors"
	"fmt"
	"reflect"

	stdStrings "strings"

	"github.com/aburizalpurnama/travel/internal/pkg/strings"
	"gorm.io/gorm"
)

// ParseFilter secara dinamis menambahkan klausa WHERE ke query GORM dari sebuah struct filter.
func ParseFilter(db *gorm.DB, filter any) (*gorm.DB, error) {
	if db == nil {
		return nil, errors.New("gormhelper error: db cannot be nil")
	}
	if filter == nil {
		return db, nil
	}

	val := reflect.ValueOf(filter)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	if !val.IsValid() || val.Kind() != reflect.Struct {
		return nil, errors.New("gormhelper error: filter type is not a struct or is nil")
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		// Hanya proses field yang bertipe pointer dan nilainya tidak nil
		if value.Kind() == reflect.Pointer && !value.IsNil() {

			var columnName string
			queryTag := field.Tag.Get("query")

			// Ignore field dengan tag `query:"-"`.
			// Ignorance ini dapat dimanfaatkan untuk filter yang nama attributnya tidak sama dengan nama kolom pada database, misal filter untuk mendapatkan data dengan rentang created_on tertentu menggunakan start_date dan end_date. User perlu melakukan penyesuaian query secara manual di bagian repository untuk case tsb.
			if queryTag == "-" {
				continue
			}

			if queryTag != "" {
				columnName = stdStrings.Split(queryTag, ";")[0]
			}

			if columnName == "" {
				columnName = strings.ToSnakeCase(field.Name)
			}

			actualValue := value.Elem().Interface()

			// Kondisi khusus untuk 'search'
			if columnName == "search" {
				var searchableFields []string
				searchTag := field.Tag.Get("search")
				if searchTag != "" {
					searchableFields = stdStrings.Split(searchTag, ",")
				}

				searchTerm, ok := actualValue.(string)
				if ok && searchTerm != "" && len(searchableFields) > 0 {
					var orClauses []string
					var orArgs []any

					// Buat klausa OR untuk setiap field yang bisa dicari
					for _, searchableField := range searchableFields {
						orClauses = append(orClauses, fmt.Sprintf("%s ILIKE ?", searchableField))
						orArgs = append(orArgs, "%"+searchTerm+"%")
					}

					// Gabungkan semua klausa dengan " OR "
					queryString := stdStrings.Join(orClauses, " OR ")
					db = db.Where(queryString, orArgs...)
				}

			} else {
				// Untuk filter biasa, gunakan query equality (=)
				db = db.Where(fmt.Sprintf("%s = ?", columnName), actualValue)
			}
		}
	}

	return db, nil
}
