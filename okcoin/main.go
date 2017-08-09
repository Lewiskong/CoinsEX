package okcoin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// various coin infos
var (
	infos map[string]*CoinInfo
)

// buy and sell type
const (
	BtcCoin = iota
	LtcCoin
	EthCoin
)

// way of trade
const (
	TradeBuy = iota
	TradeBuyMarket
	TradeSell
	TradeSellMarket
)

type TradeWay uint
type CoinType uint

func Init() {

	initConf()
	// err := Buy(EthCoin, 1888, 1)
	// err := Buy(EthCoin, 1000, 1)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	infos = make(map[string]*CoinInfo)
	infos["btc"] = &CoinInfo{}
	infos["ltc"] = &CoinInfo{}
	infos["eth"] = &CoinInfo{}
	for k, _ := range infos {
		v := infos[k]
		v.init(3, time.Second)
		go v.start(k)
	}

	// <-time.After(time.Second * 8)

}

func GetUserInfo() {
	param := []string{fmt.Sprintf("api_key=%s", okconfig.AppKey)}
	reqStr := getRequestStr(param)
	rsp, _ := http.Post("https://www.okcoin.cn/api/v1/userinfo.do", "application/x-www-form-urlencoded", strings.NewReader(reqStr))
	bts, _ := ioutil.ReadAll(rsp.Body)
	fmt.Println(string(bts))

}

func Buy(coin CoinType, price, amount float64) error {
	err := Trade(coin, price, amount, TradeBuy)
	return err
}

func BuyMarket(coin CoinType, price float64) error {
	err := Trade(coin, price, 0, TradeBuyMarket)
	return err
}

func Sell(coin CoinType, price, amount float64) error {
	err := Trade(coin, price, amount, TradeSell)
	return err
}

func SellMarket(coin CoinType, amount float64) error {
	err := Trade(coin, -1, amount, TradeSellMarket)
	return err
}

func Trade(coin CoinType, price, amount float64, way TradeWay) error {
	symbol, err := _getCoinSymbol(coin)
	if err != nil {
		return err
	}
	tradeWay, err := _getTradeWay(way)
	if err != nil {
		return err
	}
	param := []string{
		"api_key=" + okconfig.AppKey,
		"symbol=" + symbol,
		"type=" + tradeWay,
		"price=" + strconv.FormatFloat(price, 'f', -1, 32),
		"amount=" + strconv.FormatFloat(amount, 'f', -1, 32),
	}
	reqStr := getRequestStr(param)
	rsp, err := http.Post("https://www.okcoin.cn/api/v1/trade.do", "application/x-www-form-urlencoded", strings.NewReader(reqStr))
	if err != nil {
		return err
	}
	bts, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	type result struct {
		Result     bool    `json:"result"`
		Order_id   float32 `json:"order_id"`
		Error_code float32 `json:"error_code"`
	}
	ret := result{}

	err = json.Unmarshal(bts, &ret)
	fmt.Println(string(bts[:]))
	fmt.Println(ret)
	if err != nil {
		return err
	}
	if !ret.Result {
		return errors.New(errCodeConfig[ret.Error_code])
	}
	return nil
}

func _getCoinSymbol(coin CoinType) (string, error) {
	var symbol string
	switch coin {
	case BtcCoin:
		symbol = "btc_cny"
	case LtcCoin:
		symbol = "ltc_cny"
	case EthCoin:
		symbol = "eth_cny"
	default:
		return "", errors.New("Unsupported coin type")
	}
	return symbol, nil
}

func _getTradeWay(way TradeWay) (wayName string, err error) {
	switch way {
	case TradeBuy:
		wayName = "buy"
	case TradeBuyMarket:
		wayName = "buy_market"
	case TradeSell:
		wayName = "sell"
	case TradeSellMarket:
		wayName = "sell_market"
	default:
		return "", errors.New("unsupported trade way")
	}
	return
}
