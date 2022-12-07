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
	// taobao.tbk.dg.material.optional
	result, err := opentaobao.Execute("taobao.tbk.dg.material.optional", opentaobao.Parameter{
		"q":         "鸿星尔克男鞋板鞋2022新款白色厚底休闲鞋子潮时尚男运动鞋滑板鞋",
		"adzone_id": "466816091",
		"platform":  "2",
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
