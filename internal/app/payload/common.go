package payload

// CommonGetAllRequest defines standard pagination and sorting parameters for list retrieval endpoints.
type CommonGetAllRequest struct {
	Page      *int    `query:"page"`
	Size      *int    `query:"size"`
	OrderBy   *string `query:"order_by"`
	OrderType *string `query:"order_type"`
}

// SetDefault applies default values for pagination and sorting if they are not provided.
// Defaults: Page=1, Size=10, OrderBy="id", OrderType="desc".
func (req *CommonGetAllRequest) SetDefault() {
	if req.Page == nil {
		page := 1
		req.Page = &page
	}

	if req.Size == nil {
		size := 10
		req.Size = &size
	}

	if req.OrderBy == nil {
		orderBy := "id"
		req.OrderBy = &orderBy
	}

	if req.OrderType == nil {
		orderType := "desc"
		req.OrderType = &orderType
	}
}
