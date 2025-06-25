package wallet

import (
	"context"
	"fmt"
)

func (cl *Client) GetWalletByID(ctx context.Context, id string) (WalletResponse, error) {
	var wallet WalletResponse

	url := cl.buildUrl(fmt.Sprintf("/wallets/%s", id), nil)

	_, err := cl.httpClient.R().
		SetContext(ctx).
		SetResult(&wallet).
		Get(url)

	if err != nil {
		return WalletResponse{}, fmt.Errorf("failed to get wallet by ID: %w", err)
	}

	return wallet, nil
}

func (cl *Client) UpdateWalletStatus(ctx context.Context, id string, req UpdateWalletStatusRequest) (WalletResponse, error) {
	var wallet WalletResponse

	url := cl.buildUrl(fmt.Sprintf("/wallets/%s/status", id), nil)

	_, err := cl.httpClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&wallet).
		Patch(url)

	if err != nil {
		return WalletResponse{}, err
	}

	return wallet, nil
}

func (cl *Client) CreateWallet(ctx context.Context, req CreateWalletRequest) (WalletResponse, error) {
	var wallet WalletResponse

	url := cl.buildUrl("/wallets", nil)

	_, err := cl.httpClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&wallet).
		Post(url)

	if err != nil {
		return WalletResponse{}, fmt.Errorf("failed to create wallet: %w", err)
	}

	return wallet, nil
}

func (cl *Client) ListWallets(ctx context.Context, query ListWalletsRequest) (WalletsResponse, error) {
	var wallets WalletsResponse

	url := cl.buildUrl("/wallets", query)

	_, err := cl.httpClient.R().
		SetContext(ctx).
		SetResult(&wallets).
		Get(url)

	if err != nil {
		return WalletsResponse{}, fmt.Errorf("failed to list wallets: %w", err)
	}

	return wallets, nil
}

func (cl *Client) GetWalletWithBalance(ctx context.Context, id string) (WalletResponse, error) {
	var wallet WalletResponse

	url := cl.buildUrl(fmt.Sprintf("/wallets/%s/balance", id), nil)

	_, err := cl.httpClient.R().
		SetContext(ctx).
		SetResult(&wallet).
		Get(url)

	if err != nil {
		return WalletResponse{}, fmt.Errorf("failed to get wallet with balance: %w", err)
	}

	return wallet, nil
}
