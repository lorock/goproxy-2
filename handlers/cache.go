package handlers

import (
	_ "bufio"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/panjf2000/goproxy/config"
	"github.com/panjf2000/goproxy/interface"
	"github.com/panjf2000/goproxy/tool"
	"github.com/valyala/fasthttp"
)

var cachePool api.CachePool
var cacheLog *logrus.Logger

func init() {
	logPath := config.RuntimeViper.GetString("server.log_path")
	cacheLog, _ = tool.InitLog(path.Join(logPath, "cache.log"))
}

// RegisterCachePool register a new cache-pool.
func RegisterCachePool(c api.CachePool) {
	cachePool = c
}

//CacheHandler handles "Get" request
func (ps *ProxyServer) CacheHandler(ctx *fasthttp.RequestCtx) {
	req := &ctx.Request
	resp := &ctx.Response

	var uri = string(req.RequestURI())

	c := cachePool.Get(uri)

	if c != nil {
		if c.Verify() {
			cacheLog.WithFields(logrus.Fields{
				"request url": uri,
			}).Debug("Found cache!")
			c.WriteTo(resp)
			return
		} else {
			cacheLog.WithFields(logrus.Fields{
				"request url": uri,
			}).Debug("Delete cache!")
			cachePool.Delete(uri)
		}
	}

	RmProxyReqHeaders(req)
	if err := ps.Client.Do(req, resp); err != nil {
		proxyLog.WithError(err)
		return
	}

	newResp := new(fasthttp.Response)
	resp.CopyTo(newResp)

	cacheLog.WithFields(logrus.Fields{
		"request url": uri,
	}).Debug("Check out this cache and then stores it if it is right!")
	go cachePool.CheckAndStore(uri, ctx)

	resp.Header.Reset()
	CopyHeaders(&resp.Header, &newResp.Header)

	cacheLog.WithFields(logrus.Fields{
		"response bytes": len(newResp.Body()),
		"request url":    req.URI().String(),
	}).Info("response has been copied successfully!")
}
