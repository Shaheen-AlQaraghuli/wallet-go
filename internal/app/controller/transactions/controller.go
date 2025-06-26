package transactions

import (
	"context"

	svcModels "github.com/Shaheen-AlQaraghuli/wallet-go/internal/app/models"
	_ "github.com/Shaheen-AlQaraghuli/wallet-go/internal/util/http/apierror"
	jsonlib "github.com/Shaheen-AlQaraghuli/wallet-go/internal/util/http/errors/json"
	"github.com/Shaheen-AlQaraghuli/wallet-go/internal/util/pagination"
	"github.com/Shaheen-AlQaraghuli/wallet-go/pkg/wallet"
	"github.com/gin-gonic/gin"
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

// GetTransactionByID godoc
//
// @Summary      Get transaction by ID
// @Description  Get transaction by ID
// @ID getTransactionByID
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Transaction ID"
// @Success      200  {object}  wallet.TransactionResponse
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

	ctx.JSON(200, wallet.TransactionResponse{
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
// @Success      200    {object}  wallet.TransactionResponse
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

	var req wallet.UpdateTransactionStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		jsonlib.SendApiValidationError(ctx, err)

		return
	}

	transaction, err := c.transactionSvc.UpdateTransactionStatus(ctx, id, string(req.Status))
	if err != nil {
		jsonlib.SendGenericAPIError(ctx, err)

		return
	}

	ctx.JSON(200, wallet.TransactionResponse{
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
// @Param        query  query      wallet.ListTransactionsRequest  true  "Query parameters"
// @Success      200    {object}  wallet.TransactionsResponse
// @Failure      400    {object}  apierror.Error
// @Failure      404    {object}  apierror.Error
// @Failure      422    {object}  apierror.Error
// @Failure      500    {object}  apierror.Error
// @Router       /transactions [get]
func (c *Controller) ListTransactions(ctx *gin.Context) {
	var req wallet.ListTransactionsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		jsonlib.SendApiValidationError(ctx, err)

		return
	}

	transactions, pagination, err := c.transactionSvc.ListTransactions(ctx, svcModels.QueryTransactions{}.FromRequest(req))
	if err != nil {
		jsonlib.SendGenericAPIError(ctx, err)

		return
	}

	ctx.JSON(200, wallet.TransactionsResponse{
		Transactions: transactions.ToResponse(),
		Metadata: wallet.Metadata{
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
// @Param        transaction  body      wallet.CreateTransactionRequest  true  "Transaction data"
// @Success      201    {object}  wallet.TransactionResponse
// @Failure      400    {object}  apierror.Error
// @Failure      404    {object}  apierror.Error
// @Failure      422    {object}  apierror.Error
// @Failure      500    {object}  apierror.Error
// @Router       /transactions [post]
func (c *Controller) CreateTransaction(ctx *gin.Context) {
	var req wallet.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		jsonlib.SendApiValidationError(ctx, err)

		return
	}

	transaction, err := c.transactionSvc.CreateTransaction(ctx, svcModels.CreateTransactionRequest{}.FromRequest(req))
	if err != nil {
		jsonlib.SendGenericAPIError(ctx, err)

		return
	}

	ctx.JSON(201, wallet.TransactionResponse{
		Transaction: transaction.ToResponse(),
	})
}
