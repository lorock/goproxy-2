package api

import (
	"time"

	"github.com/valyala/fasthttp"
)

type CachePool interface {
	Get(uri string) Cache
	Delete(uri string)
	CheckAndStore(uri string, ctx *fasthttp.RequestCtx)
	Clear(d time.Duration)
}

type Cache interface {
	Verify() bool
	WriteTo(rw *fasthttp.Response) (int64, error)
}
