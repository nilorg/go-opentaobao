package opentaobao

import "net/http"

// Options 可选参数列表
type Options struct {
	Router       string // 环境请求地址
	AppKey       string // 应用Key
	AppSecret    string // 秘密
	JSONSimplify bool   // 是否开启JSON返回值的简化模式
	httpClient   *http.Client
}

// Option 为可选参数赋值的函数
type Option func(*Options)

// NewOptions 创建可选参数
func NewOptions(opts ...Option) Options {
	opt := Options{
		Router:       "http://gw.api.taobao.com/router/rest",
		JSONSimplify: true,
		httpClient:   &http.Client{},
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// WithRouter 设置环境请求地址
func WithRouter(router string) Option {
	return func(o *Options) {
		o.Router = router
	}
}

// WithAppKey 设置应用Key
func WithAppKey(appKey string) Option {
	return func(o *Options) {
		o.AppKey = appKey
	}
}

// WithAppSecret 设置秘密
func WithAppSecret(appSecret string) Option {
	return func(o *Options) {
		o.AppSecret = appSecret
	}
}

// WithJSONSimplify 设置是否开启JSON返回值的简化模式
func WithJSONSimplify(jsonSimplify bool) Option {
	return func(o *Options) {
		o.JSONSimplify = jsonSimplify
	}
}

// WithHTTPClient 设置HTTP客户端
func WithHTTPClient(httpClient *http.Client) Option {
	return func(o *Options) {
		o.httpClient = httpClient
	}
}
