package transactions

import (
	"context"
	svcModels "wallet/internal/app/models"
	"wallet/internal/util/pagination"
	pkg "wallet/pkg/wallet"
	_  "wallet/internal/util/http/apierror"
	"github.com/gin-gonic/gin"
	jsonlib "wallet/internal/util/http/errors/json"
)

type transactionService interface {
	GetTransactionByID(ctx context.Context, id string) (svcModels.Transaction, error)
	UpdateTransactionStatus(ctx context.Context, id string, status string) (svcModels.Transaction, error)
	ListTransactions(ctx context.Context, query svcModels.QueryTransactions) (svcModels.Transactions, *pagination.Pagination, error)
	CreateTransaction(ctx context.Context, transaction svcModels.CreateTransactionRequest) (svcModels.Transaction, error)
}

type Controller struct {
	transactionSvc transactionService
}

func New(transactionSvc transactionService) *Controller {
 return &Controller{
  transactionSvc: transactionSvc,
 }
}

// GetTransactionByID godoc
//
// @Summary      Get transaction by ID
// @Description  Get transaction by ID
// @ID getTransactionByID
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Transaction ID"
// @Success      200  {object}  pkg.TransactionResponse
// @Failure      400  {object}  apierror.Error
// @Failure      404  {object}  apierror.Error
// @Failure      422  {object}  apierror.Error
// @Failure      500  {object}  apierror.Error
// @Router       /transactions/{id} [get]
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

// UpdateTransactionStatus godoc
//
// @Summary      Update transaction status
// @Description  Update the status of a transaction
// @ID updateTransactionStatus
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id     path      string  true  "Transaction ID"
// @Param        status body      string  true  "New status for the transaction"
// @Success      200    {object}  pkg.TransactionResponse
// @Failure      400    {object}  apierror.Error
// @Failure      404    {object}  apierror.Error
// @Failure      422    {object}  apierror.Error
// @Failure      500    {object}  apierror.Error
// @Router       /transactions/{id}/status [put]
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

// ListTransactions godoc
//
// @Summary      List transactions
// @Description  List transactions with pagination
// @ID listTransactions
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        query  query      pkg.ListTransactionsRequest  true  "Query parameters"
// @Success      200    {object}  pkg.TransactionsResponse
// @Failure      400    {object}  apierror.Error
// @Failure      404    {object}  apierror.Error
// @Failure      422    {object}  apierror.Error
// @Failure      500    {object}  apierror.Error
// @Router       /transactions [get]
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

// CreateTransaction godoc
//
// @Summary      Create transaction
// @Description  Create a new transaction
// @ID createTransaction
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        transaction  body      pkg.CreateTransactionRequest  true  "Transaction data"
// @Success      201    {object}  pkg.TransactionResponse
// @Failure      400    {object}  apierror.Error
// @Failure      404    {object}  apierror.Error
// @Failure      422    {object}  apierror.Error
// @Failure      500    {object}  apierror.Error
// @Router       /transactions [post]
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
