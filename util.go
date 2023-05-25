package opentaobao

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"sort"
	"strings"

	"github.com/nilorg/sdk/convert"
)

// 获取签名
func Sign(params Parameter, appSecret string) string {
	// 获取Key
	keys := []string{}
	for k := range params {
		keys = append(keys, k)
	}
	// 排序asc
	sort.Strings(keys)
	// 把所有参数名和参数值串在一起
	query := bytes.NewBufferString(appSecret)
	for _, k := range keys {
		query.WriteString(k)
		query.WriteString(interfaceToString(params[k]))
	}
	query.WriteString(appSecret)
	// 使用MD5加密
	h := md5.New()
	io.Copy(h, query)
	// 把二进制转化为大写的十六进制
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

func interfaceToString(src interface{}) string {
	if src == nil {
		panic(ErrTypeIsNil)
	}
	switch v := src.(type) {
	case string:
		return v
	case uint8, uint16, uint32, uint64, int, int8, int32, int64, float32, float64:
		return convert.ToString(v)
	}
	data, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}
	return string(data)
}
