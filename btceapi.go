// Package btceapi provides the API for the btc-e.com stock exchange.
package btceapi

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var ApiURL string = "https://btc-e.com"

type BtceAPI struct {
	Key    string
	Secret string
}

func (btcApi BtceAPI) GetInfo() (UserInfo, error) {
	params := make(map[string]string)
	params["method"] = "getInfo"

	res := UserInfo{}
	err := query(btcApi, "getInfo", params, &res)

	return res, err
}

func (btcApi BtceAPI) GetTransHistory(filterParams FilterParams) (TransHistory, error) {
	params := getParams(filterParams)
	res := TransHistory{}
	err := query(btcApi, "TransHistory", params, &res)

	return res, err
}

func (btcApi BtceAPI) GetTradeHistory(filterParams FilterParams) (TradeHistory, error) {
	params := getParams(filterParams)
	res := TradeHistory{}
	err := query(btcApi, "TradeHistory", params, &res)

	return res, err
}

func (btcApi BtceAPI) GetOrderList(filterParams FilterParams) (OrderList, error) {
	params := getParams(filterParams)
	res := OrderList{}
	err := query(btcApi, "OrderList", params, &res)

	return res, err
}

func (btcApi BtceAPI) Trade(pair string, tradeType string, rate float64, amount float64) (TradeAnswer, error) {
	params := make(map[string]string)

	params["pair"] = pair
	params["type"] = tradeType
	params["rate"] = strconv.FormatFloat(rate, 'f', -1, 64)
	params["amount"] = strconv.FormatFloat(amount, 'f', -1, 64)

	res := TradeAnswer{}
	err := query(btcApi, "Trade", params, &res)
	return res, err
}

func (btcApi BtceAPI) CancelOrder(orderId string) (CancelOrderAnswer, error) {
	params := make(map[string]string)

	params["order_id"] = orderId

	res := CancelOrderAnswer{}
	err := query(btcApi, "CancelOrder", params, &res)
	return res, err
}

func GetDepthV2(pair string) (Depth, error) {
	url := fmt.Sprintf("api/2/%s/depth", pair)
	res := Depth{}
	err := makeGetCall(url, &res)
	return res, err
}


func GetDepthV3(pair string) (Depth, error) {

	url := fmt.Sprintf("api/3/depth/%s", pair)
	res := DepthV3{}
	err := makeGetCall(url, &res)
	return res, err
}

func GetTickerV2(pair string) (Ticker, error) {
	url := fmt.Sprintf("api/2/%s/ticker", pair)
	res := Ticker{}
	err := makeGetCall(url, &res)
	return res, err
}

func GetTickerV3(pair string) (TickerV3, error) {
	url := fmt.Sprintf("api/3/ticker/%s", pair)
	res := TickerV3{}
	err := makeGetCall(url, &res)
	return res, err
}

func GetTradesV2(pair string) (TradeList, error) {
	url := fmt.Sprintf("api/2/%s/trades", pair)
	res := TradeList{}
	err := makeGetCall(url, &res)
	return res, err
}

func GetTradesV3(pair string) (TradeList, error) {
	url := fmt.Sprintf("api/3/trades/%s", pair)
	res := TradeListV3{}
	err := makeGetCall(url, &res)
	return res, err
}

func GetFeeV2(pair string) (Fee, error) {
	url := fmt.Sprintf("api/2/%s/fee", pair)
	res := Fee{}
	err := makeGetCall(url, &res)
	return res, err
}

func GetFeeV3(pair string) (Fee, error) {
	url := fmt.Sprintf("api/3/fee/%s", pair)
	res := FeeV3{}
	err := makeGetCall(url, &res)
	return res, err
}

func getParams(filterParams FilterParams) map[string]string {
	params := make(map[string]string)
	defaultTime := time.Time{}

	if filterParams.From > 0 {
		params["from"] = string(filterParams.From)
	}
	if filterParams.Count > 0 {
		params["count"] = string(filterParams.Count)
	}
	if filterParams.FromID > 0 {
		params["from_id"] = string(filterParams.FromID)
	}
	if filterParams.EndID > 0 {
		params["end_id"] = string(filterParams.EndID)
	}
	params["order"] = orderName(filterParams.OrderAsc)

	if filterParams.Since != defaultTime {
		params["since"] = string(filterParams.Since.Unix())
	}
	if filterParams.End != defaultTime {
		params["end"] = string(filterParams.End.Unix())
	}

	return params
}

func orderName(orderAsc bool) string {
	if orderAsc {
		return "ASC"
	}

	return "DESC"
}

func query(btcAPI BtceAPI, method string, params map[string]string, result interface{}) error {
	params["method"] = method
	params["nonce"] = strconv.FormatInt(time.Now().Unix(), 10)
	client := http.Client{Timeout: time.Duration(15 * time.Second)}

	data := url.Values{}

	for key := range params {
		data.Add(key, params[key])
	}
	dataEncode := data.Encode()

	bfs := bytes.NewBufferString(dataEncode)
	req, _ := http.NewRequest("POST", getFullURL("tapi"), bfs)
	sign := signData([]byte(btcAPI.Secret), bfs.Bytes())

	req.Header.Set("Key", btcAPI.Key)
	req.Header.Set("Sign", sign)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	respData := RawResponse{}

	err = json.Unmarshal(bodyBytes, &respData)
	if err != nil {
		return err
	}

	if respData.Success == 1 {
		err = json.Unmarshal(respData.Return, &result)
		if err != nil {
			return err
		}
	} else {
		return nil
	}

	return nil
}

func getFullURL(url string) string {
	return fmt.Sprintf("%s/%s", ApiURL, url)
}

func makeGetCall(url string, result interface{}) error {
	fullURL := getFullURL(url)
	resp, err := http.Get(fullURL)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	return nil
}

func signData(key []byte, data []byte) string {
	hashMaker := hmac.New(sha512.New, key)
	hashMaker.Write(data)
	return strings.ToLower(hex.EncodeToString(hashMaker.Sum(nil)))
}
