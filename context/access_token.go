package context

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/trrtly/esign/util"
)

const (
	//AccessTokenURL 获取access_token的接口
	AccessTokenURL = "%s/v1/oauth2/access_token"
)

//ResAccessToken struct
type ResAccessToken struct {
	CommonError

	Data struct {
		Token        string      `json:"token"`
		RefreshToken string      `json:"refreshToken"`
		ExpiresIn    ExpiresTime `json:"expiresIn"`
	} `json:"data"`
}

// ExpiresTime ExpiresTime
type ExpiresTime string

//GetAccessTokenFunc 获取 access token 的函数签名
type GetAccessTokenFunc func(ctx *Context) (accessToken string, err error)

// GetAccessTokenURL 获取 access_token 请求地址
func (ctx *Context) GetAccessTokenURL() string {
	return fmt.Sprintf(AccessTokenURL, ctx.GetDomain())
}

//SetAccessTokenLock 设置读写锁（一个appid一个读写锁）
func (ctx *Context) SetAccessTokenLock(l *sync.RWMutex) {
	ctx.accessTokenLock = l
}

//SetGetAccessTokenFunc 设置自定义获取accessToken的方式, 需要自己实现缓存
func (ctx *Context) SetGetAccessTokenFunc(f GetAccessTokenFunc) {
	ctx.accessTokenFunc = f
}

//GetAccessToken 获取access_token
func (ctx *Context) GetAccessToken() (accessToken string, err error) {
	ctx.accessTokenLock.Lock()
	defer ctx.accessTokenLock.Unlock()

	if ctx.accessTokenFunc != nil {
		return ctx.accessTokenFunc(ctx)
	}
	accessTokenCacheKey := fmt.Sprintf("esign_access_token_%s", ctx.Appid)
	val := ctx.Cache.Get(accessTokenCacheKey)
	if val != nil {
		accessToken = val.(string)
		return
	}

	//从服务器获取
	var resAccessToken ResAccessToken
	resAccessToken, err = ctx.GetAccessTokenFromServer()
	if err != nil {
		return
	}

	accessToken = resAccessToken.Data.Token
	return
}

//GetAccessTokenFromServer 从服务器获取token
func (ctx *Context) GetAccessTokenFromServer() (resAccessToken ResAccessToken, err error) {
	url := fmt.Sprintf("%s?grantType=client_credentials&appId=%s&secret=%s", ctx.GetAccessTokenURL(), ctx.Appid, ctx.Secret)
	var body []byte
	body, err = util.HTTPGet(url)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &resAccessToken)
	if err != nil {
		return
	}
	if resAccessToken.Code != 0 {
		err = fmt.Errorf("get access_token error : code=%v , msg=%v", resAccessToken.Code, resAccessToken.Message)
		return
	}

	accessTokenCacheKey := fmt.Sprintf("esign_access_token_%s", ctx.Appid)
	expires := resAccessToken.Data.ExpiresIn.ConverToSeconds() - time.Now().Unix() - 1500
	err = ctx.Cache.Set(accessTokenCacheKey, resAccessToken.Data.Token, time.Duration(expires)*time.Second)
	return
}

// ConverToSeconds token 过期时间转化为秒
func (exp ExpiresTime) ConverToSeconds() int64 {
	microSeconds, _ := strconv.ParseInt(string(exp), 10, 64)
	return int64(microSeconds / 1000)
}
