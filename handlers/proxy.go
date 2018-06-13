package handlers

import (
	_ "bufio"
	"io"
	"net"
	"net/http"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/panjf2000/goproxy/cache"
	"github.com/panjf2000/goproxy/config"
	"github.com/panjf2000/goproxy/tool"
	"github.com/valyala/fasthttp"
)

type ProxyServer struct {
	Addr string
	// Browser records user's name
	Browser string
	Client  *fasthttp.HostClient
	Server  *fasthttp.Server
}

var proxyLog *logrus.Logger

func init() {
	logPath := config.RuntimeViper.GetString("server.log_path")
	proxyLog, _ = tool.InitLog(path.Join(logPath, "proxy.log"))

}

// NewProxyServer returns a new proxyserver.
func NewProxyServer() *ProxyServer {
	if config.RuntimeViper.GetBool("server.cache") {
		RegisterCachePool(cache.NewCachePool(config.RuntimeViper.GetString("redis.redis_host"),
			config.RuntimeViper.GetString("redis.redis_pass"), config.RuntimeViper.GetInt("redis.idle_timeout"),
			config.RuntimeViper.GetInt("redis.max_active"), config.RuntimeViper.GetInt("redis.max_idle")))
	}

	//server := &fasthttp.Server{
	//	Name:         config.RuntimeViper.GetString("server.name"),
	//	Handler:      proxyServer.HandleFastHTTP,
	//	ReadTimeout:  10 * time.Second,
	//	WriteTimeout: 10 * time.Second,
	//}
	proxyServer := &ProxyServer{
		Addr: config.RuntimeViper.GetString("server.port"),
		Client: &fasthttp.HostClient{
			IsTLS: false,
			Addr:  "",
			// set other options here if required - most notably timeouts.
			// ReadTimeout: 60, // 如果在生产环境启用会出现多次请求现象
		},
	}
	return proxyServer
	//return &http.Server{
	//	Addr:           config.RuntimeViper.GetString("server.port"),
	//	Handler:        &ProxyServer{Travel: &http.Transport{Proxy: http.ProxyFromEnvironment, DisableKeepAlives: true}},
	//	ReadTimeout:    10 * time.Second,
	//	WriteTimeout:   10 * time.Second,
	//	MaxHeaderBytes: 1 << 20,
	//}
}
func (ps *ProxyServer) ListenAndServe() error {
	return fasthttp.ListenAndServe(ps.Addr, ps.HandleFastHTTP)
}

func (ps *ProxyServer) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	proxyLog.WithFields(logrus.Fields{
		"request user":   ps.Browser,
		"request method": string(ctx.Method()),
		"request url":    string(ctx.Host()),
	}).Info("request's detail !")
	req := &ctx.Request
	resp := &ctx.Response
	RmProxyReqHeaders(req)
	if string(ctx.Method()) == http.MethodGet && config.RuntimeViper.GetBool("server.cache") {
		ps.CacheHandler(ctx)
	} else {
		if err := ps.Client.Do(req, resp); err != nil {
			proxyLog.WithError(err)
		}
	}
	RmProxyRespHeaders(resp)

	proxyLog.WithFields(logrus.Fields{
		"response bytes": len(resp.Body()),
		"request url":    req.URI().String(),
	}).Info("response has been copied successfully!")
	ps.Done(req)
}

var HTTP200 = []byte("HTTP/1.1 200 Connection Established\r\n\r\n")

// HttpsHandler handles any connection which need connect method.
// 处理https连接，主要用于CONNECT方法
func (ps *ProxyServer) HttpsHandler(rw http.ResponseWriter, req *http.Request) {
	proxyLog.WithFields(logrus.Fields{
		"user": ps.Browser,
		"host": req.URL.Host,
	}).Info("http user tried to connect host!")

	hj, _ := rw.(http.Hijacker)
	Client, _, err := hj.Hijack() //获取客户端与代理服务器的tcp连接
	if err != nil {
		proxyLog.WithFields(logrus.Fields{
			"user":        ps.Browser,
			"request uri": req.RequestURI,
		}).Error("http user failed to get tcp connection!")
		http.Error(rw, "Failed", http.StatusBadRequest)
		return
	}

	Remote, err := net.Dial("tcp", req.URL.Host) //建立服务端和代理服务器的tcp连接
	if err != nil {
		proxyLog.WithFields(logrus.Fields{
			"user":        ps.Browser,
			"request uri": req.RequestURI,
		}).Error("http user failed to connect this uri!")
		http.Error(rw, "Failed", http.StatusBadGateway)
		return
	}

	Client.Write(HTTP200)

	go copyRemoteToClient(ps.Browser, Remote, Client)
	go copyRemoteToClient(ps.Browser, Client, Remote)
}

func copyRemoteToClient(User string, Remote, Client net.Conn) {
	defer func() {
		Remote.Close()
		Client.Close()
	}()

	nr, err := io.Copy(Remote, Client)
	if err != nil && err != io.EOF {
		proxyLog.WithFields(logrus.Fields{
			"client": User,
			"error":  err,
		}).Error("occur an error when handling CONNECT Method")
		return
	}
	proxyLog.WithFields(logrus.Fields{
		"user":           User,
		"nr":             nr,
		"remote_address": Remote.RemoteAddr(),
		"client_address": Client.RemoteAddr(),
	}).Info("transport the bytes between client and remote!")
}
