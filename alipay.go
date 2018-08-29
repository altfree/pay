package pay

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	Gateway      = "https://openapi.alipay.com/gateway.do?"
	AppId        = "商户id"
	Format       = "JSON"
	ReturnUrl    = "https://api.ncwc.com.cn"
	NotifyUrl    = "https://api.ncwc.com.cn"
	Charset      = "UTF-8"
	Version      = "1.0"
	SignType     = "RSA2"
	PublicKey    = "//Users/alt/go/src/golang/example/publickey.pem" //支付宝公钥 用于回调签名验证
	PrivateKey   = "/Users/alt/go/src/golang/example/privatekey.pem" //商户私钥用于创建支付宝订单
	ProductCode  = "FAST_INSTANT_TRADE_PAY"                          //签约产品
	PcPayMethod  = "alipay.trade.page.pay"                           //pc支付接口名称
	WapPayMethod = "alipay.trade.wap.pay"                            //pc支付接口名称

)

func GetPayUrl(param map[string]string) (string, error) {

	if param["out_trade_no"] == "" {

		return "", EmptyError("订单号不可以为空")

	}
	if len(param["out_trade_no"]) > 64 {

		return "", EmptyError("订单号长度不能超过64字符")

	}
	if param["total_amount"] == "" {

		return "", EmptyError("请输入订单金额")

	}
	if param["subject"] == "" {

		return "", EmptyError("订单标题不能为空")

	}
	param["product_code"] = ProductCode
	url, err := commonQes(WapPayMethod, param)
	if err != nil {

		return "", err

	}
	return url, nil

}

//组装请求参数
func commonQes(method string, param map[string]string) (string, error) {

	data := make(map[string]string)
	//请求业务参数排序
	content, _ := json.Marshal(param)
	data["biz_content"] = string(content)
	data["app_id"] = AppId
	data["method"] = method
	data["format"] = Format
	data["return_url"] = ReturnUrl
	data["charset"] = Charset
	data["version"] = Version
	data["notify_url"] = NotifyUrl
	data["sign_type"] = SignType
	data["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	strJoin := ascii(data)
	signStr, err := alipaySignature([]byte(strJoin), SignType)
	if err != nil {
		return "", err
	}
	data["sign"] = signStr
	requestUrl := url.Values{}
	for k, v := range data {
		requestUrl.Add(k, v)
	}

	return Gateway + requestUrl.Encode(), nil

}

//ASCII排序 并链接成字符串
func ascii(data map[string]string) string {
	var key []string
	var content string
	for k, _ := range data {
		key = append(key, k)
	}
	sort.Strings(key)

	for _, v := range key {

		if data[v] == "" {
			continue
		}
		content += v + "=" + data[v] + "&"

	}
	content = strings.TrimRight(content, "&")
	return content
}

//签名并base64编码
func alipaySignature(content []byte, sign string) (string, error) {
	key, err := ioutil.ReadFile(PrivateKey)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(key)
	if block == nil {
		return "", EmptyError("pem解码后的私钥错误")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	if SignType == "RSA2" {
		h := sha256.New()
		h.Write([]byte(content))
		digest := h.Sum(nil)
		s, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA256, digest)
		if err != nil {
			return "", err
		}
		data := base64.StdEncoding.EncodeToString(s)
		return data, nil
	}
	//默认采用RSA1签名方式
	h := sha1.New()
	h.Write([]byte(content))
	digest := h.Sum(nil)
	s, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA1, digest)
	if err != nil {
		return "", err
	}
	data := base64.StdEncoding.EncodeToString(s)
	return data, nil

}

//签名验证
func CheckNotify(notify map[string]string) (bool, error) {

	sign := notify["sign"]
	signType := notify["sign_type"]

	delete(notify, "sign") //销毁签名字符串
	delete(notify, "sign_type")

	sort := ascii(notify)
	publicKey, err := ioutil.ReadFile(PublicKey)
	if err != nil {
		return false, err
	}

	deKey, _ := pem.Decode(publicKey)
	parseKey, err := x509.ParsePKIXPublicKey(deKey.Bytes)
	rsaPub, _ := parseKey.(*rsa.PublicKey)

	//对返回的签名进行base64解码
	deSign, _ := base64.StdEncoding.DecodeString(sign)

	if signType == "RSA2" {

		//进行SHA256哈希
		digest := sha256.New()
		digest.Write([]byte(sort))
		hashData := digest.Sum(nil)
		err = rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, hashData, deSign)

	} else {
		//进行SHA1哈希
		digest := sha1.New()
		digest.Write([]byte(sort))
		hashData := digest.Sum(nil)
		err = rsa.VerifyPKCS1v15(rsaPub, crypto.SHA1, hashData, deSign)
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

//定义返回错误信息
func EmptyError(msg string) error {

	return errors.New(msg)

}
