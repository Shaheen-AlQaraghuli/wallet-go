package transactions

import (
	"context"

	"github.com/Shaheen-AlQaraghuli/wallet-go/internal/app/models"
	"github.com/Shaheen-AlQaraghuli/wallet-go/internal/app/repositories"
	"github.com/Shaheen-AlQaraghuli/wallet-go/internal/util/dblib"
	"github.com/Shaheen-AlQaraghuli/wallet-go/internal/util/pagination"
	"gorm.io/gorm"
)

type Repository struct {
	dblib.TxManager
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		TxManager: dblib.NewTxManager(db),
	}
}

func (r *Repository) Create(ctx context.Context, transaction models.Transaction) (models.Transaction, error) {
	if err := r.DB(ctx).Create(&transaction).Error; err != nil {
		return models.Transaction{}, err
	}

	return transaction, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (models.Transaction, error) {
	var transaction models.Transaction
	if err := r.DB(ctx).First(&transaction, "id = ?", id).Error; err != nil {
		return models.Transaction{}, err
	}

	return transaction, nil
}

func (r *Repository) Update(ctx context.Context, transaction models.Transaction) (models.Transaction, error) {
	if err := r.DB(ctx).Save(&transaction).Error; err != nil {
		return models.Transaction{}, err
	}

	return transaction, nil
}

func (r *Repository) List(ctx context.Context, query models.QueryTransactions) (
	[]models.Transaction, *pagination.Pagination, error) {
	var transactions []models.Transaction

	queryBuilder := r.DB(ctx).Model(&models.Transaction{})
	applyFilters(queryBuilder, query)

	paginator := repositories.GetPaginator(query.Paginator)

	if err := queryBuilder.
		Order("created_at DESC").
		Scopes(repositories.Paginate(paginator)).
		Find(&transactions).Error; err != nil {
		return []models.Transaction{}, &pagination.Pagination{}, err
	}

	total, err := repositories.CountTotal(queryBuilder, paginator, len(transactions))
	if err != nil {
		return []models.Transaction{}, &pagination.Pagination{}, err
	}

	return transactions, pagination.NewPagination(
		*paginator.Page,
		len(transactions),
		int(total),
		*paginator.PerPage,
	), nil
}

func (r *Repository) ListAllTransactions(ctx context.Context, walletID string) (models.Transactions, error) {
	var transactions models.Transactions

	if err := r.DB(ctx).
		Where("wallet_id = ?", walletID).
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func applyFilters(db *gorm.DB, query models.QueryTransactions) {
	if len(query.IDs) > 0 {
		db = db.Where("id IN ?", query.IDs)
	}

	if len(query.WalletIDs) > 0 {
		db = db.Where("wallet_id IN ?", query.WalletIDs)
	}

	if len(query.Statuses) > 0 {
		db = db.Where("status IN ?", query.Statuses)
	}

	if len(query.Types) > 0 {
		db = db.Where("type IN ?", query.Types)
	}

	if query.CreatedAtFrom != nil {
		db = db.Where("created_at >= ?", *query.CreatedAtFrom)
	}

	if query.CreatedAtTo != nil {
		db = db.Where("created_at <= ?", *query.CreatedAtTo)
	}
}
