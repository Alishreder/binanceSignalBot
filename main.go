package main

import (
	`encoding/json`
	`fmt`
	`io`
	`net/http`
	`strconv`
	`time`
)

const host = "https://api.binance.com"
const getProductsUrl = "https://www.binance.com/exchange-api/v2/public/asset-service/product/get-products"

const endpointGetPriceChange = "/api/v3/ticker"

const USDT = "USDT"
const marketCapLimit = 100000000

type responseResult struct {
	Code          string      `json:"code"`
	Message       interface{} `json:"message"`
	MessageDetail interface{} `json:"messageDetail"`
	Data          []pairInfo  `json:"data"`
	Success       bool        `json:"success"`
}

type pairInfo struct {
	S    string      `json:"s"`
	St   string      `json:"st"`
	B    string      `json:"b"`
	Q    string      `json:"q"`
	Ba   string      `json:"ba"`
	Qa   string      `json:"qa"`
	I    string      `json:"i"`
	Ts   string      `json:"ts"`
	An   string      `json:"an"`
	Qn   string      `json:"qn"`
	O    string      `json:"o"`
	H    string      `json:"h"`
	L    string      `json:"l"`
	C    string      `json:"c"`
	V    string      `json:"v"`
	Qv   string      `json:"qv"`
	Y    int         `json:"y"`
	As   float64     `json:"as"`
	Pm   string      `json:"pm"`
	Pn   string      `json:"pn"`
	Cs   int         `json:"cs"`
	Tags []string    `json:"tags"`
	Pom  bool        `json:"pom"`
	Pomt interface{} `json:"pomt"`
	Lc   bool        `json:"lc"`
	G    bool        `json:"g"`
	Sd   bool        `json:"sd"`
	R    bool        `json:"r"`
	Hd   bool        `json:"hd"`
	Etf  bool        `json:"etf"`
}

type pairData struct {
	Symbol string
	Prices []float64
}

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

func main() {

	pairs := getPairs()

	for {


		for _, pair := range pairs {
			priceChangePercentage, err := getTokenPriceChange(pair.Symbol, 1)
			if err != nil {
				continue
			}
			if priceChangePercentage > 4 {
				fmt.Printf("AAAAALLLLLEEEEERRRRRTTTTT %s change price on %f per hour\n", pair.Symbol, priceChangePercentage)
			}

		}
		time.Sleep(time.Hour)

	}
}

func getTokenPriceChange(symbol string, hours int) (float64, error) {
	url := host + endpointGetPriceChange + "?symbol=" + symbol + "&windowSize=" + strconv.Itoa(hours) + "h"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return 0, fmt.Errorf("can't create request: %s", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return 0, fmt.Errorf("can't read response body: %s", err)
	}

	if resp.StatusCode == 418 || resp.StatusCode == 429 {
		err = fmt.Errorf("you're temporary baned, retry after %s seconds", resp.Header.Get("Retry-After"))
		panic(err)
	}

	var response getPriceChangeResult

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err, string(body))
		return 0, fmt.Errorf("can't unmarshal body: %s", err)
	}

	priceChange, err := strconv.ParseFloat(response.PriceChangePercent, 32)
	if err != nil {
		fmt.Println(err, response.PriceChangePercent)
		return 0, fmt.Errorf("can't convert priceChange to float: %s", err)
	}

	return priceChange, nil
}

func getPairs() []pairData {
	resp, err := http.Get(getProductsUrl)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var response responseResult

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var result []pairData

	for _, pair := range response.Data {

		if pair.Q != USDT {
			continue
		}

		if !isMarketCapValid(pair.C, pair.Cs) {
			continue
		}

		currPair := pairData{
			Symbol: pair.S,
		}

		result = append(result, currPair)
	}

	return result
}

func isMarketCapValid(c string, cs int) bool {
	currentPrice, err := strconv.ParseFloat(c, 32)
	if err != nil {
		fmt.Println(err)
		return false
	}

	marketCap := float64(cs) * currentPrice

	if marketCap < marketCapLimit {
		return false
	}

	return true
}
