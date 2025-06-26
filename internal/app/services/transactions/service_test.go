package transactions

import (
	"context"
	"errors"
	"testing"
	"time"

	"wallet/internal/app/models"
	"wallet/internal/app/services/transactions/mocks"
	"wallet/internal/util/pagination"
	"wallet/pkg/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetTransactionByID(t *testing.T) {
	tests := []struct {
		name                string
		transactionID       string
		mockSetup           func(*mocks.MockTransactionRepo)
		expectedTransaction models.Transaction
		expectedError       string
	}{
		{
			name:          "successful get transaction",
			transactionID: "txn-123",
			mockSetup: func(m *mocks.MockTransactionRepo) {
				m.On("GetByID", mock.Anything, "txn-123").Return(models.Transaction{
					ID:       "txn-123",
					WalletID: "wallet-123",
					Amount:   1000,
					Type:     string(types.TransactionTypeCredit),
					Status:   string(types.TransactionStatusCompleted),
				}, nil)
			},
			expectedTransaction: models.Transaction{
				ID:       "txn-123",
				WalletID: "wallet-123",
				Amount:   1000,
				Type:     string(types.TransactionTypeCredit),
				Status:   string(types.TransactionStatusCompleted),
			},
		},
		{
			name:          "transaction not found",
			transactionID: "txn-999",
			mockSetup: func(m *mocks.MockTransactionRepo) {
				m.On("GetByID", mock.Anything, "txn-999").Return(models.Transaction{}, errors.New("transaction not found"))
			},
			expectedError: "transaction not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockWalletRepo := mocks.NewMockWalletRepo(t)
			mockTransactionRepo := mocks.NewMockTransactionRepo(t)
			mockCache := mocks.NewMockCacheClient(t)

			tt.mockSetup(mockTransactionRepo)

			// Create service (correct parameter order: walletRepo, transactionRepo, cache, now)
			service := NewService(mockWalletRepo, mockTransactionRepo, mockCache, time.Now)

			// Execute
			result, err := service.GetTransactionByID(context.Background(), tt.transactionID)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTransaction, result)
			}
		})
	}
}

func TestService_ListTransactions(t *testing.T) {
	tests := []struct {
		name                 string
		query                models.QueryTransactions
		mockSetup            func(*mocks.MockTransactionRepo)
		expectedTransactions models.Transactions
		expectedPagination   *pagination.Pagination
		expectedError        string
	}{
		{
			name: "successful list transactions",
			query: models.QueryTransactions{
				WalletIDs: []string{"wallet-123"},
			},
			mockSetup: func(m *mocks.MockTransactionRepo) {
				m.On("List", mock.Anything, mock.MatchedBy(func(query models.QueryTransactions) bool {
					return len(query.WalletIDs) == 1 && query.WalletIDs[0] == "wallet-123"
				})).Return([]models.Transaction{
					{
						ID:       "txn-1",
						WalletID: "wallet-123",
						Amount:   1000,
						Type:     string(types.TransactionTypeCredit),
						Status:   string(types.TransactionStatusCompleted),
					},
					{
						ID:       "txn-2",
						WalletID: "wallet-123",
						Amount:   500,
						Type:     string(types.TransactionTypeDebit),
						Status:   string(types.TransactionStatusCompleted),
					},
				}, &pagination.Pagination{
					Page:    1,
					PerPage: 10,
					Total:   2,
				}, nil)
			},
			expectedTransactions: models.Transactions{
				{
					ID:       "txn-1",
					WalletID: "wallet-123",
					Amount:   1000,
					Type:     string(types.TransactionTypeCredit),
					Status:   string(types.TransactionStatusCompleted),
				},
				{
					ID:       "txn-2",
					WalletID: "wallet-123",
					Amount:   500,
					Type:     string(types.TransactionTypeDebit),
					Status:   string(types.TransactionStatusCompleted),
				},
			},
			expectedPagination: &pagination.Pagination{
				Page:    1,
				PerPage: 10,
				Total:   2,
			},
		},
		{
			name: "database error",
			query: models.QueryTransactions{
				WalletIDs: []string{"wallet-999"},
			},
			mockSetup: func(m *mocks.MockTransactionRepo) {
				m.On("List", mock.Anything, mock.Anything).Return([]models.Transaction{}, &pagination.Pagination{}, errors.New("database error"))
			},
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockWalletRepo := mocks.NewMockWalletRepo(t)
			mockTransactionRepo := mocks.NewMockTransactionRepo(t)
			mockCache := mocks.NewMockCacheClient(t)

			tt.mockSetup(mockTransactionRepo)

			// Create service (correct parameter order: walletRepo, transactionRepo, cache, now)
			service := NewService(mockWalletRepo, mockTransactionRepo, mockCache, time.Now)

			// Execute
			result, pag, err := service.ListTransactions(context.Background(), tt.query)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTransactions, result)
				assert.Equal(t, tt.expectedPagination, pag)
			}
		})
	}
}

func TestService_CreateTransaction(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	nowFunc := func() time.Time { return fixedTime }

	tests := []struct {
		name                string
		request             models.CreateTransactionRequest
		mockSetup           func(*mocks.MockWalletRepo, *mocks.MockTransactionRepo, *mocks.MockCacheClient)
		expectedTransaction models.Transaction
		expectedError       string
	}{
		{
			name: "successful transaction creation",
			request: models.CreateTransactionRequest{
				WalletID:       "wallet-123",
				Amount:         1000,
				Type:           string(types.TransactionTypeCredit),
				IdempotencyKey: "idempotency-key-123",
			},
			mockSetup: func(wr *mocks.MockWalletRepo, tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
				// Mock idempotency lock
				unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
				c.On("Mutex", mock.Anything, "idempotency:idempotency-key-123").Return(unlockFunc, nil)

				// Check idempotency - no existing transaction
				c.On("GetIdempotentTransaction", mock.Anything, "idempotency-key-123").Return((*models.Transaction)(nil), nil)

				// Check wallet exists
				wr.On("GetByID", mock.Anything, "wallet-123").Return(models.Wallet{
					ID:       "wallet-123",
					OwnerID:  "owner-123",
					Currency: "USD",
					Status:   "active",
				}, nil)

				// Mock wallet lock
				c.On("Mutex", mock.Anything, "wallet-123").Return(unlockFunc, nil)

				// Mock balance calculation
				tr.On("ListAllTransactions", mock.Anything, "wallet-123").Return(models.Transactions{
					{Amount: 500, Type: string(types.TransactionTypeCredit), Status: string(types.TransactionStatusCompleted)},
				}, nil)

				// Create transaction
				tr.On("Create", mock.Anything, mock.MatchedBy(func(txn models.Transaction) bool {
					return txn.WalletID == "wallet-123" && txn.Amount == 1000 && txn.ID != ""
				})).Return(models.Transaction{
					ID:       "01HQXC8ZJQN5T1FNBC123456",
					WalletID: "wallet-123",
					Amount:   1000,
					Type:     string(types.TransactionTypeCredit),
					Status:   string(types.TransactionStatusPending),
				}, nil)

				// Set cache for idempotency and balance
				c.On("SetIdempotentTransaction", mock.Anything, "idempotency-key-123", mock.Anything).Return(nil)
				c.On("SetBalance", mock.Anything, "wallet-123", 500).Return(nil)
			},
			expectedTransaction: models.Transaction{
				ID:       "01HQXC8ZJQN5T1FNBC123456",
				WalletID: "wallet-123",
				Amount:   1000,
				Type:     string(types.TransactionTypeCredit),
				Status:   string(types.TransactionStatusPending),
			},
		},
		{
			name: "duplicate idempotency key",
			request: models.CreateTransactionRequest{
				WalletID:       "wallet-123",
				Amount:         1000,
				Type:           string(types.TransactionTypeCredit),
				IdempotencyKey: "existing-key",
			},
			mockSetup: func(wr *mocks.MockWalletRepo, tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
				// Mock idempotency lock
				unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
				c.On("Mutex", mock.Anything, "idempotency:existing-key").Return(unlockFunc, nil)

				// Return existing transaction for idempotency key
				existingTxn := models.Transaction{
					ID:       "existing-txn",
					WalletID: "wallet-123",
					Amount:   1000,
					Type:     string(types.TransactionTypeCredit),
					Status:   string(types.TransactionStatusCompleted),
				}
				c.On("GetIdempotentTransaction", mock.Anything, "existing-key").Return(&existingTxn, nil)
			},
			expectedTransaction: models.Transaction{
				ID:       "existing-txn",
				WalletID: "wallet-123",
				Amount:   1000,
				Type:     string(types.TransactionTypeCredit),
				Status:   string(types.TransactionStatusCompleted),
			},
		},
		{
			name: "wallet not found",
			request: models.CreateTransactionRequest{
				WalletID:       "wallet-999",
				Amount:         1000,
				Type:           string(types.TransactionTypeCredit),
				IdempotencyKey: "idempotency-key-456",
			},
			mockSetup: func(wr *mocks.MockWalletRepo, tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
				// Mock idempotency lock
				unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
				c.On("Mutex", mock.Anything, "idempotency:idempotency-key-456").Return(unlockFunc, nil)

				// Check idempotency - no existing transaction
				c.On("GetIdempotentTransaction", mock.Anything, "idempotency-key-456").Return((*models.Transaction)(nil), nil)

				// Wallet not found
				wr.On("GetByID", mock.Anything, "wallet-999").Return(models.Wallet{}, errors.New("wallet not found"))
			},
			expectedError: "wallet not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockWalletRepo := mocks.NewMockWalletRepo(t)
			mockTransactionRepo := mocks.NewMockTransactionRepo(t)
			mockCache := mocks.NewMockCacheClient(t)

			tt.mockSetup(mockWalletRepo, mockTransactionRepo, mockCache)

			// Create service (correct parameter order: walletRepo, transactionRepo, cache, now)
			service := NewService(mockWalletRepo, mockTransactionRepo, mockCache, nowFunc)

			// Execute
			result, err := service.CreateTransaction(context.Background(), tt.request)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTransaction.WalletID, result.WalletID)
				assert.Equal(t, tt.expectedTransaction.Amount, result.Amount)
				assert.Equal(t, tt.expectedTransaction.Type, result.Type)
				assert.NotEmpty(t, result.ID)
			}
		})
	}
}

func TestService_RunningBalance(t *testing.T) {
	tests := []struct {
		name            string
		walletID        string
		mockSetup       func(*mocks.MockTransactionRepo, *mocks.MockCacheClient)
		expectedBalance int
		expectedError   string
	}{
		{
			name:     "successful balance calculation from cache",
			walletID: "wallet-123",
			mockSetup: func(tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
				balance := 1500
				c.On("GetBalance", mock.Anything, "wallet-123").Return(&balance, nil)
			},
			expectedBalance: 1500,
		},
		{
			name:     "successful balance calculation from database",
			walletID: "wallet-123",
			mockSetup: func(tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
				// Cache miss
				c.On("GetBalance", mock.Anything, "wallet-123").Return((*int)(nil), nil)

				// Mock wallet lock
				unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
				c.On("Mutex", mock.Anything, "wallet-123").Return(unlockFunc, nil)

				// Calculate from transactions - return transactions with balance: 1000 - 300 + 800 = 1500
				tr.On("ListAllTransactions", mock.Anything, "wallet-123").Return(models.Transactions{
					{Amount: 1000, Type: string(types.TransactionTypeCredit), Status: string(types.TransactionStatusCompleted)},
					{Amount: 300, Type: string(types.TransactionTypeDebit), Status: string(types.TransactionStatusCompleted)},
					{Amount: 800, Type: string(types.TransactionTypeCredit), Status: string(types.TransactionStatusCompleted)},
				}, nil)

				// Set cache
				c.On("SetBalance", mock.Anything, "wallet-123", 1500).Return(nil)
			},
			expectedBalance: 1500,
		},
		{
			name:     "database error",
			walletID: "wallet-999",
			mockSetup: func(tr *mocks.MockTransactionRepo, c *mocks.MockCacheClient) {
				c.On("GetBalance", mock.Anything, "wallet-999").Return((*int)(nil), nil)

				// Mock wallet lock
				unlockFunc := func(ctx context.Context) (bool, error) { return true, nil }
				c.On("Mutex", mock.Anything, "wallet-999").Return(unlockFunc, nil)

				tr.On("ListAllTransactions", mock.Anything, "wallet-999").Return(models.Transactions{}, errors.New("database error"))
			},
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockWalletRepo := mocks.NewMockWalletRepo(t)
			mockTransactionRepo := mocks.NewMockTransactionRepo(t)
			mockCache := mocks.NewMockCacheClient(t)

			tt.mockSetup(mockTransactionRepo, mockCache)

			// Create service (correct parameter order: walletRepo, transactionRepo, cache, now)
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
