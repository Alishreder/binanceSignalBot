package crypto

import (
	"fmt"
	"time"
)

type pairData struct {
	Symbol string
}

type PriceSender struct {
	PriceChanges chan string
	pairs        []pairData
}

func NewPriceSender() PriceSender {
	pairs := getPairs()
	return PriceSender{
		PriceChanges: make(chan string, len(pairs)),
		pairs:        pairs,
	}
}

func (p *PriceSender) TrackPriceChange() {

	for {

		fmt.Println("start checking price percentage change")

		for _, pair := range p.pairs {
			priceChangePercentage, err := getTokenPriceChange(pair.Symbol, 1)
			if err != nil {
				continue
			}
			if priceChangePercentage > 5 {
				p.PriceChanges <- fmt.Sprintf("AAAAALLLLLEEEEERRRRRTTTTT %s change price on %f per hour\n", pair.Symbol, priceChangePercentage)
			}

		}

		fmt.Println("that's all, time to sleep ZZZ")

		time.Sleep(time.Hour)

	}
}
