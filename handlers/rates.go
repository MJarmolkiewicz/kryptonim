package handlers

import (
	"errors"
	"kryptonim/models"
	"math/big"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vargspjut/wlog"
)

//go:generate mockery --name=RateServiceI --output ../internal/mocks --with-expecter
type RateServiceI interface {
	FetchRatesUSD() (map[string]*big.Float, error)
}

type currencyHandler struct {
	rateService RateServiceI
}

func NewCurrencyHandler(rs RateServiceI) *currencyHandler {
	return &currencyHandler{rateService: rs}
}

func (ch *currencyHandler) GetFIATRates(c *gin.Context) {
	currenciesParam := c.Query("currencies")
	if currenciesParam == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	currencies := strings.Split(currenciesParam, ",")
	if len(currencies) < 2 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	usdRates, err := ch.rateService.FetchRatesUSD()
	if err != nil {
		wlog.Errorf("failed to fetch USD rates: %s", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	result, err := computeConversionMatrix(usdRates, currencies)
	if err != nil {
		wlog.Errorf("failed to compute conversion matrix: %s", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, result)
}

func computeConversionMatrix(usdRates map[string]*big.Float, currencies []string) ([]models.ExchangeRate, error) {
	for _, c := range currencies {
		if _, ok := usdRates[c]; !ok {
			return nil, errors.New("invalid currency: " + c)
		}
	}

	var result []models.ExchangeRate
	roundTo := big.NewFloat(1e6)

	for i := 0; i < len(currencies); i++ {
		from := currencies[i]
		for j := 0; j < len(currencies); j++ {
			if i == j {
				continue
			}
			to := currencies[j]

			//usdRates[to] / usdRates[from]
			rate := new(big.Float).Quo(usdRates[to], usdRates[from])

			// Round to 6 decimal places
			scaled := new(big.Float).Mul(rate, roundTo)
			intPart, _ := scaled.Int(nil) // truncate fractional part
			rounded := new(big.Float).Quo(new(big.Float).SetInt(intPart), roundTo)

			// Convert to string with fixed 6 decimals
			rateStr := rounded.Text('f', 6)

			result = append(result, models.ExchangeRate{
				From: from,
				To:   to,
				Rate: rateStr,
			})
		}
	}
	return result, nil
}

func (ch *currencyHandler) ExchangeCrypto(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	amountStr := c.Query("amount")

	if from == "" || to == "" || amountStr == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	fromRate, ok1 := models.CryptoRates[from]
	toRate, ok2 := models.CryptoRates[to]
	if !ok1 || !ok2 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	amount := new(big.Float)
	amount.SetPrec(128).SetString(amountStr)

	fromRateVal := new(big.Float).SetPrec(128)
	fromRateVal.SetString(fromRate.Rate)

	toRateVal := new(big.Float).SetPrec(128)
	toRateVal.SetString(toRate.Rate)

	// USD = amount * fromRate
	usd := new(big.Float).Mul(amount, fromRateVal)
	// result = USD / toRate
	result := new(big.Float).Quo(usd, toRateVal)

	// Zaokrąglij wg precision
	formatted := new(big.Float)
	formatted.SetPrec(128).SetString(result.Text('f', toRate.Decimals))

	c.JSON(http.StatusOK, gin.H{
		"from":   from,
		"to":     to,
		"amount": formatted,
	})
}
