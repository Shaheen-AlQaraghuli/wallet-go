package main

import (
	"time"
	
	transactionSvc "wallet/internal/app/services/transactions"
	walletSvc "wallet/internal/app/services/wallets"

	transactionCtrl "wallet/internal/app/controller/transactions"
	walletCtrl "wallet/internal/app/controller/wallets"

	transactionsRepo "wallet/internal/app/repositories/transactions"
	walletRepo "wallet/internal/app/repositories/wallets"

	cacher "wallet/internal/app/cache"

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
