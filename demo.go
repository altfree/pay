package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"pay"
)

//微信支付异步通知

var wechat pay.WechatPayConfig
var alipay pay.AlipayConfig

func init() {

	wechat.WechatAppId = "wx9257b6ca5978ec62"
	wechat.WechatMchId = "1460719102"
	wechat.WechatKey = "mfeknfk1ok2wsaso9jkoamda30kdDK22"
	wechat.WechatNotifyUrl = "https://www.baidu.com"
	alipay.Gateway = "https://openapi.alipay.com/gateway.do?"
	alipay.AppId = "2018021102179240"
	alipay.Format = "JSON"
	alipay.Charset = "UTF-8"
	alipay.Version = "1.0"
	alipay.SignType = "RSA2"
	alipay.PublicKey = "/Users/alt/go/src/golang/example/publickey.pem"
	alipay.PrivateKey = "/Users/alt/go/src/golang/example/privatekey.pem"

}

func Wenotify(w http.ResponseWriter, r *http.Request) {

	body := r.Body
	defer body.Close()
	//读取请求主体参数到data
	data, _ := ioutil.ReadAll(body)
	_, err := wechat.WechatNotify(string(data))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("签名验证成功")

	io.WriteString(w, "hello golang")
}

func main() {

	//支付宝获取支付链接
	m := map[string]string{"subject": "testlow", "out_trade_no": "yn0923j23iwwn23jk", "total_amount": "0.01"}
	x, err := alipay.AlipayGetPayUrl(m)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(x)
	//支付宝退款
	refundData := map[string]string{"out_trade_no": "yn0923j23in23jk", "refund_amount": "0.01"}
	res, err := alipay.AlipayTradeRefund(refundData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

}
