package pay

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"strconv"
	"strings"
)

type Wechat interface {
	JsPay() (string, error)
	AppPay() (string, error)
	ScanPay() (string, error)
}

const (
	WechatAppId         = "公众号appid"
	WechatMchId         = "商户号"
	WechatNotifyUrl     = "http://example.com/payments/wechat-notify"
	WechatCreatTradeUrl = "https://api.mch.weixin.qq.com/pay/unifiedorder"  //统一下单
	WechatRefund        = "https://api.mch.weixin.qq.com/secapi/pay/refund" //申请退款
	WechatQueryRefund   = "https://api.mch.weixin.qq.com/pay/refundquery"   //退款查询
	WechatAppPay        = "APP"
	WechatJsPay         = "JSAPI"
	WechatScan          = "NATIVE"
	WechatSignType      = "HMAC-SHA256" //签名方式默认为md5
	WechatKey           = "mfeknfk1ok2wsaso9jkoamda30kdDK22"
)

//定义支付数据解析xml的结构体
type ReuqestPayParam struct {
	ReturnCode string `xml:"return_code" json:"return_code"`
	ReturnMsg  string `xml:"return_msg"  json:"return_msg"`
	ResultCode string `xml:"result_code"  json:"result_code"`
	ErrCode    string `xml:"err_code"  json:"err_code"`
	ErrCodeDes string `xml:"err_code_des"  json:"err_code_des"`
	TradeType  string `xml:"trade_type"  json:"trade_type"`
	PrePayId   string `xml:"prepay_id"  json:"prepay_id"`
	CodeUrl    string `xml:"code_url"  json:"code_url"`
}

//定义退款查询数据解析xml的结构体
type RefundParam struct {
	ReturnCode    string `xml:"return_code" json:"return_code"`
	ReturnMsg     string `xml:"return_msg"  json:"return_msg"`
	ResultCode    string `xml:"result_code"  json:"result_code"`
	ErrCode       string `xml:"err_code"  json:"err_code"`
	ErrCodeDes    string `xml:"err_code_des"  json:"err_code_des"`
	Totalfee      string `xml:"total_fee"  json:"total_fee"`
	RefundCount   string `xml:"refund_count"  json:"refund_count"`
	CashFee       string `xml:"cash_fee"  json:"cash_fee"`
	FeeType       string `xml:"fee_type"  json:"fee_type"`
	OutRefundNo   string `xml:"out_refund_no_$n"  json:"out_refund_no_$n"`
	RufundId      string `xml:"refund_id_$n"  json:"refund_id_$n"`
	RedundChannel string `xml:"refund_channel_$n"  json:"refund_channel_$n"`
	RefundFee     string `xml:"refund_fee_$n"  json:"refund_fee_$n"`
	RefundStatus  string `xml:"refund_status_$n"  json:"refund_status_$n"`
	RefundRecv    string `xml:"refund_recv_accout_$n"  json:"refund_recv_accout_$n"`
}

//定义支付数据解析xml的结构体
type QueryRefundParam struct {
	ReturnCode  string `xml:"return_code" json:"return_code"`
	ReturnMsg   string `xml:"return_msg"  json:"return_msg"`
	ResultCode  string `xml:"result_code"  json:"result_code"`
	ErrCode     string `xml:"err_code"  json:"err_code"`
	ErrCodeDes  string `xml:"err_code_des"  json:"err_code_des"`         //错误信息描述
	RufundId    string `xml:"refund_id"  json:"refund_id"`               //微信退款单号
	RefundFee   string `xml:"refund_fee"  json:"refund_fee"`             //退款金额
	Totalfee    string `xml:"total_fee"  json:"total_fee"`               //订单总金额
	OutRefundNo string `xml:"out_refund_no_$n"  json:"out_refund_no_$n"` //商户退款订单号
}

var publicParameter map[string]string

//准备发起的交易信息 请按照微信开发文档的字段进行传值 https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1
func creatTrade(genre string, param map[string]string) (string, error) {

	param["appid"] = WechatAppId
	param["mch_id"] = WechatMchId
	// param["sign_type"] = WechatSignType
	// param["body"] = WechatAppId
	// param["out_trade_no"] = WechatAppId
	// param["total_fee"] = WechatAppId
	// param["spbill_create_ip"] = WechatAppId
	param["notify_url"] = WechatNotifyUrl
	param["trade_type"] = genre
	if param["appid"] == "" {

		return "", EmptyError("请传入公众号appid")

	}
	if param["mch_id"] == "" {

		return "", EmptyError("请传入商户号mch_id")

	}
	if param["nonce_str"] == "" {

		return "", EmptyError("请传入随机字符串")

	}
	if param["out_trade_no"] == "" {

		return "", EmptyError("请传入商户订单号")

	}
	amount, _ := strconv.Atoi(param["total_fee"])
	if amount < 1 {

		return "", EmptyError("订单金额最低为1分钱")

	}
	if param["spbill_create_ip"] == "" {

		return "", EmptyError("请传入用户地址ip")

	}
	signStr := signature(param)
	param["sign"] = signStr //转换为大写
	res, err := paramFormat(WechatCreatTradeUrl, param)
	if err != nil {

		return "", err

	}
	//处理返回参数的格式
	var data ReuqestPayParam
	xml.Unmarshal(res, &data)
	res, _ = json.Marshal(data)
	return string(res), nil

}

//签名   默认使用MD5
func signature(param map[string]string) string {

	sort := ascii(param)
	waitSign := sort + "&key=" + WechatKey
	signMethod := md5.New()
	signMethod.Write([]byte(waitSign))
	signStr := signMethod.Sum(nil)
	return strings.ToUpper(hex.EncodeToString(signStr))
}

//格式化请求参数
func paramFormat(url string, param map[string]string) ([]byte, error) {
	xmlData := "<xml>"
	for k, v := range param {

		xmlData += "<" + k + ">" + v + "</" + k + ">"

	}
	xmlData += "</xml>"
	res, err := CurlGetRes(WechatCreatTradeUrl, xmlData)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//申请退款
func ApplyRefund(refund map[string]string) (string, error) {

	refund["appid"] = WechatAppId
	refund["mch_id"] = WechatMchId
	if refund["appid"] == "" {

		return "", EmptyError("请传入公众号appid")

	}
	if refund["mch_id"] == "" {

		return "", EmptyError("请传入商户号mch_id")

	}
	if refund["nonce_str"] == "" {

		return "", EmptyError("请传入随机字符串")

	}
	if refund["out_refund_no"] == "" {

		return "", EmptyError("退款订单号不可以为空")

	}
	amount, _ := strconv.Atoi(refund["refund_fee"])
	if amount < 1 {

		return "", EmptyError("退款金额最低为1分钱")

	}
	if refund["out_trade_no"] == "" && refund["transaction_id"] == "" {

		return "", EmptyError("微信订单号和商户订单号不能同时为空")

	}
	signStr := signature(refund)
	refund["sign"] = signStr
	res, err := paramFormat(WechatRefund, refund)
	if err != nil {

		return "", err

	}
	//处理返回参数的格式
	var data RefundParam
	xml.Unmarshal(res, &data)
	res, _ = json.Marshal(data)
	return string(res), nil

}

func JsPay(param map[string]string) (string, error) {

	if param["openid"] == "" {

		return "", EmptyError("JSAPI支付必须传入用户openid")
	}

	payParam, err := creatTrade(WechatJsPay, param)

	if err != nil {
		return "", err
	}
	return payParam, nil
}
func AppPay(param map[string]string) (string, error) {

	payParam, err := creatTrade(WechatAppPay, param)
	if err != nil {
		return "", err
	}
	return payParam, nil
}

func ScanPay(param map[string]string) (string, error) {

	payParam, err := creatTrade(WechatScan, param)
	if err != nil {
		return "", err
	}
	return payParam, nil
}

//通知支付结果
func WechatNotify(param map[string]string) (bool, error) {

	signValue := param["sign"]
	signType := param["sign_type"]
	delete(param, "sign")
	delete(param, "sign_type")
	var signStr string
	if signType == "HMAC-SHA256" {

		return false, EmptyError("这里进行sha256签名")

	} else {

		signStr = signature(param)
	}
	if signValue == signStr {

		return true, nil

	}
	return false, EmptyError("签名不一致")

}