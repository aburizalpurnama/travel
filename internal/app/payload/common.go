package payload

type CommonGetAllRequest struct {
	Page      *int    `query:"page"`
	Size      *int    `query:"size"`
	OrderBy   *string `query:"order_by"`
	OrderType *string `query:"order_type"`
}

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
