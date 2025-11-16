package paginator

// GetOffset calculates the database query offset based on the page number and page size.
// It assumes a 1-based page index. If page is 1 or less, the offset defaults to 0.
func GetOffset(page, size int) (offset int) {
	if page > 1 {
		offset = (page * size) - size
	}

	return
}
