package btceapi

import (
	"encoding/json"
	"time"
)

type Rights struct {
	Info  int `json:"info"`
	Trade int `json:"trade"`
}

type FilterParams struct {
	From     int
	Count    int
	FromID   int
	EndID    int
	OrderAsc bool
	Since    time.Time
	End      time.Time
}

type UserInfo struct {
	Funds            map[string]float64 `json:"funds"`
	Rights           Rights             `json:"rights"`
	TransactionCount int64              `json:"transaction_count"`
	OpenOrders       int64              `json:"open_orders"`
	ServerTime       float64            `json:"server_time"`
}

type Transaction struct {
	Type        int     `json:"type"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Description string  `json:"desc"`
	Status      int     `json:"status"`
	Timestamp   int64   `json:"timestamp"`
}

type TransHistory map[string]Transaction

type Trade struct {
	Pair        string  `json:"pair"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Rate        float64 `json:"rate"`
	OrderID     int     `json:"order_id"`
	IsYourOrder int     `json:"is_your_order"`
	Timestamp   int64   `json:"timestamp"`
}

type TradeHistory map[string]Trade

type Order struct {
	Pair             string  `json:"pair"`
	Type             string  `json:"type"`
	StartAmount      float64 `json:"start_amount"`
	Amount           float64 `json:"amount"`
	Rate             float64 `json:"rate"`
	TimestampCreated int64   `json:"timestamp_created"`
	Status           int     `json:"status"`
}

type TradeAnswer struct {
	Received float64            `json:"received"`
	Remains  float64            `json:"remains"`
	OrderID  int                `json:"order_id"`
	Funds    map[string]float64 `json:"funds"`
}

type CancelOrderAnswer struct {
	OrderID int                `json:"order_id"`
	Funds   map[string]float64 `json:"funds"`
}

type OrderList map[string]Order

type Depth struct {
	Asks []float64 `json:"asks"`
	Bids []float64 `json:"bids"`
}

type DepthV3 map[string]DepthValuesV3

type DepthValuesV3 struct {
	Asks [][]float64 `json:"asks"`
	Bids [][]float64 `json:"bids"`
}


type Ticker struct {
	Ticker TickerItem `json:"ticker"`
}

type TickerItem struct {
	Average       float64 `json:"avg"`
	Buy           float64 `json:"buy"`
	High          float64 `json:"high"`
	Last          float64 `json:"last"`
	Low           float64 `json:"low"`
	Sell          float64 `json:"sell"`
	Volume        float64 `json:"vol"`
	VolumeCurrent float64 `json:"vol_cur"`
	ServerTime    int64   `json:"server_time"`
}

type TickerV3  map[string]TickerItemV3

type TickerItemV3 struct {
	High float64 `json:"high"`
	Low float64 `json:"low"`
	Avg float64 `json:"avg"`
	Volume float64 `json:"vol"`
	VolCurrent float64 `json:"vol_cur"`
	Last float64 `json:"last"`
	Buy float64 `json:"buy"`
	Sell float64 `json:"sell"`
	Updated int64 `json:"updated"`
}


type Fee struct {
	Trade float64 `json:"trade"`
}
type FeeV3 map[string]float64


type TradeList []TradeInfo

type TradeInfo struct {
	Amount        float64 `json:"amount"`
	Price         float64 `json:"price"`
	Date          int64   `json:"date"`
	Item          string  `json:"item"`
	PriceCurrency string  `json:"price_currency"`
	Tid           int32   `json:"tid"`
	Type          string  `json:"type"`
}

type TradeListV3 map[string][]TradeInfoV3

type TradeInfoV3 struct {
	Type string `json:"type"`
	Price float64 `json:"price"`
	Amount float64 `json:"amount"`
	Tid int32 `json:"tid"`
	Timestamp int64 `json:"timestamp"`
}


type RawResponse struct {
	Success int             `json:"success"`
	Return  json.RawMessage `json:"return"`
}
