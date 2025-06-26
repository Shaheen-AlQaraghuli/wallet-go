package wallets

import (
	"context"

	"gorm.io/gorm"
	"wallet/internal/app/models"
	"wallet/internal/app/repositories"
	"wallet/internal/util/dblib"
	"wallet/internal/util/pagination"
)

type Repository struct {
	dblib.TxManager
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		TxManager: dblib.NewTxManager(db),
	}
}

func (r *Repository) Create(ctx context.Context, wallet models.Wallet) (models.Wallet, error) {
	if err := r.DB(ctx).Create(&wallet).Error; err != nil {
		return models.Wallet{}, err
	}

	return wallet, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (models.Wallet, error) {
	var wallet models.Wallet
	if err := r.DB(ctx).First(&wallet, "id = ?", id).Error; err != nil {
		return models.Wallet{}, err
	}

	return wallet, nil
}

func (r *Repository) Update(ctx context.Context, wallet models.Wallet) (models.Wallet, error) {
	if err := r.DB(ctx).Save(&wallet).Error; err != nil {
		return models.Wallet{}, err
	}

	return wallet, nil
}

func (r *Repository) List(ctx context.Context, query models.QueryWallets) (
	[]models.Wallet, *pagination.Pagination, error) {
	var wallets []models.Wallet

	queryBuilder := r.DB(ctx).Model(&models.Wallet{})
	applyFilters(queryBuilder, query)

	paginator := repositories.GetPaginator(query.Paginator)

	if err := queryBuilder.Scopes(repositories.Paginate(paginator)).Find(&wallets).Error; err != nil {
		return nil, nil, err
	}

	total, err := repositories.CountTotal(queryBuilder, paginator, len(wallets))
	if err != nil {
		return nil, nil, err
	}

	return wallets, pagination.NewPagination(
		*paginator.Page,
		len(wallets),
		int(total),
		*paginator.PerPage,
	), nil
}

func applyFilters(db *gorm.DB, query models.QueryWallets) {
	if len(query.IDs) > 0 {
		db.Where("id IN ?", query.IDs)
	}

	if len(query.OwnerIDs) > 0 {
		db.Where("owner_id IN ?", query.OwnerIDs)
	}

	if len(query.Currencies) > 0 {
		db.Where("currency IN ?", query.Currencies)
	}
}
