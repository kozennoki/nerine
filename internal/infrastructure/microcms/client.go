package microcms

import (
	"github.com/microcmsio/microcms-go-sdk"
)

type Client struct {
	client *microcms.Client
}

func NewClient(apiKey, serviceID string) *Client {
	client := microcms.New(serviceID, apiKey)
	return &Client{
		client: client,
	}
}