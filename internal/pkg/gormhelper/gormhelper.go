package gormhelper

import (
	"errors"
	"fmt"
	"reflect"
	stdStrings "strings"

	"github.com/aburizalpurnama/travel/internal/pkg/strings"
	"gorm.io/gorm"
)

// ParseFilter dynamically adds WHERE clauses to a GORM query based on a filter struct.
// It inspects the struct fields and their tags to construct the query.
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

		// Only process pointer fields that are not nil
		if value.Kind() == reflect.Pointer && !value.IsNil() {

			var columnName string
			queryTag := field.Tag.Get("query")

			// Ignore fields with tag `query:"-"`.
			// This is useful for filters where the attribute name does not match the database column name directly,
			// e.g., filtering data within a date range using start_date and end_date.
			// The user must handle the query customization manually in the repository for such cases.
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

			// Special handling for 'search' field
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

					// Create OR clauses for each searchable field
					for _, searchableField := range searchableFields {
						orClauses = append(orClauses, fmt.Sprintf("%s ILIKE ?", searchableField))
						orArgs = append(orArgs, "%"+searchTerm+"%")
					}

					// Combine all clauses with " OR "
					queryString := stdStrings.Join(orClauses, " OR ")
					db = db.Where(queryString, orArgs...)
				}

			} else {
				// For standard filters, use equality query (=)
				db = db.Where(fmt.Sprintf("%s = ?", columnName), actualValue)
			}
		}
	}

	return db, nil
}
