package wallets

import (
	"context"

	svcModels "github.com/Shaheen-AlQaraghuli/wallet-go/internal/app/models"
	_ "github.com/Shaheen-AlQaraghuli/wallet-go/internal/util/http/apierror"
	jsonlib "github.com/Shaheen-AlQaraghuli/wallet-go/internal/util/http/errors/json"
	"github.com/Shaheen-AlQaraghuli/wallet-go/internal/util/pagination"
	"github.com/Shaheen-AlQaraghuli/wallet-go/pkg/wallet"
	"github.com/gin-gonic/gin"
)

type walletService interface {
	GetWalletByID(ctx context.Context, id string) (svcModels.Wallet, error)
	UpdateWalletStatus(ctx context.Context, id, status string) (svcModels.Wallet, error)
	ListWallets(ctx context.Context, query svcModels.QueryWallets) (svcModels.Wallets, *pagination.Pagination, error)
	CreateWallet(ctx context.Context, wallet svcModels.CreateWalletRequest) (svcModels.Wallet, error)
	GetWalletWithBalance(ctx context.Context, id string) (svcModels.Wallet, error)
}

type Controller struct {
	walletSvc walletService
}

func New(walletSvc walletService) *Controller {
	return &Controller{
		walletSvc: walletSvc,
	}
}

// GetWalletByID godoc
//
// @Summary      Get wallet by ID
// @Description  Get wallet by ID
// @ID getWalletByID
// @Tags         wallets
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Wallet ID"
// @Success      200  {object}  wallet.WalletResponse
// @Failure      400  {object}  apierror.Error
// @Failure      404  {object}  apierror.Error
// @Failure      422  {object}  apierror.Error
// @Failure      500  {object}  apierror.Error
// @Router       /v1/wallets/{id} [get]
func (c *Controller) GetWalletByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		jsonlib.SendBadRequestError(ctx, "Wallet ID is required")

		return
	}

	walletResp, err := c.walletSvc.GetWalletByID(ctx, id)
	if err != nil {
		jsonlib.SendGenericAPIError(ctx, err)

		return
	}

	ctx.JSON(200, wallet.WalletResponse{
		Wallet: walletResp.ToResponse(),
	})
}

// UpdateWalletStatus godoc
//
// @Summary      Update wallet status
// @Description  Update wallet status
// @ID updateWalletStatus
// @Tags         wallets
// @Accept       json
// @Produce      json
// @Param        id     path      string  true  "Wallet ID"
// @Param        status body      string  true  "New wallet status"
// @Success      200    {object}  wallet.WalletResponse
// @Failure      400    {object}  apierror.Error
// @Failure      404    {object}  apierror.Error
// @Failure      422    {object}  apierror.Error
// @Failure      500    {object}  apierror.Error
// @Router       /v1/wallets/{id}/status [patch]
func (c *Controller) UpdateWalletStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		jsonlib.SendBadRequestError(ctx, "Wallet ID is required")

		return
	}

	req := wallet.UpdateWalletStatusRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		jsonlib.SendApiValidationError(ctx, err)

		return
	}

	walletResp, err := c.walletSvc.UpdateWalletStatus(ctx, id, string(req.Status))
	if err != nil {
		jsonlib.SendGenericAPIError(ctx, err)

		return
	}

	ctx.JSON(200, wallet.WalletResponse{
		Wallet: walletResp.ToResponse(),
	})
}

// ListWallets godoc
//
// @Summary      List wallets
// @Description  List wallets with pagination
// @ID listWallets
// @Tags         wallets
// @Accept       json
// @Produce      json
// @Param        query  query     wallet.ListWalletsRequest  false  "Query params"
// @Success      200    {object}  wallet.WalletsResponse
// @Failure      400    {object}  apierror.Error
// @Failure      500    {object}  apierror.Error
// @Router       /v1/wallets [get]
func (c *Controller) ListWallets(ctx *gin.Context) {
	var query wallet.ListWalletsRequest
	if err := ctx.ShouldBindQuery(&query); err != nil {
		jsonlib.SendApiValidationError(ctx, err)

		return
	}

	wallets, pagination, err := c.walletSvc.ListWallets(ctx, svcModels.QueryWallets{}.FromRequest(query))
	if err != nil {
		jsonlib.SendGenericAPIError(ctx, err)

		return
	}

	ctx.JSON(200, wallet.WalletsResponse{
		Wallets: wallets.ToResponse(),
		Metadata: wallet.Metadata{
			Pagination: *pagination,
		},
	})
}

// CreateWallet godoc
//
// @Summary      Create a new wallet
// @Description  Create a new wallet with initial balance
// @ID createWallet
// @Tags		 wallets
// @Accept       json
// @Produce      json
// @Param        wallet body      wallet.CreateWalletRequest  true  "Wallet data"
// @Success      201    {object}  wallet.WalletResponse
// @Failure      400    {object}  apierror.Error
// @Failure      422    {object}  apierror.Error
// @Failure      500    {object}  apierror.Error
// @Router       /v1/wallets [post]
func (c *Controller) CreateWallet(ctx *gin.Context) {
	var req wallet.CreateWalletRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		jsonlib.SendApiValidationError(ctx, err)

		return
	}

	walletResp, err := c.walletSvc.CreateWallet(ctx, svcModels.CreateWalletRequest{}.FromRequest(req))
	if err != nil {
		jsonlib.SendGenericAPIError(ctx, err)

		return
	}

	ctx.JSON(200, wallet.WalletResponse{
		Wallet: walletResp.ToResponse(),
	})
}

// GetWalletWithBalance godoc
//
// @Summary      Get wallet with balance
// @Description  Get wallet with balance by ID
// @ID getWalletWithBalance
// @Tags         wallets
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Wallet ID"
// @Success      200  {object}  wallet.WalletResponse
// @Failure      400  {object}  apierror.Error
// @Failure      404  {object}  apierror.Error
// @Failure      422  {object}	apierror.Error
// @Failure      500  {object}  apierror.Error
// @Router       /v1/wallets/{id}/balance [get]
func (c *Controller) GetWalletWithBalance(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		jsonlib.SendBadRequestError(ctx, "Wallet ID is required")

		return
	}

	walletResp, err := c.walletSvc.GetWalletWithBalance(ctx, id)
	if err != nil {
		jsonlib.SendGenericAPIError(ctx, err)

		return
	}

	ctx.JSON(200, wallet.WalletResponse{
		Wallet: walletResp.ToResponse(),
	})
}
