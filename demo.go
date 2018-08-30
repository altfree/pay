package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

//微信支付异步通知

func Wenotify(w http.ResponseWriter, r *http.Request) {

	for k, v := range r.Form {
		fmt.Printf("%s=%s\n", k, v)
	}
	io.WriteString(w, "hello world")

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
	// wechatData := map[string]string{"nonce_str": "3223321321mdwojwed2", "out_trade_no": "2211212xxs", "total_fee": "1", "spbill_create_ip": "192.168.1.1", "body": "测试数据"}
	// str, err := pay.ScanPay(wechatData)
	//微信退款
	// refund := map[string]string{"nonce_str": "3223321321mdwojwed2", "out_trade_no": "2211212xxs", "total_fee": "1", "out_refund_no": "0932090932", "refund_fee": "1"}
	// str, err := pay.ApplyRefund(refund)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(str)
	http.HandleFunc("/pay/test", Wenotify)
	log.Fatal(http.ListenAndServe(":80", nil))

}
