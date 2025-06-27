package transactions

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Shaheen-AlQaraghuli/wallet-go/internal/app/models"
	"github.com/Shaheen-AlQaraghuli/wallet-go/internal/app/services/transactions/mocks"
	"github.com/Shaheen-AlQaraghuli/wallet-go/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRunningBalance(t *testing.T) {
    tests := []struct {
        name            string
        walletID        string
        mockSetup       func(*mocks.MockTransactionRepo, *mocks.MockCacheClient)
        expectedBalance int
        expectedError   string
    }{
        {
            name:     "balance from cache - cache hit",
            walletID: "wallet-123",
            mockSetup: func(tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                balance := 1200
                c.On("GetBalance", mock.Anything, "wallet-123").Return(&balance, nil)
            },
            expectedBalance: 1200,
        },
        {
            name:     "credit transactions only counted when completed",
            walletID: "wallet-123",
            mockSetup: func(tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Cache miss
                c.On("GetBalance", mock.Anything, "wallet-123").Return((*int)(nil), nil)

                // Mock wallet lock
                unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
                c.On("Mutex", mock.Anything, "wallet-123").Return(unlockFunc, nil)

                // Transactions: only completed credit should count
                // Expected balance: 1000 (completed credit) + 0 (pending credit not counted) = 1000
                tr.On("ListAllTransactions", mock.Anything, "wallet-123").Return(models.Transactions{
                    {Amount: 1000, Type: string(types.TransactionTypeCredit), Status: string(types.TransactionStatusCompleted)},
                    {Amount: 500, Type: string(types.TransactionTypeCredit), Status: string(types.TransactionStatusPending)}, // Not counted
                    {Amount: 300, Type: string(types.TransactionTypeCredit), Status: string(types.TransactionStatusFailed)},   // Not counted
                }, nil)

                // Set cache with calculated balance
                c.On("SetBalance", mock.Anything, "wallet-123", 1000).Return(nil)
            },
            expectedBalance: 1000,
        },
        {
            name:     "debit transactions counted when pending or completed, not failed",
            walletID: "wallet-123",
            mockSetup: func(tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Cache miss
                c.On("GetBalance", mock.Anything, "wallet-123").Return((*int)(nil), nil)

                // Mock wallet lock
                unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
                c.On("Mutex", mock.Anything, "wallet-123").Return(unlockFunc, nil)

                // Transactions: debit pending and completed should count, failed should not
                // Expected balance: 2000 - 300 (completed debit) - 200 (pending debit) + 0 (failed debit not counted) = 1500
                tr.On("ListAllTransactions", mock.Anything, "wallet-123").Return(models.Transactions{
                    {Amount: 2000, Type: string(types.TransactionTypeCredit), Status: string(types.TransactionStatusCompleted)},
                    {Amount: 300, Type: string(types.TransactionTypeDebit), Status: string(types.TransactionStatusCompleted)},
                    {Amount: 200, Type: string(types.TransactionTypeDebit), Status: string(types.TransactionStatusPending)},
                    {Amount: 100, Type: string(types.TransactionTypeDebit), Status: string(types.TransactionStatusFailed)}, // Not counted
                }, nil)

                // Set cache with calculated balance
                c.On("SetBalance", mock.Anything, "wallet-123", 1500).Return(nil)
            },
            expectedBalance: 1500,
        },
        {
            name:     "mixed transaction types with various statuses",
            walletID: "wallet-123",
            mockSetup: func(tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Cache miss
                c.On("GetBalance", mock.Anything, "wallet-123").Return((*int)(nil), nil)

                // Mock wallet lock
                unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
                c.On("Mutex", mock.Anything, "wallet-123").Return(unlockFunc, nil)

                // Credits (only completed): 1000 + 800 = 1800
                // Debits (pending + completed): 300 + 150 = 450
                // Failed transactions: ignored
                // Expected balance: 1800 - 450 = 1350
                tr.On("ListAllTransactions", mock.Anything, "wallet-123").Return(models.Transactions{
                    // Completed credits (counted)
                    {Amount: 1000, Type: string(types.TransactionTypeCredit), Status: string(types.TransactionStatusCompleted)},
                    {Amount: 800, Type: string(types.TransactionTypeCredit), Status: string(types.TransactionStatusCompleted)},
                    // Pending credit (not counted)
                    {Amount: 500, Type: string(types.TransactionTypeCredit), Status: string(types.TransactionStatusPending)},
                    // Failed credit (not counted)
                    {Amount: 200, Type: string(types.TransactionTypeCredit), Status: string(types.TransactionStatusFailed)},
                    // Completed debit (counted)
                    {Amount: 300, Type: string(types.TransactionTypeDebit), Status: string(types.TransactionStatusCompleted)},
                    // Pending debit (counted)
                    {Amount: 150, Type: string(types.TransactionTypeDebit), Status: string(types.TransactionStatusPending)},
                    // Failed debit (not counted)
                    {Amount: 250, Type: string(types.TransactionTypeDebit), Status: string(types.TransactionStatusFailed)},
                }, nil)

                // Set cache with calculated balance
                c.On("SetBalance", mock.Anything, "wallet-123", 1350).Return(nil)
            },
            expectedBalance: 1350,
        },
        {
            name:     "zero balance with failed transactions",
            walletID: "wallet-123",
            mockSetup: func(tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Cache miss
                c.On("GetBalance", mock.Anything, "wallet-123").Return((*int)(nil), nil)

                // Mock wallet lock
                unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
                c.On("Mutex", mock.Anything, "wallet-123").Return(unlockFunc, nil)

                // All transactions are failed or pending credits - should result in 0 balance
                tr.On("ListAllTransactions", mock.Anything, "wallet-123").Return(models.Transactions{
                    {Amount: 1000, Type: string(types.TransactionTypeCredit), Status: string(types.TransactionStatusPending)}, // Not counted
                    {Amount: 500, Type: string(types.TransactionTypeCredit), Status: string(types.TransactionStatusFailed)},   // Not counted
                    {Amount: 300, Type: string(types.TransactionTypeDebit), Status: string(types.TransactionStatusFailed)},    // Not counted
                }, nil)

                // Set cache with calculated balance
                c.On("SetBalance", mock.Anything, "wallet-123", 0).Return(nil)
            },
            expectedBalance: 0,
        },
        {
            name:     "cache error - fallback to database calculation",
            walletID: "wallet-123",
            mockSetup: func(tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Cache error
                c.On("GetBalance", mock.Anything, "wallet-123").Return((*int)(nil), errors.New("cache connection error"))
            },
            expectedError: "cache connection error",
        },
        {
            name:     "database error during transaction list",
            walletID: "wallet-123",
            mockSetup: func(tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Cache miss
                c.On("GetBalance", mock.Anything, "wallet-123").Return((*int)(nil), nil)

                // Mock wallet lock
                unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
                c.On("Mutex", mock.Anything, "wallet-123").Return(unlockFunc, nil)

                // Database error
                tr.On("ListAllTransactions", mock.Anything, "wallet-123").Return(models.Transactions{}, errors.New("database connection error"))
            },
            expectedError: "database connection error",
        },
        {
            name:     "mutex lock error",
            walletID: "wallet-123",
            mockSetup: func(tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Cache miss
                c.On("GetBalance", mock.Anything, "wallet-123").Return((*int)(nil), nil)

                // Mutex lock error
                c.On("Mutex", mock.Anything, "wallet-123").Return(nil, errors.New("failed to acquire lock"))
            },
            expectedError: "failed to acquire lock",
        },
        {
            name:     "empty transaction list",
            walletID: "wallet-123",
            mockSetup: func(tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Cache miss
                c.On("GetBalance", mock.Anything, "wallet-123").Return((*int)(nil), nil)

                // Mock wallet lock
                unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
                c.On("Mutex", mock.Anything, "wallet-123").Return(unlockFunc, nil)

                // No transactions
                tr.On("ListAllTransactions", mock.Anything, "wallet-123").Return(models.Transactions{}, nil)

                // Set cache with zero balance
                c.On("SetBalance", mock.Anything, "wallet-123", 0).Return(nil)
            },
            expectedBalance: 0,
        },
        {
            name:     "cache set error - balance still returned",
            walletID: "wallet-123",
            mockSetup: func(tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Cache miss
                c.On("GetBalance", mock.Anything, "wallet-123").Return((*int)(nil), nil)

                // Mock wallet lock
                unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
                c.On("Mutex", mock.Anything, "wallet-123").Return(unlockFunc, nil)

                // Return transactions
                tr.On("ListAllTransactions", mock.Anything, "wallet-123").Return(models.Transactions{
                    {Amount: 1000, Type: string(types.TransactionTypeCredit), Status: string(types.TransactionStatusCompleted)},
                }, nil)

                // Cache set error (should not affect balance return)
                c.On("SetBalance", mock.Anything, "wallet-123", 1000).Return(errors.New("cache set error"))
            },
            expectedBalance: 1000, // Balance should still be calculated correctly even if cache set fails
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup mocks
            mockWalletRepo := mocks.NewMockWalletRepo(t)
            mockTransactionRepo := mocks.NewMockTransactionRepo(t)
            mockCache := mocks.NewMockCacheClient(t)

            tt.mockSetup(mockTransactionRepo, mockCache)

            // Create service
            service := NewService(mockWalletRepo, mockTransactionRepo, mockCache, time.Now)

            // Execute
            result, err := service.RunningBalance(context.Background(), tt.walletID)

            // Assert
            if tt.expectedError != "" {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectedError)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expectedBalance, result)
            }
        })
    }
}

func TestCreateTransaction(t *testing.T) {
    fixedTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
    
    tests := []struct {
        name          string
        request       models.CreateTransactionRequest
        mockSetup     func(*mocks.MockWalletRepo, *mocks.MockTransactionRepo, *mocks.MockCacheClient)
        expectedError string
        expectSuccess bool
    }{
        {
            name: "successful credit transaction",
            request: models.CreateTransactionRequest{
                WalletID:       "wallet-123",
                Amount:         1000,
                Type:           string(types.TransactionTypeCredit),
                IdempotencyKey: "idempotency-123",
            },
            mockSetup: func(wr *mocks.MockWalletRepo, tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Mock idempotency check
                unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
                c.On("Mutex", mock.Anything, "idempotency:idempotency-123").Return(unlockFunc, nil)
                c.On("GetIdempotentTransaction", mock.Anything, "idempotency-123").Return((*models.Transaction)(nil), nil)

                // Mock wallet retrieval - active wallet
                wallet := models.Wallet{
                    ID:     "wallet-123",
                    Status: string(types.WalletStatusActive),
                }
                wr.On("GetByID", mock.Anything, "wallet-123").Return(wallet, nil)

                // Mock wallet lock
                c.On("Mutex", mock.Anything, "wallet-123").Return(unlockFunc, nil)

                // Mock existing transactions (current balance: 500)
                tr.On("ListAllTransactions", mock.Anything, "wallet-123").Return(models.Transactions{
                    {Amount: 500, Type: string(types.TransactionTypeCredit), Status: string(types.TransactionStatusCompleted)},
                }, nil)

                // Mock transaction creation
                expectedTransaction := models.Transaction{
                    WalletID: "wallet-123",
                    Amount:   1000,
                    Type:     string(types.TransactionTypeCredit),
                    Status:   string(types.TransactionStatusPending),
                }
                tr.On("Create", mock.Anything, mock.MatchedBy(func(t models.Transaction) bool {
                    return t.WalletID == expectedTransaction.WalletID &&
                        t.Amount == expectedTransaction.Amount &&
                        t.Type == expectedTransaction.Type &&
                        t.Status == expectedTransaction.Status
                })).Return(expectedTransaction, nil)

                // Mock idempotency cache
                c.On("SetIdempotentTransaction", mock.Anything, "idempotency-123", expectedTransaction).Return(nil)

                // Mock balance cache update (credit transactions don't update cache immediately)
                c.On("SetBalance", mock.Anything, "wallet-123", 500).Return(nil)
            },
            expectSuccess: true,
        },
        {
            name: "successful debit transaction with sufficient balance",
            request: models.CreateTransactionRequest{
                WalletID:       "wallet-123",
                Amount:         300,
                Type:           string(types.TransactionTypeDebit),
                IdempotencyKey: "idempotency-456",
            },
            mockSetup: func(wr *mocks.MockWalletRepo, tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Mock idempotency check
                unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
                c.On("Mutex", mock.Anything, "idempotency:idempotency-456").Return(unlockFunc, nil)
                c.On("GetIdempotentTransaction", mock.Anything, "idempotency-456").Return((*models.Transaction)(nil), nil)

                // Mock wallet retrieval - active wallet
                wallet := models.Wallet{
                    ID:     "wallet-123",
                    Status: string(types.WalletStatusActive),
                }
                wr.On("GetByID", mock.Anything, "wallet-123").Return(wallet, nil)

                // Mock wallet lock
                c.On("Mutex", mock.Anything, "wallet-123").Return(unlockFunc, nil)

                // Mock existing transactions (current balance: 1000)
                tr.On("ListAllTransactions", mock.Anything, "wallet-123").Return(models.Transactions{
                    {Amount: 1000, Type: string(types.TransactionTypeCredit), Status: string(types.TransactionStatusCompleted)},
                }, nil)

                // Mock transaction creation
                expectedTransaction := models.Transaction{
                    WalletID: "wallet-123",
                    Amount:   300,
                    Type:     string(types.TransactionTypeDebit),
                    Status:   string(types.TransactionStatusPending),
                }
                tr.On("Create", mock.Anything, mock.MatchedBy(func(t models.Transaction) bool {
                    return t.WalletID == expectedTransaction.WalletID &&
                        t.Amount == expectedTransaction.Amount &&
                        t.Type == expectedTransaction.Type &&
                        t.Status == expectedTransaction.Status
                })).Return(expectedTransaction, nil)

                // Mock idempotency cache
                c.On("SetIdempotentTransaction", mock.Anything, "idempotency-456", expectedTransaction).Return(nil)

                // Mock balance cache update (debit transactions update cache: 1000 - 300 = 700)
                c.On("SetBalance", mock.Anything, "wallet-123", 700).Return(nil)
            },
            expectSuccess: true,
        },
        {
            name: "should not allow debit when there is no balance",
            request: models.CreateTransactionRequest{
                WalletID:       "wallet-123",
                Amount:         500,
                Type:           string(types.TransactionTypeDebit),
                IdempotencyKey: "idempotency-789",
            },
            mockSetup: func(wr *mocks.MockWalletRepo, tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Mock idempotency check
                unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
                c.On("Mutex", mock.Anything, "idempotency:idempotency-789").Return(unlockFunc, nil)
                c.On("GetIdempotentTransaction", mock.Anything, "idempotency-789").Return((*models.Transaction)(nil), nil)

                // Mock wallet retrieval - active wallet
                wallet := models.Wallet{
                    ID:     "wallet-123",
                    Status: string(types.WalletStatusActive),
                }
                wr.On("GetByID", mock.Anything, "wallet-123").Return(wallet, nil)

                // Mock wallet lock
                c.On("Mutex", mock.Anything, "wallet-123").Return(unlockFunc, nil)

                // Mock existing transactions (current balance: 0)
                tr.On("ListAllTransactions", mock.Anything, "wallet-123").Return(models.Transactions{}, nil)
            },
            expectedError: "insufficient funds",
        },
        {
            name: "should not allow debit when balance is insufficient",
            request: models.CreateTransactionRequest{
                WalletID:       "wallet-123",
                Amount:         1500,
                Type:           string(types.TransactionTypeDebit),
                IdempotencyKey: "idempotency-999",
            },
            mockSetup: func(wr *mocks.MockWalletRepo, tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Mock idempotency check
                unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
                c.On("Mutex", mock.Anything, "idempotency:idempotency-999").Return(unlockFunc, nil)
                c.On("GetIdempotentTransaction", mock.Anything, "idempotency-999").Return((*models.Transaction)(nil), nil)

                // Mock wallet retrieval - active wallet
                wallet := models.Wallet{
                    ID:     "wallet-123",
                    Status: string(types.WalletStatusActive),
                }
                wr.On("GetByID", mock.Anything, "wallet-123").Return(wallet, nil)

                // Mock wallet lock
                c.On("Mutex", mock.Anything, "wallet-123").Return(unlockFunc, nil)

                // Mock existing transactions (current balance: 1000)
                tr.On("ListAllTransactions", mock.Anything, "wallet-123").Return(models.Transactions{
                    {Amount: 1000, Type: string(types.TransactionTypeCredit), Status: string(types.TransactionStatusCompleted)},
                }, nil)
            },
            expectedError: "insufficient funds",
        },
        {
            name: "should not allow transactions on inactive wallets",
            request: models.CreateTransactionRequest{
                WalletID:       "wallet-456",
                Amount:         100,
                Type:           string(types.TransactionTypeCredit),
                IdempotencyKey: "idempotency-inactive",
            },
            mockSetup: func(wr *mocks.MockWalletRepo, tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Mock idempotency check
                unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
                c.On("Mutex", mock.Anything, "idempotency:idempotency-inactive").Return(unlockFunc, nil)
                c.On("GetIdempotentTransaction", mock.Anything, "idempotency-inactive").Return((*models.Transaction)(nil), nil)

                // Mock wallet retrieval - inactive wallet
                wallet := models.Wallet{
                    ID:     "wallet-456",
                    Status: string(types.WalletStatusInactive),
                }
                wr.On("GetByID", mock.Anything, "wallet-456").Return(wallet, nil)
            },
            expectedError: "cannot create transaction for non active wallets",
        },
        {
            name: "should not allow transactions on frozen wallets",
            request: models.CreateTransactionRequest{
                WalletID:       "wallet-789",
                Amount:         100,
                Type:           string(types.TransactionTypeCredit),
                IdempotencyKey: "idempotency-frozen",
            },
            mockSetup: func(wr *mocks.MockWalletRepo, tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Mock idempotency check
                unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
                c.On("Mutex", mock.Anything, "idempotency:idempotency-frozen").Return(unlockFunc, nil)
                c.On("GetIdempotentTransaction", mock.Anything, "idempotency-frozen").Return((*models.Transaction)(nil), nil)

                // Mock wallet retrieval - frozen wallet
                wallet := models.Wallet{
                    ID:     "wallet-789",
                    Status: string(types.WalletStatusFrozen),
                }
                wr.On("GetByID", mock.Anything, "wallet-789").Return(wallet, nil)
            },
            expectedError: "cannot create transaction for non active wallets",
        },
        {
            name: "idempotency - return existing transaction",
            request: models.CreateTransactionRequest{
                WalletID:       "wallet-123",
                Amount:         1000,
                Type:           string(types.TransactionTypeCredit),
                IdempotencyKey: "existing-key",
            },
            mockSetup: func(wr *mocks.MockWalletRepo, tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Mock idempotency check - existing transaction found
                unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
                c.On("Mutex", mock.Anything, "idempotency:existing-key").Return(unlockFunc, nil)
                
                existingTransaction := &models.Transaction{
                    ID:       "existing-txn-123",
                    WalletID: "wallet-123",
                    Amount:   1000,
                    Type:     string(types.TransactionTypeCredit),
                    Status:   string(types.TransactionStatusCompleted),
                }
                c.On("GetIdempotentTransaction", mock.Anything, "existing-key").Return(existingTransaction, nil)
            },
            expectSuccess: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup mocks
            mockWalletRepo := mocks.NewMockWalletRepo(t)
            mockTransactionRepo := mocks.NewMockTransactionRepo(t)
            mockCache := mocks.NewMockCacheClient(t)

            tt.mockSetup(mockWalletRepo, mockTransactionRepo, mockCache)

            // Create service
            service := NewService(mockWalletRepo, mockTransactionRepo, mockCache, func() time.Time { return fixedTime })

            // Execute
            result, err := service.CreateTransaction(context.Background(), tt.request)

            // Assert
            if tt.expectedError != "" {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectedError)
                assert.Empty(t, result.ID) // Should not return a valid transaction
            } else if tt.expectSuccess {
                assert.NoError(t, err)
                assert.NotEmpty(t, result.WalletID)
                assert.Equal(t, tt.request.Amount, result.Amount)
                assert.Equal(t, tt.request.Type, result.Type)
            }
        })
    }
}

func TestUpdateTransactionStatus(t *testing.T) {
    tests := []struct {
        name          string
        transactionID string
        newStatus     string
        mockSetup     func(*mocks.MockWalletRepo, *mocks.MockTransactionRepo, *mocks.MockCacheClient)
        expectedError string
        expectSuccess bool
    }{
        {
            name:          "credit transaction completed - balance should be updated in cache",
            transactionID: "txn-123",
            newStatus:     string(types.TransactionStatusCompleted),
            mockSetup: func(wr *mocks.MockWalletRepo, tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Mock transaction retrieval - pending credit transaction
                transaction := models.Transaction{
                    ID:       "txn-123",
                    WalletID: "wallet-123",
                    Amount:   1000,
                    Type:     string(types.TransactionTypeCredit),
                    Status:   string(types.TransactionStatusPending),
                }
                tr.On("GetByID", mock.Anything, "txn-123").Return(transaction, nil)

                // Mock wallet lock
                unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
                c.On("Mutex", mock.Anything, "wallet-123").Return(unlockFunc, nil)

                // Mock database transaction
                tr.On("Tx", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
                    fn := args.Get(1).(func(context.Context) error)
                    fn(context.Background())
                }).Return(nil)

                // Mock transaction update
                updatedTransaction := transaction
                updatedTransaction.Status = string(types.TransactionStatusCompleted)
                tr.On("Update", mock.Anything, mock.MatchedBy(func(t models.Transaction) bool {
                    return t.ID == "txn-123" && t.Status == string(types.TransactionStatusCompleted)
                })).Return(updatedTransaction, nil)

                // Mock cache balance check - cache exists with current balance 500
                currentBalance := 500
                c.On("GetBalance", mock.Anything, "wallet-123").Return(&currentBalance, nil)

                // Mock cache balance update - should add credit amount: 500 + 1000 = 1500
                c.On("SetBalance", mock.Anything, "wallet-123", 1500).Return(nil)
            },
            expectSuccess: true,
        },
        {
            name:          "debit transaction failed - balance should be updated in cache",
            transactionID: "txn-456",
            newStatus:     string(types.TransactionStatusFailed),
            mockSetup: func(wr *mocks.MockWalletRepo, tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Mock transaction retrieval - pending debit transaction
                transaction := models.Transaction{
                    ID:       "txn-456",
                    WalletID: "wallet-123",
                    Amount:   300,
                    Type:     string(types.TransactionTypeDebit),
                    Status:   string(types.TransactionStatusPending),
                }
                tr.On("GetByID", mock.Anything, "txn-456").Return(transaction, nil)

                // Mock wallet lock
                unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
                c.On("Mutex", mock.Anything, "wallet-123").Return(unlockFunc, nil)

                // Mock database transaction
                tr.On("Tx", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
                    fn := args.Get(1).(func(context.Context) error)
                    fn(context.Background())
                }).Return(nil)

                // Mock transaction update
                updatedTransaction := transaction
                updatedTransaction.Status = string(types.TransactionStatusFailed)
                tr.On("Update", mock.Anything, mock.MatchedBy(func(t models.Transaction) bool {
                    return t.ID == "txn-456" && t.Status == string(types.TransactionStatusFailed)
                })).Return(updatedTransaction, nil)

                // Mock cache balance check - cache exists with current balance 700 (after debit was applied)
                currentBalance := 700
                c.On("GetBalance", mock.Anything, "wallet-123").Return(&currentBalance, nil)

                // Mock cache balance update - should restore debit amount: 700 + 300 = 1000
                c.On("SetBalance", mock.Anything, "wallet-123", 1000).Return(nil)
            },
            expectSuccess: true,
        },
        {
            name:          "credit transaction failed - balance should be updated in cache",
            transactionID: "txn-789",
            newStatus:     string(types.TransactionStatusFailed),
            mockSetup: func(wr *mocks.MockWalletRepo, tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Mock transaction retrieval - pending credit transaction
                transaction := models.Transaction{
                    ID:       "txn-789",
                    WalletID: "wallet-123",
                    Amount:   500,
                    Type:     string(types.TransactionTypeCredit),
                    Status:   string(types.TransactionStatusPending),
                }
                tr.On("GetByID", mock.Anything, "txn-789").Return(transaction, nil)

                // Mock wallet lock
                unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
                c.On("Mutex", mock.Anything, "wallet-123").Return(unlockFunc, nil)

                // Mock database transaction
                tr.On("Tx", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
                    fn := args.Get(1).(func(context.Context) error)
                    fn(context.Background())
                }).Return(nil)

                // Mock transaction update
                updatedTransaction := transaction
                updatedTransaction.Status = string(types.TransactionStatusFailed)
                tr.On("Update", mock.Anything, mock.MatchedBy(func(t models.Transaction) bool {
                    return t.ID == "txn-789" && t.Status == string(types.TransactionStatusFailed)
                })).Return(updatedTransaction, nil)
            },
            expectSuccess: true,
        },
        {
            name:          "no cache update when balance not in cache",
            transactionID: "txn-999",
            newStatus:     string(types.TransactionStatusCompleted),
            mockSetup: func(wr *mocks.MockWalletRepo, tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Mock transaction retrieval - pending credit transaction
                transaction := models.Transaction{
                    ID:       "txn-999",
                    WalletID: "wallet-123",
                    Amount:   1000,
                    Type:     string(types.TransactionTypeCredit),
                    Status:   string(types.TransactionStatusPending),
                }
                tr.On("GetByID", mock.Anything, "txn-999").Return(transaction, nil)

                // Mock wallet lock
                unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
                c.On("Mutex", mock.Anything, "wallet-123").Return(unlockFunc, nil)

                // Mock database transaction
                tr.On("Tx", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
                    fn := args.Get(1).(func(context.Context) error)
                    fn(context.Background())
                }).Return(nil)

                // Mock transaction update
                updatedTransaction := transaction
                updatedTransaction.Status = string(types.TransactionStatusCompleted)
                tr.On("Update", mock.Anything, mock.MatchedBy(func(t models.Transaction) bool {
                    return t.ID == "txn-999" && t.Status == string(types.TransactionStatusCompleted)
                })).Return(updatedTransaction, nil)

                // Mock cache balance check - no cache entry
                c.On("GetBalance", mock.Anything, "wallet-123").Return((*int)(nil), nil)

                // No SetBalance call should be made when balance is not in cache
            },
            expectSuccess: true,
        },
        {
            name:          "same status - no update needed",
            transactionID: "txn-same",
            newStatus:     string(types.TransactionStatusCompleted),
            mockSetup: func(wr *mocks.MockWalletRepo, tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Mock transaction retrieval - already completed transaction
                transaction := models.Transaction{
                    ID:       "txn-same",
                    WalletID: "wallet-123",
                    Amount:   1000,
                    Type:     string(types.TransactionTypeCredit),
                    Status:   string(types.TransactionStatusCompleted),
                }
                tr.On("GetByID", mock.Anything, "txn-same").Return(transaction, nil)

                // No other mocks needed as the function should return early
            },
            expectSuccess: true,
        },
        {
            name:          "invalid status transition",
            transactionID: "txn-invalid",
            newStatus:     string(types.TransactionStatusPending),
            mockSetup: func(wr *mocks.MockWalletRepo, tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Mock transaction retrieval - completed transaction
                transaction := models.Transaction{
                    ID:       "txn-invalid",
                    WalletID: "wallet-123",
                    Amount:   1000,
                    Type:     string(types.TransactionTypeCredit),
                    Status:   string(types.TransactionStatusCompleted),
                }
                tr.On("GetByID", mock.Anything, "txn-invalid").Return(transaction, nil)

                // Mock wallet lock
                unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
                c.On("Mutex", mock.Anything, "wallet-123").Return(unlockFunc, nil)
            },
            expectedError: "invalid status transition from completed to pending",
        },
        {
            name:          "transaction not found",
            transactionID: "txn-not-found",
            newStatus:     string(types.TransactionStatusCompleted),
            mockSetup: func(wr *mocks.MockWalletRepo, tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Mock transaction retrieval - not found
                tr.On("GetByID", mock.Anything, "txn-not-found").Return(models.Transaction{}, errors.New("transaction not found"))
            },
            expectedError: "transaction not found",
        },
        {
            name:          "debit transaction completed - no cache update needed",
            transactionID: "txn-debit-completed",
            newStatus:     string(types.TransactionStatusCompleted),
            mockSetup: func(wr *mocks.MockWalletRepo, tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
                // Mock transaction retrieval - pending debit transaction
                transaction := models.Transaction{
                    ID:       "txn-debit-completed",
                    WalletID: "wallet-123",
                    Amount:   300,
                    Type:     string(types.TransactionTypeDebit),
                    Status:   string(types.TransactionStatusPending),
                }
                tr.On("GetByID", mock.Anything, "txn-debit-completed").Return(transaction, nil)

                // Mock wallet lock
                unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
                c.On("Mutex", mock.Anything, "wallet-123").Return(unlockFunc, nil)

                // Mock database transaction
                tr.On("Tx", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Run(func(args mock.Arguments) {
                    fn := args.Get(1).(func(context.Context) error)
                    fn(context.Background())
                }).Return(nil)

                // Mock transaction update
                updatedTransaction := transaction
                updatedTransaction.Status = string(types.TransactionStatusCompleted)
                tr.On("Update", mock.Anything, mock.MatchedBy(func(t models.Transaction) bool {
                    return t.ID == "txn-debit-completed" && t.Status == string(types.TransactionStatusCompleted)
                })).Return(updatedTransaction, nil)
            },
            expectSuccess: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup mocks
            mockWalletRepo := mocks.NewMockWalletRepo(t)
            mockTransactionRepo := mocks.NewMockTransactionRepo(t)
            mockCache := mocks.NewMockCacheClient(t)

            tt.mockSetup(mockWalletRepo, mockTransactionRepo, mockCache)

            // Create service
            service := NewService(mockWalletRepo, mockTransactionRepo, mockCache, time.Now)

            // Execute
            result, err := service.UpdateTransactionStatus(context.Background(), tt.transactionID, tt.newStatus)

            // Assert
            if tt.expectedError != "" {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectedError)
            } else if tt.expectSuccess {
                assert.NoError(t, err)
                assert.Equal(t, tt.transactionID, result.ID)
                if tt.newStatus != result.Status && tt.name != "same status - no update needed" {
                    assert.Equal(t, tt.newStatus, result.Status)
                }
            }
        })
    }
}
