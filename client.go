package opentaobao

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type Client struct {
	opts Options
}

func NewClient(opts ...Option) *Client {
	client := new(Client)
	client.opts = NewOptions(opts...)
	return client
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	return c.opts.httpClient.Do(req)
}

// 检查配置
func (c *Client) checkConfig() error {
	if c.opts.AppKey == "" {
		return errors.New("AppKey不能为空")
	}
	if c.opts.AppSecret == "" {
		return errors.New("AppSecret不能为空")
	}
	if c.opts.Router == "" {
		return errors.New("Router不能为空")
	}
	return nil
}

// execute 执行API接口
func (c *Client) execute(ctx context.Context, param Parameter) (bytes []byte, err error) {
	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, "POST", c.opts.Router, strings.NewReader(param.getRequestData()))
	if err != nil {
		return
	}
	var httpResp *http.Response
	httpResp, err = c.do(req)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != 200 {
		err = fmt.Errorf("http status code: %d", httpResp.StatusCode)
		return
	}
	bytes, err = io.ReadAll(httpResp.Body)
	return
}

// Execute 执行API接口
func (c *Client) Execute(ctx context.Context, method string, param Parameter) (res gjson.Result, err error) {
	err = c.checkConfig()
	if err != nil {
		return
	}
	if c.opts.JSONSimplify {
		param["simplify"] = "true"
	}
	param["method"] = method
	param.setRequestData(ctx, c.opts.AppKey, c.opts.AppSecret)

	var bodyBytes []byte
	bodyBytes, err = c.execute(ctx, param)
	if err != nil {
		return
	}
	return bytesToResult(bodyBytes)
}

// Parameter 参数
type Parameter map[string]any

func bytesToResult(bytes []byte) (res gjson.Result, err error) {
	res = gjson.ParseBytes(bytes)
	responseError := res.Get("error_response")
	if ok := responseError.Exists(); ok {
		subMsg := responseError.Get("sub_msg")
		if ok := subMsg.Exists(); ok {
			err = errors.New(subMsg.String())
		} else {
			err = errors.New(responseError.Get("msg").String())
		}
	}
	return
}

func (p Parameter) setRequestData(ctx context.Context, appKey string, appSecret string) {
	p["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	p["format"] = "json"
	p["app_key"] = appKey
	p["v"] = "2.0"
	p["sign_method"] = "md5"
	if session, ok := FromSessionContext(ctx); ok {
		p["session"] = session
	}
	// 设置签名
	p["sign"] = Sign(p, appSecret)
}

// 获取请求数据
func (p Parameter) getRequestData() string {
	// 公共参数
	args := url.Values{}
	// 请求参数
	for key, val := range p {
		args.Set(key, interfaceToString(val))
	}
	return args.Encode()
}
