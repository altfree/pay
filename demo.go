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

func Wenotify(w http.ResponseWriter, r *http.Request) {

	body := r.Body
	defer body.Close()
	//读取请求主体参数到data
	data, _ := ioutil.ReadAll(body)
	_, err := pay.WechatNotify(string(data))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("签名验证成功")

	io.WriteString(w, "hello golang")
}

func main() {

	// //支付宝获取支付链接
	// m := map[string]string{"subject": "testlow", "out_trade_no": "ewewweewweweew999", "total_amount": "0.01"}
	// x, err := pay.AlipayGetPayUrl(m)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(x)
	// //支付宝退款
	// refundData := map[string]string{"out_trade_no": "ewewweewweweew999", "refund_amount": "0.01"}
	// res, err := pay.AlipayTradeRefund(refundData)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(res)
	//微信支付
	wechatData := map[string]string{"nonce_str": "3223321321wwmdwojwed2", "out_trade_no": "232323211212xxs", "total_fee": "1", "spbill_create_ip": "192.168.1.1", "body": "测试数据"}
	str, err := pay.ScanPay(wechatData)
	//微信退款 待完成
	// refund := map[string]string{"nonce_str": "3223321321mdwojwed2", "out_trade_no": "2211212xxs", "total_fee": "1", "out_refund_no": "0932090932", "refund_fee": "1"}
	// str, err := pay.ApplyRefund(refund)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(str)
	// http.HandleFunc("/pay/test", Wenotify)
	// log.Fatal(http.ListenAndServe(":9091", nil))
	// x, err := pay.WechatNotify(str)
	// if err != nil {

	// 	fmt.Println(err)
	// }
	// fmt.Println(x)

}
