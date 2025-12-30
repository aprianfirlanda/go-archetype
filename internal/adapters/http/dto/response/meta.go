package response

type PaginationMeta struct {
	Page       int   `json:"page" example:"1"`
	Limit      int   `json:"limit" example:"10"`
	From       int64 `json:"from" example:"11"`
	To         int64 `json:"to" example:"20"`
	TotalItems int64 `json:"total_items" example:"125"`
	TotalPages int   `json:"total_pages" example:"13"`
	HasNext    bool  `json:"has_next" example:"true"`
	HasPrev    bool  `json:"has_prev" example:"false"`
}

func NewPaginationMeta(
	page int,
	limit int,
	totalItems int64,
) PaginationMeta {
	totalPages := int((totalItems + int64(limit) - 1) / int64(limit))

	from := int64((page-1)*limit + 1)
	to := from + int64(limit) - 1
	if totalItems == 0 {
		from = 0
		to = 0
	}
	if to > totalItems {
		to = totalItems
	}

	return PaginationMeta{
		Page:       page,
		Limit:      limit,
		From:       from,
		To:         to,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}
