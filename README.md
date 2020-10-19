# icbc

[![GoDoc](https://godoc.org/github.com/go-wheels/icbc?status.svg)](https://godoc.org/github.com/go-wheels/icbc)
[![License](https://img.shields.io/github/license/go-wheels/icbc)](LICENSE)

工行开放平台 SDK for Go

## 使用指南

安装工行开放平台 SDK

```shell script
go get -u github.com/go-wheels/icbc
```

初始化工行开放平台 SDK

```go
options := icbc.Options{
    AppID:            "应用ID", 
    AppPrivateKey:    "应用私钥",
    GatewayPublicKey: "网关公钥",
}
client, _ := icbc.NewClient(options)
```

## 二维码扫码支付 API

二维码生成

```go
msgID := "消息通讯唯一编号"

reqBiz := &icbc.QrcodeGenerateRequestV2Biz{
    MerID:           "商户号",
    StoreCode:       "e生活号",
    OutTradeNo:      "商户订单号",
    OrderAmt:        "订单总金额 (单位: 分)",
    TradeDate:       "商户订单生成日期 (格式: yyyyMMdd)",
    TradeTime:       "商户订单生成时间 (格式: HHmmss)",
    Attach:          "商户附加数据 (原样返回)",
    PayExpire:       "二维码有效期 (单位: 秒)",
    NotifyURL:       "商户接收支付成功通知消息URL (当notify_flag为1时必填)",
    TporderCreateIP: "商户订单生成机器IP",
    SpFlag:          "扫码后是否需要跳转分行 (取值: 0 或 1)",
    NotifyFlag:      "商户是否开启通知接口 (取值: 0 或 1)",
}

respBiz := &icbc.QrcodeGenerateResponseV2Biz{}

client.Execute(msgID, requestBiz, respBiz)

log.Printf("%#v", respBiz)
```

二维码查询

```go
msgID := "消息通讯唯一编号"

reqBiz := &icbc.QrcodeQueryRequestV2Biz{
    MerID:      "商户号",
    CustID:     "支付时工行返回的用户唯一标识",
    OutTradeNo: "商户订单号",
    OrderID:    "行内订单号",
} // 商户订单号或行内订单号必须其中一个不为空

respBiz := &icbc.QrcodeQueryResponseV2Biz{}

client.Execute(msgID, requestBiz, respBiz)

log.Printf("%#v", respBiz)
```

通知验签

```go
http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
    err := client.VerifyNotification(r)
    if err != nil {
        // 验签失败
        return
    }
    // 验签成功
    log.Printf("%#v", params)
})
```
