package okcoin

import (
	"errors"
	"reflect"
	"strings"
	// "encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"time"
)

type CoinInfoItem struct {
	Date string `json:"date"`
	Buy  string `json:"ticker.buy"`
	Sell string `json:"ticker.sell"`
	High string `json:"ticker.high"`
	Low  string `json:"ticker.low"`
	Last string `json:"ticker.last"`
	Vol  string `json:"ticker.vol"`
}

type CoinInfo struct {
	size     int64
	interval time.Duration
	data     []CoinInfoItem
}

const (
	BTC_TICKER = "https://www.okcoin.cn/api/v1/ticker.do?symbol=btc_cny"
	LTC_TICKER = ""
	ETH_TICKER = ""
)

func (c *CoinInfo) init(size int64, interval time.Duration) {
	c.size = size
	c.interval = interval
}

//https://www.okcoin.cn/api/v1/ticker.do?symbol=btc_cny
func (c *CoinInfo) start(name string) {
	url := fmt.Sprintf("http://www.okcoin.cn/api/v1/ticker.do?symbol=%s_cny", name)

	for {
		rsp, _ := http.Get(url)
		fmt.Println(url)
		value, _ := ioutil.ReadAll(rsp.Body)

		item := &CoinInfoItem{}
		_unmarshal(value, item)

		_push(c, item)
		<-time.After(c.interval)

	}

}

func _push(info *CoinInfo, item *CoinInfoItem) {
	if int64(len(info.data)) >= info.size {
		// shift the oldest data
		info.data = info.data[1:]
	}
	info.data = append(info.data, *item)
}

// to optimize
// unmarshal from json bytes to CoinInfo
func _unmarshal(data []byte, v interface{}) error {

	value := reflect.ValueOf(v)
	value = reflect.Indirect(value)
	tp := reflect.TypeOf(v).Elem()
	for i := 0; i < value.NumField(); i++ {
		f := value.Field(i)
		sf := tp.Field(i)
		tag := string(sf.Tag[5:])
		tag = strings.Trim(tag, "\"")
		ret := gjson.Get(string(data[:]), tag)

		switch f.Kind() {
		case reflect.String:
			f.SetString(ret.String())
		case reflect.Float64:
			f.SetFloat(ret.Float())
		case reflect.Bool:
			f.SetBool(ret.Bool())
		default:
			return errors.New("wrong type convert")
		}

	}

	return nil
}
