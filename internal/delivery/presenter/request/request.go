package request

type Pagination struct {
	Keyword   string `json:"keyword"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
	OrderBy   string `json:"order_by"`
	OrderType string `json:"order_type"`
}

func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *Pagination) GetLimit() int {
	return p.PageSize
}
