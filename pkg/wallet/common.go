package wallet

import (
	"wallet/internal/util/pagination"
)

type Metadata struct {
	Pagination pagination.Pagination `json:"pagination"`
}
