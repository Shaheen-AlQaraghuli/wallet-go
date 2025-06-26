package wallets

import (
	"context"
	svcModels "wallet/internal/app/models"
	_ "wallet/internal/util/http/apierror"
	jsonlib "wallet/internal/util/http/errors/json"
	"wallet/internal/util/pagination"
	pkg "wallet/pkg/wallet"

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
// @Success      200  {object}  pkg.WalletResponse
// @Failure      400  {object}  apierror.Error
// @Failure      404  {object}  apierror.Error
// @Failure      422  {object}  apierror.Error
// @Failure      500  {object}  apierror.Error
// @Router       /wallets/{id} [get]
func (c *Controller) GetWalletByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		jsonlib.SendBadRequestError(ctx, "Wallet ID is required")
		return
	}

	wallet, err := c.walletSvc.GetWalletByID(ctx, id)
	if err != nil {
		jsonlib.SendGenericAPIError(ctx, err)
		return
	}

	ctx.JSON(200, pkg.WalletResponse{
		Wallet: wallet.ToResponse(),
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
// @Success      200    {object}  pkg.WalletResponse
// @Failure      400    {object}  apierror.Error
// @Failure      404    {object}  apierror.Error
// @Failure      422    {object}  apierror.Error
// @Failure      500    {object}  apierror.Error
// @Router       /wallets/{id}/status [patch]
func (c *Controller) UpdateWalletStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		jsonlib.SendBadRequestError(ctx, "Wallet ID is required")
		return
	}

	req := pkg.UpdateWalletStatusRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		jsonlib.SendApiValidationError(ctx, err)
		return
	}

	wallet, err := c.walletSvc.UpdateWalletStatus(ctx, id, string(req.Status))
	if err != nil {
		jsonlib.SendGenericAPIError(ctx, err)
		return
	}

	ctx.JSON(200, pkg.WalletResponse{
		Wallet: wallet.ToResponse(),
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
// @Param        query  query     pkg.ListWalletsRequest  false  "Query params"
// @Success      200    {object}  pkg.WalletsResponse
// @Failure      400    {object}  apierror.Error
// @Failure      500    {object}  apierror.Error
// @Router       /wallets [get]
func (c *Controller) ListWallets(ctx *gin.Context) {
	var query pkg.ListWalletsRequest
	if err := ctx.ShouldBindQuery(&query); err != nil {
		jsonlib.SendApiValidationError(ctx, err)
		return
	}

	wallets, pagination, err := c.walletSvc.ListWallets(ctx, svcModels.QueryWallets{}.FromRequest(query))
	if err != nil {
		jsonlib.SendGenericAPIError(ctx, err)
		return
	}

	ctx.JSON(200, pkg.WalletsResponse{
		Wallets: wallets.ToResponse(),
		Metadata: pkg.Metadata{
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
// @Param        wallet body      pkg.CreateWalletRequest  true  "Wallet data"
// @Success      201    {object}  pkg.WalletResponse
// @Failure      400    {object}  apierror.Error
// @Failure      422    {object}  apierror.Error
// @Failure      500    {object}  apierror.Error
// @Router       /wallets [post]
func (c *Controller) CreateWallet(ctx *gin.Context) {
	var req pkg.CreateWalletRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		jsonlib.SendApiValidationError(ctx, err)
		return
	}

	wallet, err := c.walletSvc.CreateWallet(ctx, svcModels.CreateWalletRequest{}.FromRequest(req))
	if err != nil {
		jsonlib.SendGenericAPIError(ctx, err)
		return
	}

	ctx.JSON(200, pkg.WalletResponse{
		Wallet: wallet.ToResponse(),
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
// @Success      200  {object}  pkg.WalletResponse
// @Failure      400  {object}  apierror.Error
// @Failure      404  {object}  apierror.Error
// @Failure      422  {object}	apierror.Error
// @Failure      500  {object}  apierror.Error
// @Router       /wallets/{id}/balance [get]
func (c *Controller) GetWalletWithBalance(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		jsonlib.SendBadRequestError(ctx, "Wallet ID is required")
		return
	}

	wallet, err := c.walletSvc.GetWalletWithBalance(ctx, id)
	if err != nil {
		jsonlib.SendGenericAPIError(ctx, err)
		return
	}

	ctx.JSON(200, pkg.WalletResponse{
		Wallet: wallet.ToResponse(),
	})
}
