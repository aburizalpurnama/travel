package paginator

func GetOffset(page, size int) (offset int) {
	if page > 1 {
		offset = (page * size) - size
	}

	return
}
