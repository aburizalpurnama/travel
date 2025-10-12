package paginator

type OffsetBasedOption struct {
	Page *int `json:"page"`
	Size *int `json:"size"`
}

type OffsetBasedPayload struct {
	TotalItems  *int64 `json:"total_items,omitempty"`
	TotalPages  *int   `json:"total_pages,omitempty"`
	CurrentPage *int   `json:"current_page,omitempty"`
	PageSize    *int   `json:"page_size,omitempty"`
}

func GetOffset(page, size int) (offset int) {
	if page > 1 {
		offset = page * size
	}

	return
}

func Parse[T any]() {

}
