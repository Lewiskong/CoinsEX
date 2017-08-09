package okcoin

import (
	"crypto/md5"
	"encoding/hex"
	"sort"
	"strings"
)

func getSign(params []string) string {
	sort.Strings(params)
	reqParam := strings.Join(params, "&")
	reqParam += "&secret_key=" + okconfig.SecretKey

	h := md5.New()
	h.Write([]byte(reqParam[:]))
	sign := hex.EncodeToString(h.Sum(nil))

	return strings.ToUpper(sign)
}

func getRequestStr(params []string) string {
	sign := getSign(params)
	params = append(params, "sign="+sign)
	return strings.Join(params, "&")
}
