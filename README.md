# taobaogo
淘宝Api、淘宝开放平台Api请求基础SDK

# 淘宝API

[sign算法](http://open.taobao.com/doc.htm?docId=101617&docType=1)

[淘宝Session](https://oauth.taobao.com/authorize?response_type=token&client_id=24840730)

# Example 🌰
```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/nilorg/go-opentaobao/v2"
)

func main() {
	client := opentaobao.NewClient(
		opentaobao.WithAppKey(os.Getenv("APP_KEY")),
		opentaobao.WithAppSecret(os.Getenv("APP_SECRET")),
	)
	ctx := context.Background()
	// EXP: 使用session
	// ctx = opentaobao.NewSessionContext(ctx, "session")
	result, err := client.Execute(ctx, "taobao.tbk.dg.material.optional", opentaobao.Parameter{
		"q":         "鸿星尔克男鞋板鞋",
		"adzone_id": os.Getenv("ADZONE_ID"),
		"platform":  "2",
	})

	if err != nil {
		log.Printf("execute error:%s\n", err)
		return
	}
	log.Printf("result:%s\n", result.String())
}
```
