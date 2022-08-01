package crypto

import (
	"encoding/json"
	"fmt"
	"strconv"
)

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

func getPairs() []pairData {

	body, err := doHTTPGet(getProductsUrl)
	if err != nil {
		fmt.Println(err, getProductsUrl)
		return nil
	}

	var response responseResult

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return validatePairs(response.Data)
}

func validatePairs(rowPairs []pairInfo) (result []pairData) {

	for _, pair := range rowPairs {

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

	return
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
