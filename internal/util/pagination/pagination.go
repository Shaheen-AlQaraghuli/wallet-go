package pagination

import (
	"math"
)

const (
	defaultPage    = 1
	defaultPerPage = 10
)

type Pagination struct {
	Total      int `json:"total"`
	Count      int `json:"count"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalPages int `json:"total_pages"`
}

func NewPagination(page int, count int, total int, limit int) *Pagination {
	return &Pagination{
		Total:      total,
		Count:      count,
		Page:       page,
		PerPage:    limit,
		TotalPages: getTotalPages(total, limit),
	}
}

func getTotalPages(total int, limit int) int {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return max(totalPages, defaultPage)
}
