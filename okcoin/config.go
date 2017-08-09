package okcoin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type OkcoinConfig struct {
	AppKey    string `json:"api_key"`
	SecretKey string `json:"secret_key"`
}

var (
	okconfig      OkcoinConfig
	errCodeConfig = make(map[float32]string)
)

func initConf() {
	// init account config
	body, err := ioutil.ReadFile("../account.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &okconfig)
	if err != nil {
		panic("config json parse error")
	}
	fmt.Println(okconfig)

	// init err_code config
	body, err = ioutil.ReadFile("../Errors.conf")
	if err != nil {
		panic(err)
	}
	kvs := strings.Split(string(body[:]), "\n")
	for _, v := range kvs {
		kv := strings.Split(v, "\t")
		if len(kv) < 2 {
			continue
		}
		key, err := strconv.ParseFloat(kv[0], 32)
		if err != nil {
			panic(err)
		}
		errCodeConfig[float32(key)] = kv[1]
	}
	fmt.Println("init errcode complete")
}
