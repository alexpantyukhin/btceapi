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

var ApiUrl string = "https://btc-e.com"

type BtceAPI struct {
	key    string
	secret string
}

func (btcApi BtceAPI) GetInfo() (UserInfo, error) {
	params := make(map[string]string)
	params["method"] = "getInfo"

	res := UserInfo{}
	err := query(btcApi, params, &res)

	return res, err
}

func (btcApi BtceAPI) GetTransHistory(filterParams FilterParams) (TransHistory, error) {
	params := getParams("TradeHistory", filterParams)
	res := TransHistory{}
	err := query(btcApi, params, &res)

	return res, err
}

func (btcApi BtceAPI) GetTradeHistory(filterParams FilterParams) (TradeHistory, error) {
	params := getParams("TradeHistory", filterParams)
	res := TradeHistory{}
	err := query(btcApi, params, &res)

	return res, err
}

func (btcApi BtceAPI) GetOrderList(filterParams FilterParams) (OrderList, error) {
	params := getParams("OrderList", filterParams)
	res := OrderList{}
	err := query(btcApi, params, &res)

	return res, err
}

func (btcApi BtceAPI) Trade(pair string, tradeType string, rate float64, amount float64) (TradeAnswer, error) {
	params := make(map[string]string)

	params["method"] = "Trade"
	params["pair"] = pair
	params["type"] = tradeType
	params["rate"] = strconv.FormatFloat(rate, 'f', -1, 64)
	params["amount"] = strconv.FormatFloat(amount, 'f', -1, 64)

	res := TradeAnswer{}
	err := query(btcApi, params, &res)
	return res, err
}

func (btcApi BtceAPI) CancelOrder(orderId string) (CancelOrderAnswer, error) {
	params := make(map[string]string)

	params["method"] = "CancelOrder"
	params["order_id"] = orderId

	res := CancelOrderAnswer{}
	err := query(btcApi, params, &res)
	return res, err
}

func GetDepth(pair string) (Depth, error) {
	url := fmt.Sprintf("api/2/%s/depth", pair)
	res := Depth{}
	err := makeGetCall(url, &res)
	return res, err
}

func GetTicker(pair string) (Ticker, error) {
	url := fmt.Sprintf("api/2/%s/ticker", pair)
	res := Ticker{}
	err := makeGetCall(url, &res)
	return res, err
}

func GetTrades(pair string) (TradeList, error) {
	url := fmt.Sprintf("api/2/%s/trades", pair)
	res := TradeList{}
	err := makeGetCall(url, &res)
	return res, err
}

func GetFee(pair string) (Fee, error) {
	url := fmt.Sprintf("api/2/%s/fee", pair)
	res := Fee{}
	err := makeGetCall(url, &res)
	return res, err
}

func getParams(methodName string, filterParams FilterParams) map[string]string {
	params := make(map[string]string)
	params["method"] = methodName
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

func query(btcAPI BtceAPI, params map[string]string, result interface{}) error {
	params["nonce"] = strconv.FormatInt(time.Now().Unix(), 10)
	client := http.Client{Timeout: time.Duration(15 * time.Second)}

	data := url.Values{}

	for key := range params {
		data.Add(key, params[key])
	}
	dataEncode := data.Encode()

	bfs := bytes.NewBufferString(dataEncode)
	req, _ := http.NewRequest("POST", getFullURL("tapi"), bfs)
	sign := signData([]byte(btcAPI.secret), bfs.Bytes())

	req.Header.Set("Key", btcAPI.key)
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
	return fmt.Sprintf("%s/%s", ApiUrl, url)
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
