package response

type Simple struct {
	Data interface{} `json:"data"`
}

type Paginate struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

type Meta struct {
	Page       int   `json:"page" example:"1"`
	PerPage    int   `json:"per_page" example:"10"`
	TotalItems int64 `json:"total_items" example:"125"`
	TotalPages int   `json:"total_pages" example:"13"`
	HasNext    bool  `json:"has_next" example:"true"`
	HasPrev    bool  `json:"has_prev" example:"false"`
}
