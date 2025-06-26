package main

import (
	"time"

	cacher "github.com/Shaheen-AlQaraghuli/wallet-go/internal/app/cache"
	transactionCtrl "github.com/Shaheen-AlQaraghuli/wallet-go/internal/app/controller/transactions"
	walletCtrl "github.com/Shaheen-AlQaraghuli/wallet-go/internal/app/controller/wallets"
	transactionsRepo "github.com/Shaheen-AlQaraghuli/wallet-go/internal/app/repositories/transactions"
	walletRepo "github.com/Shaheen-AlQaraghuli/wallet-go/internal/app/repositories/wallets"
	transactionSvc "github.com/Shaheen-AlQaraghuli/wallet-go/internal/app/services/transactions"
	walletSvc "github.com/Shaheen-AlQaraghuli/wallet-go/internal/app/services/wallets"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func addWalletRoutes(db *gorm.DB, cache *cacher.Cache, routerGroup *gin.RouterGroup) {
	repo := walletRepo.New(db)
	transactionsRepo := transactionsRepo.New(db)
	transactionService := transactionSvc.NewService(repo, transactionsRepo, cache, time.Now)
	walletService := walletSvc.NewService(transactionService, repo, cache, time.Now)
	walletController := walletCtrl.New(walletService)

	routerGroup.GET("/wallets", walletController.ListWallets)
	routerGroup.POST("/wallets", walletController.CreateWallet)
	routerGroup.GET("/wallets/:id", walletController.GetWalletByID)
	routerGroup.PATCH("/wallets/:id/status", walletController.UpdateWalletStatus)
	routerGroup.GET("/wallets/:id/balance", walletController.GetWalletWithBalance)
}

func addTransactionRoutes(db *gorm.DB, cache *cacher.Cache, routerGroup *gin.RouterGroup) {
	repo := transactionsRepo.New(db)
	walletRepo := walletRepo.New(db)
	transactionService := transactionSvc.NewService(walletRepo, repo, cache, time.Now)
	transactionController := transactionCtrl.New(transactionService)

	routerGroup.GET("/transactions", transactionController.ListTransactions)
	routerGroup.POST("/transactions", transactionController.CreateTransaction)
	routerGroup.GET("/transactions/:id", transactionController.GetTransactionByID)
	routerGroup.PATCH("/transactions/:id/status", transactionController.UpdateTransactionStatus)
}
