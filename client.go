package opentaobao

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	simplejson "github.com/bitly/go-simplejson"
)

var (
	// AppKey 应用Key
	AppKey string
	// AppSecret 秘密
	AppSecret string
	// Router 环境请求地址
	Router string
	// Timeout ...
	Timeout time.Duration
	// CacheExpiration 缓存过期时间
	CacheExpiration = time.Hour
	// GetCache 获取缓存
	GetCache GetCacheFunc
	// SetCache 设置缓存
	SetCache SetCacheFunc
)

// Parameter 参数
type Parameter map[string]string

// copyParameter 复制参数
func copyParameter(srcParams Parameter) Parameter {
	newParams := make(Parameter)
	for key, value := range srcParams {
		newParams[key] = value
	}
	return newParams
}

// newCacheKey 创建缓存Key
func newCacheKey(params Parameter) string {
	cpParams := copyParameter(params)
	delete(cpParams, "session")
	delete(cpParams, "timestamp")
	delete(cpParams, "sign")

	keys := []string{}
	for k := range cpParams {
		keys = append(keys, k)
	}
	// 排序asc
	sort.Strings(keys)
	// 把所有参数名和参数值串在一起
	cacheKeyBuf := new(bytes.Buffer)
	for _, k := range keys {
		cacheKeyBuf.WriteString(k + "=" + cpParams[k])
	}
	h := md5.New()
	io.Copy(h, cacheKeyBuf)
	return hex.EncodeToString(h.Sum(nil))
}

// execute 执行API接口
func execute(param Parameter) (bytes []byte, err error) {
	err = checkConfig()
	if err != nil {
		return
	}

	var req *http.Request
	req, err = http.NewRequest("POST", Router, strings.NewReader(param.getRequestData()))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	httpClient := &http.Client{}
	httpClient.Timeout = Timeout
	var response *http.Response
	response, err = httpClient.Do(req)
	if err != nil {
		return
	}

	if response.StatusCode != 200 {
		err = fmt.Errorf("请求错误:%d", response.StatusCode)
		return
	}
	defer response.Body.Close()
	bytes, err = ioutil.ReadAll(response.Body)
	return
}

// Execute 执行API接口
func Execute(method string, param Parameter) (res *simplejson.Json, err error) {
	param["method"] = method
	param.setRequestData()

	var bodyBytes []byte
	bodyBytes, err = execute(param)
	if err != nil {
		return
	}

	return bytesToResult(bodyBytes)
}

func bytesToResult(bytes []byte) (res *simplejson.Json, err error) {
	res, err = simplejson.NewJson(bytes)
	if err != nil {
		return
	}

	if responseError, ok := res.CheckGet("error_response"); ok {
		if subMsg, subOk := responseError.CheckGet("sub_msg"); subOk {
			err = errors.New(subMsg.MustString())
		} else {
			err = errors.New(responseError.Get("msg").MustString())
		}
		res = nil
	}
	return
}

// ExecuteCache 执行API接口，缓存
func ExecuteCache(method string, param Parameter) (res *simplejson.Json, err error) {
	param["method"] = method
	param.setRequestData()

	cacheKey := newCacheKey(param)
	cacheBytes := GetCache(cacheKey)
	if len(cacheBytes) > 0 {
		res, err = simplejson.NewJson(cacheBytes)
		if err == nil && res != nil {
			return
		}
	}

	var bodyBytes []byte
	bodyBytes, err = execute(param)
	if err != nil {
		return
	}
	res, err = bytesToResult(bodyBytes)
	if err != nil {
		return
	}
	ejsonBody, _ := res.MarshalJSON()
	go SetCache(cacheKey, ejsonBody, CacheExpiration)
	return
}

// 检查配置
func checkConfig() error {
	if AppKey == "" {
		return errors.New("AppKey 不能为空")
	}
	if AppSecret == "" {
		return errors.New("AppSecret 不能为空")
	}
	if Router == "" {
		return errors.New("Router 不能为空")
	}
	if GetCache == nil {
		return errors.New("GetCache 不能为空")
	}
	if SetCache == nil {
		return errors.New("SetCache 不能为空")
	}
	return nil
}

func (p Parameter) setRequestData() {
	hh, _ := time.ParseDuration("8h")
	loc := time.Now().UTC().Add(hh)
	p["timestamp"] = strconv.FormatInt(loc.Unix(), 10)
	p["format"] = "json"
	p["app_key"] = AppKey
	p["v"] = "2.0"
	p["sign_method"] = "md5"
	p["partner_id"] = "Keng"
	// 设置签名
	p["sign"] = getSign(p)
}

// 获取请求数据
func (p Parameter) getRequestData() string {
	// 公共参数
	args := url.Values{}
	// 请求参数
	for key, val := range p {
		args.Set(key, val)
	}
	return args.Encode()
}

// 获取签名
func getSign(params Parameter) string {
	// 获取Key
	keys := []string{}
	for k := range params {
		keys = append(keys, k)
	}
	// 排序asc
	sort.Strings(keys)
	// 把所有参数名和参数值串在一起
	query := AppSecret
	for _, k := range keys {
		query += k + params[k]
	}
	query += AppSecret
	// 使用MD5加密
	signBytes := md5.Sum([]byte(query))
	// 把二进制转化为大写的十六进制
	return strings.ToUpper(hex.EncodeToString(signBytes[:]))
}
