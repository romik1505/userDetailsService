package statistics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type AgeStatisticsResponse struct {
	Count int    `json:"count,omitempty"`
	Name  string `json:"name,omitempty"`
	Age   int    `json:"age,omitempty"`
}

func (c Client) GetMostLikelyAge(ctx context.Context, name string) (int, error) {
	url := fmt.Sprintf("%s/?name=%s", AgeStatisticsApiURI, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var statistics AgeStatisticsResponse
	if err := json.NewDecoder(resp.Body).Decode(&statistics); err != nil {
		return 0, err
	}

	return statistics.Age, nil
}
