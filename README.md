<<<<<<< HEAD
# Pay
## 在调用相关方法时请先设置相关支付商户信息
### WechatPay


#### &nbsp; &nbsp;示例（签名方式为md5,如果需要传递可选支付参数，请阅读源代码）
    package main
    
    import (
    	"fmt"
    	"io"
    	"io/ioutil"
    	"log"
    	"net/http"
    	"github.com/Jiebodong/pay"
    	""
    	)
    //异步回调测试
    func NotifyTest(w http.Responsewriter,r *http.Request){
        
        body := r.Body
	    defer body.Close()
	    //读取请求主体参数到data
	    data, _ := ioutil.ReadAll(body)
	    _, err := pay.WechatNotify(string(data)) //传入回调参数，验证签名，签名成功则返回true，反之为false和错误信息
	    if err != nil {
	    	log.Fatal(err)
	    }
    	fmt.Println("签名验证成功")
    
    	io.WriteString(w, "微信支付回调通知")
        
    }
    	
    	
    	
    func main(){
    
        	payParam := map[string]string{"subject": "标题", "out_trade_no": "订单号", "total_amount": "金额（最小一分","openid":"用户openid（jsapi支付必传）"}
        	//backMsg 数据为json格式 
        	backMsg, err :=pay.ScanPay(payParam) //扫码支付
        	backMsg, err :=pay.ScanPay(payParam) //app支付参数
            backMsg, err :=pay.ScanPay(payParam) //jsapi支付 
            //如果想转换为map
            dataMap:=make(map[string]string)
            json.Unmarshal(res, &dataMap)  //将数据赋值给dataMap
            //支付结果通知
            pay.WechatNotify(xml数据)  
        
    }
    


### Alipay
#### &nbsp; &nbsp;示例（支持sha1/sha256签名方式）
    
    //支付宝获取支付链接
	paramData:= map[string]string{"subject": "标题", "out_trade_no": "订单号", "total_amount": "金额"}
	url, err := pay.AlipayGetPayUrl(paramData)  //返回支付url,将获取到的url调转即可到支付宝支付页面
	if err != nil {
	    fmt.Println(err)
	}
	fmt.Println(url)
	//支付宝退款
	refundData := map[string]string{"out_trade_no": "订单号", "refund_amount": "退款金额"}
    backMsg, err := pay.AlipayTradeRefund(refundData) //返回json字符串
	if err != nil {
	 	fmt.Println(err)
	}
	fmt.Println(backMsg)
	
    //支付结果，异步回调传入回调参数
    msg,err:=pay.AlipayCheckNotify("回调参数map[string]string类型")//验证签名，签名成功则返回true，反之为false和错误信息

    







=======
# 微信&支付宝支付
微信只实现了APP支付、微信公众号支付，以及扫码支付，支付结果通知的方法的功能），
>>>>>>> b6990b113273f873662df1f86b4ffe9527c48d1b
