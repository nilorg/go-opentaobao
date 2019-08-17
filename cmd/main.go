package main

import (
	"log"
	"os"
	"time"

	"github.com/nilorg/go-opentaobao"
)

func init() {
	// opentaobao.AppKey = "0000"
	opentaobao.AppKey = os.Getenv("OPEN_TAOBAO_APPKEY")
	// opentaobao.AppSecret = "xxxxx"
	opentaobao.AppSecret = os.Getenv("OPEN_TAOBAO_APPSECRET")
	opentaobao.GetCache = func(cacheKey string) []byte {
		return nil
	}
	opentaobao.SetCache = func(cacheKey string, value []byte, expiration time.Duration) bool {
		return true
	}
}

func main() {

	result, err := opentaobao.Execute("taobao.tbk.relation.refund", opentaobao.Parameter{
		"search_option": map[string]interface{}{
			"page_size":   1,
			"search_type": 4, // 1-维权发起时间，2-订单结算时间（正向订单），3-维权完成时间，4-订单创建时间
			"refund_type": 1, // 1 表示2方，2表示3方
			"start_time":  "2019-07-08 00:00:00",
			"page_no":     1,
			"biz_type":    1,
		},
	})

	// result, err := opentaobao.Execute("taobao.tbk.order.details.get", opentaobao.Parameter{
	// 	"start_time": "2019-06-17 22:00:00",
	// 	"end_time":   "2019-06-17 23:00:00",
	// })

	if err != nil {
		log.Printf("execute error:%s\n", err)
		return
	}
	data, _ := result.MarshalJSON()
	log.Printf("result:%s\n", data)
}
