package statistics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type NationalityStatisticsResponse struct {
	Count   int    `json:"count,omitempty"`
	Name    string `json:"name,omitempty"`
	Country []struct {
		CountryID   string  `json:"country_id,omitempty"`
		Probability float32 `json:"probability,omitempty"`
	} `json:"country,omitempty"`
}

func (c Client) GetMostLikelyNationality(ctx context.Context, name string) (string, error) {
	url := fmt.Sprintf("%s/?name=%s", NationalityStatisticsApiURI, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var statistics NationalityStatisticsResponse
	if err := json.NewDecoder(resp.Body).Decode(&statistics); err != nil {
		return "", err
	}

	if len(statistics.Country) == 0 {
		return "", fmt.Errorf("empty country array")
	}

	return statistics.Country[0].CountryID, nil
}
