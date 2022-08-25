package crypto

import (
	"fmt"
	"time"
)

const templatePriceChange = "ALERT!\n%s change price on %.2f%% per %s\nSee chart: %s"

type pairData struct {
	Symbol string
}

type PriceSender struct {
	PriceChanges chan string
	pairs        []pairData
}

func NewPriceSender() PriceSender {
	pairs := getPairs()
	fmt.Println(pairs)
	return PriceSender{
		PriceChanges: make(chan string, len(pairs)),
		pairs:        pairs,
	}
}

func (p *PriceSender) TrackPriceChange(windowSize string, timeToSleep time.Duration) {

	for {

		fmt.Println("start checking price percentage change for", windowSize)

		for _, pair := range p.pairs {
			priceChangePercentage, err := getTokenPriceChange(pair.Symbol, windowSize)
			if err != nil {
				continue
			}
			if priceChangePercentage > 4 {
				message := fmt.Sprintf(templatePriceChange, pair.Symbol, priceChangePercentage, windowSize, generateTradingViewURL(pair.Symbol))
				fmt.Println(message)
				p.PriceChanges <- message
			}

		}

		fmt.Printf("that's all for %s, time to sleep ZZZ\n", windowSize)

		time.Sleep(timeToSleep)

	}
}

func generateTradingViewURL(symbol string) string {
	return "https://www.tradingview.com/chart/?symbol=BINANCE%3A" + symbol
}
