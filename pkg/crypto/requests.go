package crypto

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

const codeBrokeRequestRateLimit = 429
const codeIPHasBeenBanned = 418

const host = "https://api.binance.com"
const getProductsUrl = "https://www.binance.com/exchange-api/v2/public/asset-service/product/get-products"

const endpointGetPriceChange = "/api/v3/ticker"

const USDT = "USDT"
const marketCapLimit = 100000000

type getPriceChangeResult struct {
	Symbol             string `json:"symbol"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	WeightedAvgPrice   string `json:"weightedAvgPrice"`
	OpenPrice          string `json:"openPrice"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	LastPrice          string `json:"lastPrice"`
	Volume             string `json:"volume"`
	QuoteVolume        string `json:"quoteVolume"`
	OpenTime           int64  `json:"openTime"`
	CloseTime          int64  `json:"closeTime"`
	FirstId            int    `json:"firstId"`
	LastId             int    `json:"lastId"`
	Count              int    `json:"count"`
}

func getTokenPriceChange(symbol string, windowSize string) (float64, error) {
	url := host + endpointGetPriceChange + "?symbol=" + symbol + "&windowSize=" + windowSize

	body, err := doHTTPGet(url)
	if err != nil {
		log.Printf("error while trying to get token %s price: %s", symbol, err.Error())
		return 0, fmt.Errorf("error: %w while doing request to %s", err, url)
	}

	var response getPriceChangeResult

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error while trying to unmarshal token %s price data, error: %s, body: %s", symbol, err.Error(), string(body))
		return 0, fmt.Errorf("can't unmarshal body: %w", err)
	}

	priceChange, err := strconv.ParseFloat(response.PriceChangePercent, 32)
	if err != nil {
		log.Printf("error while trying parce string to float for %s price, err: %s", symbol, err.Error())
		return 0, fmt.Errorf("can't convert priceChange to float: %w", err)
	}

	return priceChange, nil
}

func doHTTPGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read response body: %w", err)
	}

	if resp.StatusCode == codeBrokeRequestRateLimit || resp.StatusCode == codeIPHasBeenBanned {
		err = fmt.Errorf("you're temporary baned, retry after %s seconds", resp.Header.Get("Retry-After"))
		return nil, err
	}

	return body, nil
}
