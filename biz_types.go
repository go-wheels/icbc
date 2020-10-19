package icbc

type RequestBiz interface {
	ServicePath() string
}

type QrcodeQueryResponseV2Biz struct {
	ReturnCode   int    `json:"-"`
	ReturnMsg    string `json:"return_msg"`
	MsgID        string `json:"msg_id"`
	PayStatus    string `json:"pay_status"`
	CustID       string `json:"cust_id"`
	CardNo       string `json:"card_no"`
	TotalAmt     string `json:"total_amt"`
	PointAmt     string `json:"point_amt"`
	EcouponAmt   string `json:"ecoupon_amt"`
	MerDiscAmt   string `json:"mer_disc_amt"`
	CouponAmt    string `json:"coupon_amt"`
	BankDiscAmt  string `json:"bank_disc_amt"`
	PaymentAmt   string `json:"payment_amt"`
	OutTradeNo   string `json:"out_trade_no"`
	OrderID      string `json:"order_id"`
	PayTime      string `json:"pay_time"`
	TotalDiscAmt string `json:"total_disc_amt"`
}

type QrcodeQueryRequestV2Biz struct {
	MerID      string `json:"mer_id"`
	CustID     string `json:"cust_id"`
	OutTradeNo string `json:"out_trade_no"`
	OrderID    string `json:"order_id"`
}

func (QrcodeQueryRequestV2Biz) ServicePath() string {
	return "/api/qrcode/V2/query"
}

type QrcodeGenerateResponseV2Biz struct {
	ReturnCode int    `json:"-"`
	ReturnMsg  string `json:"return_msg"`
	MsgID      string `json:"msg_id"`
	Qrcode     string `json:"qrcode"`
	Attach     string `json:"attach"`
}

type QrcodeGenerateRequestV2Biz struct {
	MerID           string `json:"mer_id"`
	StoreCode       string `json:"store_code"`
	OutTradeNo      string `json:"out_trade_no"`
	OrderAmt        string `json:"order_amt"`
	TradeDate       string `json:"trade_date"`
	TradeTime       string `json:"trade_time"`
	Attach          string `json:"attach,omitempty"`
	PayExpire       string `json:"pay_expire"`
	NotifyURL       string `json:"notify_url,omitempty"`
	TporderCreateIP string `json:"tporder_create_ip"`
	SpFlag          string `json:"sp_flag,omitempty"`
	NotifyFlag      string `json:"notify_flag"`
}

func (QrcodeGenerateRequestV2Biz) ServicePath() string {
	return "/api/qrcode/V2/generate"
}
