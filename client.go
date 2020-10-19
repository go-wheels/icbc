package icbc

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

const (
	gatewayBaseURL  = "https://gw.open.icbc.com.cn"
	signTypeRSA2    = "RSA2"
	charsetUTF8     = "UTF-8"
	formatJSON      = "json"
	timestampLayout = "2006-01-02 15:04:05"
)

var (
	DefaultLocation, _ = time.LoadLocation("Asia/Shanghai")
)

type commonResponse struct {
	ResponseBizContent json.RawMessage `json:"response_biz_content"`
	Sign               string          `json:"sign"`
}

type Client struct {
	httpClient       *http.Client
	appID            string
	appPrivateKey    *rsa.PrivateKey
	gatewayPublicKey *rsa.PublicKey
}

type Options struct {
	Timeout          time.Duration
	AppID            string
	AppPrivateKey    string
	GatewayPublicKey string
}

func NewClient(options Options) (client *Client, err error) {
	appPrivateKey, err := parseRSAPrivateKey(options.AppPrivateKey)
	if err != nil {
		return
	}
	gatewayPublicKey, err := parseRSAPublicKey(options.GatewayPublicKey)
	if err != nil {
		return
	}
	client = &Client{
		httpClient:       &http.Client{Timeout: options.Timeout},
		appID:            options.AppID,
		appPrivateKey:    appPrivateKey,
		gatewayPublicKey: gatewayPublicKey,
	}
	return
}

func (c *Client) VerifyNotification(req *http.Request) (err error) {
	err = req.ParseForm()
	if err != nil {
		return
	}
	stringToSign := c.buildStringToSign(req.URL.Path, req.Form)
	sign := req.Form.Get("sign")
	err = c.verify(stringToSign, sign)
	return
}

func (c *Client) Execute(msgID string, reqBiz RequestBiz, respBiz interface{}) (err error) {
	params := make(url.Values)

	params.Set("app_id", c.appID)
	params.Set("msg_id", msgID)
	params.Set("format", formatJSON)
	params.Set("charset", charsetUTF8)
	params.Set("sign_type", signTypeRSA2)

	timestamp := time.Now().In(DefaultLocation).Format(timestampLayout)
	params.Set("timestamp", timestamp)

	bizContent, err := json.Marshal(reqBiz)
	if err != nil {
		return
	}
	params.Set("biz_content", string(bizContent))

	stringToSign := c.buildStringToSign(reqBiz.ServicePath(), params)
	sign, err := c.sign(stringToSign)
	if err != nil {
		return
	}
	params.Set("sign", sign)

	serviceURL := c.buildURL(reqBiz.ServicePath())
	resp, err := c.httpClient.PostForm(serviceURL, params)
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var commResp commonResponse
	err = json.NewDecoder(resp.Body).Decode(&commResp)
	if err != nil {
		return
	}
	err = c.verify(string(commResp.ResponseBizContent), commResp.Sign)
	if err != nil {
		return
	}
	if respBiz != nil {
		err = json.Unmarshal(commResp.ResponseBizContent, respBiz)
		if err != nil {
			return
		}
		fld := reflect.ValueOf(respBiz).Elem().FieldByName("ReturnCode")
		if fld.CanSet() && fld.Kind() == reflect.Int {
			fld.SetInt(gjson.GetBytes(commResp.ResponseBizContent, "return_code").Int())
		}
	}
	return
}

func (c *Client) verify(data, sign string) (err error) {
	sig, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return
	}
	sum := sha1.Sum([]byte(data))
	err = rsa.VerifyPKCS1v15(c.gatewayPublicKey, crypto.SHA1, sum[:], sig)
	return
}

func (c *Client) sign(data string) (sign string, err error) {
	sum := sha256.Sum256([]byte(data))
	sig, err := rsa.SignPKCS1v15(rand.Reader, c.appPrivateKey, crypto.SHA256, sum[:])
	if err != nil {
		return
	}
	sign = base64.StdEncoding.EncodeToString(sig)
	return
}

func (Client) buildStringToSign(path string, params url.Values) string {
	keys := make(sort.StringSlice, 0, len(params))
	for key := range params {
		if key != "" && key != "sign" {
			keys = append(keys, key)
		}
	}
	keys.Sort()

	var buf strings.Builder
	buf.WriteString(path)
	buf.WriteByte('?')
	for i, key := range keys {
		if i > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(key)
		buf.WriteByte('=')
		buf.WriteString(params.Get(key))
	}
	return buf.String()
}

func (Client) buildURL(path string) string {
	return gatewayBaseURL + path
}
