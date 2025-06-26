package wallet

import (
	"context"
	"fmt"
)

func (cl *Client) GetTransactionByID(ctx context.Context, id string) (TransactionResponse, error) {
	var transaction TransactionResponse

	url := cl.buildUrl(fmt.Sprintf("/transactions/%s", id), nil)

	_, err := cl.httpClient.R().
		SetContext(ctx).
		SetResult(&transaction).
		Get(url)

	if err != nil {
		return TransactionResponse{}, fmt.Errorf("failed to get transaction by ID: %w", err)
	}

	return transaction, nil
}

func (cl *Client) GetTransactions(ctx context.Context, query ListTransactionsRequest) (TransactionsResponse, error) {
	var transactions TransactionsResponse

	url := cl.buildUrl("/transactions", query)

	_, err := cl.httpClient.R().
		SetContext(ctx).
		SetResult(&transactions).
		Get(url)

	if err != nil {
		return TransactionsResponse{}, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
}

func (cl *Client) CreateTransaction(ctx context.Context, req CreateTransactionRequest) (TransactionResponse, error) {
	var transaction TransactionResponse

	url := cl.buildUrl("/transactions", nil)

	_, err := cl.httpClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&transaction).
		Post(url)

	if err != nil {
		return TransactionResponse{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	return transaction, nil
}

func (cl *Client) UpdateTransactionStatus(ctx context.Context, id string, req UpdateTransactionStatusRequest) (
	TransactionResponse, error) {
	var transaction TransactionResponse

	url := cl.buildUrl(fmt.Sprintf("/transactions/%s/status", id), nil)

	_, err := cl.httpClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&transaction).
		Patch(url)

	if err != nil {
		return TransactionResponse{}, fmt.Errorf("failed to update transaction status: %w", err)
	}

	return transaction, nil
}
