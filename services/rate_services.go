package services

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
)

type openRatesResponse struct {
	Rates map[string]float64 `json:"rates"`
}

type openExchangeFetcher struct{}

func NewRateService() *openExchangeFetcher {
	return &openExchangeFetcher{}
}

func (e openExchangeFetcher) FetchRatesUSD() (map[string]*big.Float, error) {
	AppID := ""
	url := fmt.Sprintf("https://openexchangerates.org/api/latest.json?app_id=%s", AppID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange rates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response from OpenExchangeRates: %d", resp.StatusCode)
	}

	var data openRatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to big.Float for precision
	rates := make(map[string]*big.Float)
	for currency, value := range data.Rates {
		rates[currency] = new(big.Float).SetPrec(128).SetFloat64(value)
	}

	return rates, nil
}
