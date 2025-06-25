package pagination

type Paginator struct {
	Page    *int `binding:"omitempty,gte=1"         form:"page,omitempty"     json:"page,omitempty"     url:"page,omitempty"`
	PerPage *int `binding:"omitempty,gte=1,lte=100" form:"per_page,omitempty" json:"per_page,omitempty" url:"per_page,omitempty"`
}

func (p *Paginator) GetTotal(count int) (int, bool) {
	shouldQueryForTotal := isNullOrZero(p.Page) || isNullOrZero(p.PerPage) || (count < 1 && *p.Page > 1 && *p.PerPage > 0) || count >= *p.PerPage

	if shouldQueryForTotal {
		return 0, false
	}

	page := *p.Page
	pageSize := *p.PerPage
	offset := (page - 1) * pageSize

	return offset + count, true
}

func isNullOrZero(i *int) bool {
	if i == nil {
		return true
	}

	return *i <= 0
}
