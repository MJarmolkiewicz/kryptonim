package rate_handler_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"kryptonim/handlers"
	"kryptonim/internal/mocks"
	"kryptonim/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetRates_BadRequest(t *testing.T) {
	//given
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/rates?currencies=USD,PLN", nil)

	rateService := new(mocks.RateServiceI)
	rateService.EXPECT().FetchRatesUSD().Return(nil, fmt.Errorf("no passe nada"))

	curHandler := handlers.NewCurrencyHandler(rateService)
	curHandler.GetFIATRates(c)

	assert.Equal(t, 400, w.Code)
}
func TestGetRates_MissingCurrenciesParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/rates", nil)

	rateService := new(mocks.RateServiceI)

	curHandler := handlers.NewCurrencyHandler(rateService)
	curHandler.GetFIATRates(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestGetRates_OneCurrencyOnly(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/rates?currencies=USD", nil)

	rateService := new(mocks.RateServiceI)

	curHandler := handlers.NewCurrencyHandler(rateService)
	curHandler.GetFIATRates(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
func TestGetRates_Success(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/rates?currencies=PLN,EUR", nil)

	mockResult := map[string]*big.Float{
		"PLN": big.NewFloat(21.22),
		"EUR": big.NewFloat(33.37),
	}

	rateService := new(mocks.RateServiceI)
	rateService.EXPECT().FetchRatesUSD().Return(mockResult, nil)

	curHandler := handlers.NewCurrencyHandler(rateService)
	curHandler.GetFIATRates(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.ExchangeRate
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	expected := []models.ExchangeRate([]models.ExchangeRate{
		{From: "PLN", To: "EUR", Rate: "1.572573"},
		{From: "EUR", To: "PLN", Rate: "0.635900"}})
	assert.Equal(t, expected, response)
}
