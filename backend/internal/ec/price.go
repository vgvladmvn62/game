package ec

import (
	"encoding/json"
	"strconv"
)

// Price stores details about the price.
type Price struct {
	Currency       string `json:"currency"`
	Value          int64  `json:"value"`
	FormattedValue string `json:"formatted_value"`
}

type priceJSON struct {
	Currency       string  `json:"currencyIso"`
	Value          float64 `json:"value"`
	FormattedValue string  `json:"formattedValue"`
}

// UnmarshalJSON provides custom unmarshal for Price.
func (p *Price) UnmarshalJSON(data []byte) error {
	var priceJSON priceJSON

	if err := json.Unmarshal(data, &priceJSON); err != nil {
		return err
	}

	p.Currency = priceJSON.Currency
	p.Value = int64(priceJSON.Value * 100)
	p.FormattedValue = priceJSON.FormattedValue

	return nil
}

// String returns Price string representation.
func (p *Price) String() string {
	return "Currency: " + p.Currency + "\n" +
		"Value: " + strconv.Itoa(int(p.Value)) + "\n" +
		"Formatted value: " + p.FormattedValue

}
