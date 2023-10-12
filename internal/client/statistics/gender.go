package statistics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type GenderStatisticsResponse struct {
	Count       int     `json:"count,omitempty"`
	Name        string  `json:"name,omitempty"`
	Gender      string  `json:"gender,omitempty"`
	Probability float32 `json:"probability,omitempty"`
}

func (c Client) GetMostLikelyGender(ctx context.Context, name string) (string, error) {
	url := fmt.Sprintf("%s/?name=%s", GenderStatisticsApiURI, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var statistics GenderStatisticsResponse
	if err := json.NewDecoder(resp.Body).Decode(&statistics); err != nil {
		return "", err
	}

	return statistics.Gender, nil
}
