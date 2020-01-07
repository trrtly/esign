package esign

import (
	"sync"

	"github.com/trrtly/esign/cache"
	"github.com/trrtly/esign/context"
	"github.com/trrtly/esign/verify"
)

// Esign struct
type Esign struct {
	Context *context.Context
}

// Config for user
type Config struct {
	Appid  string
	Secret string
	Debug  bool

	Cache cache.Cache
}

// NewEsign init
func NewEsign(c *Config) *Esign {
	ctx := new(context.Context)
	copyConfigToContext(c, ctx)
	return &Esign{ctx}
}

func copyConfigToContext(c *Config, context *context.Context) {
	context.Appid = c.Appid
	context.Secret = c.Secret
	context.Debug = c.Debug
	context.Cache = c.Cache
	context.SetAccessTokenLock(new(sync.RWMutex))
}

//GetAccessToken 获取access_token
func (es *Esign) GetAccessToken() (string, error) {
	return es.Context.GetAccessToken()
}

//GetVerify 认证
func (es *Esign) GetVerify() *verify.Verify {
	return verify.NewVerify(es.Context)
}
