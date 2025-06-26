package wallets

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"wallet/internal/app/models"
	"wallet/internal/app/services/wallets/mocks"
	"wallet/internal/util/pagination"
)

func TestService_GetWalletByID(t *testing.T) {
	tests := []struct {
		name           string
		walletID       string
		mockSetup      func(*mocks.MockWalletDB)
		expectedWallet models.Wallet
		expectedError  string
	}{
		{
			name:     "successful get wallet",
			walletID: "wallet-123",
			mockSetup: func(m *mocks.MockWalletDB) {
				m.On("GetByID", mock.Anything, "wallet-123").Return(models.Wallet{
					ID:       "wallet-123",
					OwnerID:  "owner-123",
					Currency: "USD",
					Status:   "active",
				}, nil)
			},
			expectedWallet: models.Wallet{
				ID:       "wallet-123",
				OwnerID:  "owner-123",
				Currency: "USD",
				Status:   "active",
			},
		},
		{
			name:     "wallet not found",
			walletID: "wallet-999",
			mockSetup: func(m *mocks.MockWalletDB) {
				m.On("GetByID", mock.Anything, "wallet-999").Return(models.Wallet{}, errors.New("wallet not found"))
			},
			expectedError: "wallet not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockWalletDB := mocks.NewMockWalletDB(t)
			mockTransactionService := mocks.NewMockTransactionService(t)
			mockCache := mocks.NewMockCache(t)

			tt.mockSetup(mockWalletDB)

			// Create service
			service := NewService(mockTransactionService, mockWalletDB, mockCache, time.Now)

			// Execute
			result, err := service.GetWalletByID(context.Background(), tt.walletID)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedWallet, result)
			}
		})
	}
}

func TestService_CreateWallet(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	nowFunc := func() time.Time { return fixedTime }

	tests := []struct {
		name           string
		request        models.CreateWalletRequest
		mockSetup      func(*mocks.MockWalletDB)
		expectedWallet models.Wallet
		expectedError  string
	}{
		{
			name: "successful wallet creation",
			request: models.CreateWalletRequest{
				OwnerID:  "owner-123",
				Currency: "USD",
				Status:   "active",
			},
			mockSetup: func(m *mocks.MockWalletDB) {
				// First, check if wallet already exists (should return empty list)
				m.On("List", mock.Anything, mock.MatchedBy(func(query models.QueryWallets) bool {
					return len(query.OwnerIDs) == 1 && query.OwnerIDs[0] == "owner-123" &&
						len(query.Currencies) == 1 && query.Currencies[0] == "USD"
				})).Return([]models.Wallet{}, &pagination.Pagination{}, nil)

				// Then create the wallet
				m.On("Create", mock.Anything, mock.MatchedBy(func(wallet models.Wallet) bool {
					return wallet.OwnerID == "owner-123" && wallet.Currency == "USD" && wallet.ID != ""
				})).Return(models.Wallet{
					ID:       "01HQXC8ZJQN5T1FNBC123456",
					OwnerID:  "owner-123",
					Currency: "USD",
					Status:   "active",
				}, nil)
			},
			expectedWallet: models.Wallet{
				ID:       "01HQXC8ZJQN5T1FNBC123456",
				OwnerID:  "owner-123",
				Currency: "USD",
				Status:   "active",
			},
		},
		{
			name: "wallet already exists",
			request: models.CreateWalletRequest{
				OwnerID:  "owner-123",
				Currency: "USD",
				Status:   "active",
			},
			mockSetup: func(m *mocks.MockWalletDB) {
				// Return existing wallet
				m.On("List", mock.Anything, mock.MatchedBy(func(query models.QueryWallets) bool {
					return len(query.OwnerIDs) == 1 && query.OwnerIDs[0] == "owner-123" &&
						len(query.Currencies) == 1 && query.Currencies[0] == "USD"
				})).Return([]models.Wallet{
					{ID: "existing-wallet", OwnerID: "owner-123", Currency: "USD"},
				}, &pagination.Pagination{}, nil)
			},
			expectedError: "wallet already exists for this owner and currency",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockWalletDB := mocks.NewMockWalletDB(t)
			mockTransactionService := mocks.NewMockTransactionService(t)
			mockCache := mocks.NewMockCache(t)

			tt.mockSetup(mockWalletDB)

			// Create service
			service := NewService(mockTransactionService, mockWalletDB, mockCache, nowFunc)

			// Execute
			result, err := service.CreateWallet(context.Background(), tt.request)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedWallet.OwnerID, result.OwnerID)
				assert.Equal(t, tt.expectedWallet.Currency, result.Currency)
				assert.Equal(t, tt.expectedWallet.Status, result.Status)
				assert.NotEmpty(t, result.ID)
			}
		})
	}
}

func TestService_UpdateWalletStatus(t *testing.T) {
	tests := []struct {
		name           string
		walletID       string
		newStatus      string
		mockSetup      func(*mocks.MockWalletDB)
		expectedWallet models.Wallet
		expectedError  string
	}{
		{
			name:      "successful status update",
			walletID:  "wallet-123",
			newStatus: "frozen",
			mockSetup: func(m *mocks.MockWalletDB) {
				// Get wallet
				m.On("GetByID", mock.Anything, "wallet-123").Return(models.Wallet{
					ID:       "wallet-123",
					OwnerID:  "owner-123",
					Currency: "USD",
					Status:   "active",
				}, nil)

				// Update wallet
				m.On("Update", mock.Anything, mock.MatchedBy(func(wallet models.Wallet) bool {
					return wallet.ID == "wallet-123" && wallet.Status == "frozen"
				})).Return(models.Wallet{
					ID:       "wallet-123",
					OwnerID:  "owner-123",
					Currency: "USD",
					Status:   "frozen",
				}, nil)
			},
			expectedWallet: models.Wallet{
				ID:       "wallet-123",
				OwnerID:  "owner-123",
				Currency: "USD",
				Status:   "frozen",
			},
		},
		{
			name:      "status unchanged - no update needed",
			walletID:  "wallet-123",
			newStatus: "active",
			mockSetup: func(m *mocks.MockWalletDB) {
				// Get wallet
				m.On("GetByID", mock.Anything, "wallet-123").Return(models.Wallet{
					ID:       "wallet-123",
					OwnerID:  "owner-123",
					Currency: "USD",
					Status:   "active",
				}, nil)
				// No Update call expected since status is the same
			},
			expectedWallet: models.Wallet{
				ID:       "wallet-123",
				OwnerID:  "owner-123",
				Currency: "USD",
				Status:   "active",
			},
		},
		{
			name:      "wallet not found",
			walletID:  "wallet-999",
			newStatus: "frozen",
			mockSetup: func(m *mocks.MockWalletDB) {
				m.On("GetByID", mock.Anything, "wallet-999").Return(models.Wallet{}, errors.New("wallet not found"))
			},
			expectedError: "wallet not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockWalletDB := mocks.NewMockWalletDB(t)
			mockTransactionService := mocks.NewMockTransactionService(t)
			mockCache := mocks.NewMockCache(t)

			tt.mockSetup(mockWalletDB)

			// Create service
			service := NewService(mockTransactionService, mockWalletDB, mockCache, time.Now)

			// Execute
			result, err := service.UpdateWalletStatus(context.Background(), tt.walletID, tt.newStatus)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedWallet, result)
			}
		})
	}
}

func TestService_GetWalletWithBalance(t *testing.T) {
	tests := []struct {
		name           string
		walletID       string
		mockSetup      func(*mocks.MockWalletDB, *mocks.MockTransactionService)
		expectedWallet models.Wallet
		expectedError  string
	}{
		{
			name:     "successful get wallet with balance",
			walletID: "wallet-123",
			mockSetup: func(m *mocks.MockWalletDB, ts *mocks.MockTransactionService) {
				// Get wallet
				m.On("GetByID", mock.Anything, "wallet-123").Return(models.Wallet{
					ID:       "wallet-123",
					OwnerID:  "owner-123",
					Currency: "USD",
					Status:   "active",
				}, nil)

				// Get balance
				ts.On("RunningBalance", mock.Anything, "wallet-123").Return(1000, nil)
			},
			expectedWallet: models.Wallet{
				ID:       "wallet-123",
				OwnerID:  "owner-123",
				Currency: "USD",
				Status:   "active",
				Balance:  intPtr(1000),
			},
		},
		{
			name:     "wallet not found",
			walletID: "wallet-999",
			mockSetup: func(m *mocks.MockWalletDB, ts *mocks.MockTransactionService) {
				m.On("GetByID", mock.Anything, "wallet-999").Return(models.Wallet{}, errors.New("wallet not found"))
			},
			expectedError: "wallet not found",
		},
		{
			name:     "balance calculation error",
			walletID: "wallet-123",
			mockSetup: func(m *mocks.MockWalletDB, ts *mocks.MockTransactionService) {
				// Get wallet
				m.On("GetByID", mock.Anything, "wallet-123").Return(models.Wallet{
					ID:       "wallet-123",
					OwnerID:  "owner-123",
					Currency: "USD",
					Status:   "active",
				}, nil)

				// Balance calculation fails
				ts.On("RunningBalance", mock.Anything, "wallet-123").Return(0, errors.New("balance calculation failed"))
			},
			expectedError: "balance calculation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockWalletDB := mocks.NewMockWalletDB(t)
			mockTransactionService := mocks.NewMockTransactionService(t)
			mockCache := mocks.NewMockCache(t)

			tt.mockSetup(mockWalletDB, mockTransactionService)

			// Create service
			service := NewService(mockTransactionService, mockWalletDB, mockCache, time.Now)

			// Execute
			result, err := service.GetWalletWithBalance(context.Background(), tt.walletID)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedWallet, result)
			}
		})
	}
}

func TestService_ListWallets(t *testing.T) {
	tests := []struct {
		name            string
		query           models.QueryWallets
		mockSetup       func(*mocks.MockWalletDB)
		expectedWallets models.Wallets
		expectedPaging  *pagination.Pagination
		expectedError   string
	}{
		{
			name: "successful list wallets",
			query: models.QueryWallets{
				OwnerIDs: []string{"owner-123"},
			},
			mockSetup: func(m *mocks.MockWalletDB) {
				m.On("List", mock.Anything, mock.MatchedBy(func(query models.QueryWallets) bool {
					return len(query.OwnerIDs) == 1 && query.OwnerIDs[0] == "owner-123"
				})).Return([]models.Wallet{
					{ID: "wallet-1", OwnerID: "owner-123", Currency: "USD", Status: "active"},
					{ID: "wallet-2", OwnerID: "owner-123", Currency: "EUR", Status: "active"},
				}, &pagination.Pagination{
					Total:      2,
					Count:      2,
					Page:       1,
					PerPage:    10,
					TotalPages: 1,
				}, nil)
			},
			expectedWallets: models.Wallets{
				{ID: "wallet-1", OwnerID: "owner-123", Currency: "USD", Status: "active"},
				{ID: "wallet-2", OwnerID: "owner-123", Currency: "EUR", Status: "active"},
			},
			expectedPaging: &pagination.Pagination{
				Total:      2,
				Count:      2,
				Page:       1,
				PerPage:    10,
				TotalPages: 1,
			},
		},
		{
			name:  "database error",
			query: models.QueryWallets{},
			mockSetup: func(m *mocks.MockWalletDB) {
				m.On("List", mock.Anything, mock.Anything).Return(nil, nil, errors.New("database error"))
			},
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockWalletDB := mocks.NewMockWalletDB(t)
			mockTransactionService := mocks.NewMockTransactionService(t)
			mockCache := mocks.NewMockCache(t)

			tt.mockSetup(mockWalletDB)

			// Create service
			service := NewService(mockTransactionService, mockWalletDB, mockCache, time.Now)

			// Execute
			wallets, paging, err := service.ListWallets(context.Background(), tt.query)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedWallets, wallets)
				assert.Equal(t, tt.expectedPaging, paging)
			}
		})
	}
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}
