package handlers

import (
	"github.com/valyala/fasthttp"
)

// CopyHeaders copy headers from source to destination.
// Nothing would be returned.
func CopyHeaders(dst, src *fasthttp.ResponseHeader) {
	src.CopyTo(dst)
}

// RmProxyHeaders remove Hop-by-hop headers.
func RmProxyReqHeaders(req *fasthttp.Request) {
	req.SetRequestURI("")
	req.Header.Del("Proxy-Connection")
	req.Header.Del("Connection")
	req.Header.Del("Keep-Alive")
	req.Header.Del("Proxy-Authenticate")
	req.Header.Del("Proxy-Authorization")
	req.Header.Del("TE")
	req.Header.Del("Trailers")
	req.Header.Del("Transfer-Encoding")
	req.Header.Del("Upgrade")
}

func RmProxyRespHeaders(resp *fasthttp.Response) {
	resp.Header.Del("Proxy-Connection")
	resp.Header.Del("Connection")
	resp.Header.Del("Keep-Alive")
	resp.Header.Del("Proxy-Authenticate")
	resp.Header.Del("Proxy-Authorization")
	resp.Header.Del("TE")
	resp.Header.Del("Trailers")
	resp.Header.Del("Transfer-Encoding")
	resp.Header.Del("Upgrade")
}
