package repositories

import (
	"gorm.io/gorm"
	"wallet/internal/util/pagination"
)

const (
	DefaultPageSize = 10
	MaxPerPage      = 10000
)

func Paginate(paginator pagination.Paginator) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := 1
		pageSize := DefaultPageSize

		if paginator.Page != nil {
			page = *paginator.Page
		}

		if paginator.PerPage != nil {
			pageSize = *paginator.PerPage
		}

		if pageSize > MaxPerPage {
			pageSize = MaxPerPage
		}

		offset := (page - 1) * pageSize

		paginator.Page = &page
		paginator.PerPage = &pageSize

		return db.Offset(offset).Limit(pageSize)
	}
}

func GetPaginator(paginator pagination.Paginator) pagination.Paginator {
	page := 1
	perPage := DefaultPageSize

	if paginator.Page == nil || *paginator.Page < 1 {
		paginator.Page = &page
	}

	if paginator.PerPage == nil || *paginator.PerPage < 1 {
		paginator.PerPage = &perPage
	}

	return paginator
}

func CountTotal(db *gorm.DB, paginator pagination.Paginator, count int) (int64, error) {
	if total, ok := paginator.GetTotal(count); ok {
		return int64(total), nil
	}

	var total int64
	err := db.Offset(-1).Limit(-1).Count(&total).Error

	return total, err
}
