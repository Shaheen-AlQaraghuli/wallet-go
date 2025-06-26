package transactions

import (
	"context"

	"github.com/gin-gonic/gin"
	svcModels "wallet/internal/app/models"
	_ "wallet/internal/util/http/apierror"
	jsonlib "wallet/internal/util/http/errors/json"
	"wallet/internal/util/pagination"
	pkg "wallet/pkg/wallet"
)

type transactionService interface {
	GetTransactionByID(ctx context.Context, id string) (svcModels.Transaction, error)
	UpdateTransactionStatus(ctx context.Context, id string, status string) (svcModels.Transaction, error)
	ListTransactions(ctx context.Context, query svcModels.QueryTransactions) (
		svcModels.Transactions, *pagination.Pagination, error)
	CreateTransaction(ctx context.Context, transaction svcModels.CreateTransactionRequest) (
		svcModels.Transaction, error)
}

type Controller struct {
	transactionSvc transactionService
}

func New(transactionSvc transactionService) *Controller {
	return &Controller{
		transactionSvc: transactionSvc,
	}
}

// @Router       /transactions/{id} [get].
func (c *Controller) GetTransactionByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		jsonlib.SendBadRequestError(ctx, "Transaction ID is required")

		return
	}

	transaction, err := c.transactionSvc.GetTransactionByID(ctx, id)
	if err != nil {
		jsonlib.SendGenericAPIError(ctx, err)

		return
	}

	ctx.JSON(200, pkg.TransactionResponse{
		Transaction: transaction.ToResponse(),
	})
}

// @Router       /transactions/{id}/status [put].
func (c *Controller) UpdateTransactionStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		jsonlib.SendBadRequestError(ctx, "Transaction ID is required")

		return
	}

	var req pkg.UpdateTransactionStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		jsonlib.SendApiValidationError(ctx, err)

		return
	}

	transaction, err := c.transactionSvc.UpdateTransactionStatus(ctx, id, string(req.Status))
	if err != nil {
		jsonlib.SendGenericAPIError(ctx, err)

		return
	}

	ctx.JSON(200, pkg.TransactionResponse{
		Transaction: transaction.ToResponse(),
	})
}

// @Router       /transactions [get].
func (c *Controller) ListTransactions(ctx *gin.Context) {
	var req pkg.ListTransactionsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		jsonlib.SendApiValidationError(ctx, err)

		return
	}

	transactions, pagination, err := c.transactionSvc.ListTransactions(ctx, svcModels.QueryTransactions{}.FromRequest(req))
	if err != nil {
		jsonlib.SendGenericAPIError(ctx, err)

		return
	}

	ctx.JSON(200, pkg.TransactionsResponse{
		Transactions: transactions.ToResponse(),
		Metadata: pkg.Metadata{
			Pagination: *pagination,
		},
	})
}

// @Router       /transactions [post].
func (c *Controller) CreateTransaction(ctx *gin.Context) {
	var req pkg.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		jsonlib.SendApiValidationError(ctx, err)

		return
	}

	transaction, err := c.transactionSvc.CreateTransaction(ctx, svcModels.CreateTransactionRequest{}.FromRequest(req))
	if err != nil {
		jsonlib.SendGenericAPIError(ctx, err)

		return
	}

	ctx.JSON(201, pkg.TransactionResponse{
		Transaction: transaction.ToResponse(),
	})
}
