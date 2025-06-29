package models

type ExchangeRate struct {
	From string `json:"from"`
	To   string `json:"to"`
	Rate string `json:"rate"`
}

type CryptoRate struct {
	Decimals int
	Rate     string
}

var CryptoRates = map[string]CryptoRate{
	"BEER":  {Decimals: 18, Rate: "0.00002461"},
	"FLOKI": {Decimals: 18, Rate: "0.0001428"},
	"GATE":  {Decimals: 18, Rate: "6.87"},
	"USDT":  {Decimals: 6, Rate: "0.999"},
	"WBTC":  {Decimals: 8, Rate: "57037.22"},
}
