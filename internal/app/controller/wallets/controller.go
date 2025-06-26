package wallets

import (
	"context"

	svcModels "github.com/Shaheen-AlQaraghuli/wallet-go/internal/app/models"
	_ "github.com/Shaheen-AlQaraghuli/wallet-go/internal/util/http/apierror"
	jsonlib "github.com/Shaheen-AlQaraghuli/wallet-go/internal/util/http/errors/json"
	"github.com/Shaheen-AlQaraghuli/wallet-go/internal/util/pagination"
	pkg "github.com/Shaheen-AlQaraghuli/wallet-go/pkg/wallet"
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

// @Router       /wallets/{id} [get].
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

// @Router       /wallets/{id}/status [patch].
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

// @Router       /wallets [get].
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

// @Router       /wallets [post].
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

// @Router       /wallets/{id}/balance [get].
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
