package wallet

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
)

type NoContentResponse struct{}

type Client struct {
	baseURL    string
	apiVersion string
	httpClient *resty.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		apiVersion: "v1",
		httpClient: resty.New(),
	}
}

func (cl *Client) buildUrl(path string, queryParams any) string {
	queries, err := query.Values(queryParams)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s/%s%s?%s", cl.baseURL, cl.apiVersion, path, queries.Encode())
}
