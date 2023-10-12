package statistics

import (
	"context"
	"net/http"
)

type Client struct {
	Client http.Client
}

type IClient interface {
	GetMostLikelyAge(ctx context.Context, name string) (int, error)
	GetMostLikelyGender(ctx context.Context, name string) (string, error)
	GetMostLikelyNationality(ctx context.Context, name string) (string, error)
}

const (
	AgeStatisticsApiURI         = "https://api.agify.io"
	GenderStatisticsApiURI      = "https://api.genderize.io"
	NationalityStatisticsApiURI = "https://api.nationalize.io"
)

func NewClient() *Client {
	return &Client{
		Client: http.Client{},
	}
}
