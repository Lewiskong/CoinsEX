package main

import (
	".."
	"fmt"
)

func main() {
	okcoin.Init()
	testSell()
}

func testBuy() {
	err := okcoin.Buy(okcoin.BtcCoin, 1000, 1)
	if err != nil {
		fmt.Println(err)
	}
}

func testBuyMarket() {
	err := okcoin.BuyMarket(okcoin.BtcCoin, 100000)
	if err != nil {
		fmt.Println(err)
	}
}

func testSell() {
	err := okcoin.Sell(okcoin.BtcCoin, 22222, 10)
	if err != nil {
		fmt.Println(err)
	}
}

func testSellMarket() {
	err := okcoin.SellMarket(okcoin.BtcCoin, 20)
	if err != nil {
		fmt.Println(err)
	}
}
