package htx_wss

import "encoding/json"

type Values struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}
type Ticker struct {
	Exchange string  `json:"exchange"`
	Symbol   string  `json:"symbol"`
	Buy      Values  `json:"buy"`
	Sell     Values  `json:"sell"`
	UpdateAt float64 `json:"update_at"`
}
type ReciveData struct {
	Exchange string `json:"exchange"`
	Symbol   string `json:"symbol"`
	Ticker   Ticker `json:"ticker"`
}

type ReciveBalanceMsg struct {
	Exchange  string  `json:"exchange"`
	Symbol    string  `json:"symbol"`
	Available float64 `json:"available"`
	AccountId int64   `json:"accountId"`
	Balance   float64 `json:"balance"`
}

type ReciveAccountsMsg struct {
	Symbol        string      `json:"symbol"`
	ContractCode  string      `json:"contract_code"`
	MarginBalance json.Number `json:"margin_balance"`
	MarginStatic  json.Number `json:"margin_static"`
}

type ReciveSwapAccountsMsg struct {
	Symbol      string  `json:"symbol"`
	FreeBalance float64 `json:"free_balance"` // 可用保金
	LockBalance float64 `json:"lock_balance"` // 冻结保金
	LiquidPrice float64 `json:"liquid_price"` // 强平价格
	MarginRatio float64 `json:"margin_ratio"` // 保证金率
	UpdateAt    float64 `json:"update_at"`    // 更新时间
}

type ReciveSwapPositionMsg struct {
	Symbol     string  `json:"symbol"`
	SellVolume int64   `json:"sell_volume"` // 持仓张数
	BuyVolume  int64   `json:"buy_volume"`  // 持仓张数
	UpdateAt   float64 `json:"update_at"`   // 更新时间
}
type ReciveSwapFundingRateMsg struct {
	Symbol      string  `json:"symbol"`
	FundingRate float64 `json:"funding_rate"` // buy or sell
	FundingTime int64   `json:"funding_time"` // 10位时间戳
	UpdateAt    float64 `json:"update_at"`    // 更新时间
}
type ReciveSpotOrderMsg struct {
	Exchange    string  `json:"exchange"`
	Symbol      string  `json:"symbol"`
	OrderId     string  `json:"order_id"`
	OrderType   string  `json:"order_type"`   // buy-market: 市价买单 sell-market: 市价卖单
	OrderPrice  float64 `json:"order_price"`  // 下单价格
	TradePrice  float64 `json:"trade_price"`  // 成交价格
	OrderValue  float64 `json:"order_value"`  // 下单金额
	TradeValue  float64 `json:"trade_value"`  // 成交金额
	OrderVolume float64 `json:"order_volume"` // 下单数量
	TradeVolume float64 `json:"trade_volume"` // 成交数量
	Status      int64   `json:"status"`       // 订单状态 1-已下单 2-已成交 8-已撤单
	CreateAt    int64   `json:"create_at"`    // 创建时间
	FilledAt    int64   `json:"filled_at"`    // 成交时间
	CancelAt    int64   `json:"cancel_at"`    // 撤单时间
}
type ReciveSwapOrderMsg struct {
	Exchange    string  `json:"exchange"`
	Symbol      string  `json:"symbol"`
	OrderId     string  `json:"order_id"`
	OrderType   string  `json:"order_type"`   // buy-market: 市价买单 sell-market: 市价卖单
	OrderPrice  float64 `json:"order_price"`  // 下单价格
	TradePrice  float64 `json:"trade_price"`  // 成交价格
	OrderValue  float64 `json:"order_value"`  // 下单金额
	TradeValue  float64 `json:"trade_value"`  // 成交金额
	OrderVolume float64 `json:"order_volume"` // 下单数量
	TradeVolume float64 `json:"trade_volume"` // 成交数量
	Status      int64   `json:"status"`       // 订单状态 1-已下单 2-已成交 8-已撤单
	CreateAt    int64   `json:"create_at"`    // 创建时间
	FilledAt    int64   `json:"filled_at"`    // 成交时间
	CancelAt    int64   `json:"cancel_at"`    // 撤单时间
}
