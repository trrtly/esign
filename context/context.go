package context

import (
	"sync"

	"github.com/trrtly/esign/cache"
)

// Context struct
type Context struct {
	Appid  string
	Secret string
	Debug  bool

	Cache cache.Cache

	//accessTokenLock 读写锁 同一个Appid一个
	accessTokenLock *sync.RWMutex

	//accessTokenFunc 自定义获取 access token 的方法
	accessTokenFunc GetAccessTokenFunc
}

// CommonError struct
type CommonError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

const (
	apiDomainProduct = "https://openapi.esign.cn"
	apiDomainDebug   = "https://smlopenapi.esign.cn"
)

// GetDomain get root domain
func (ctx *Context) GetDomain() string {
	if ctx.Debug {
		return apiDomainDebug
	}
	return apiDomainProduct
}

// GetRequestHeader 通用请求头参数
func (ctx *Context) GetRequestHeader() (header map[string]string, err error) {
	accessToken, err := ctx.GetAccessToken()
	if err != nil {
		return
	}
	header = map[string]string{
		"X-Tsign-Open-App-Id": ctx.Appid,
		"X-Tsign-Open-Token":  accessToken,
	}
	return
}
