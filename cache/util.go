package cache

import (
	"log"
	"strings"

	"github.com/valyala/fasthttp"
)

//checks whether request ask to be stored as cache
func IsReqCache(req *fasthttp.Request) bool {
	log.Printf("http request header:%v", req.Header)
	cacheControl := string(req.Header.Peek("Cache-Control"))
	contentType := string(req.Header.Peek("Content-Type"))
	if cacheControl == "" && contentType == "" {
		return true
	} else if len(cacheControl) > 0 {
		if strings.Index(cacheControl, "private") != -1 ||
			strings.Index(cacheControl, "no-cache") != -1 ||
			strings.Index(cacheControl, "no-store") != -1 ||
			strings.Index(cacheControl, "must-revalidate") != -1 ||
			(strings.Index(cacheControl, "max-age") == -1 &&
				strings.Index(cacheControl, "s-maxage") == -1 &&
				string(req.Header.Peek("Etag")) == "" &&
				string(req.Header.Peek("Last-Modified")) == "" &&
				(string(req.Header.Peek("Expires")) == "" || string(req.Header.Peek("Expires")) == "0")) {
			return false
		}

	} else if len(contentType) > 0 {
		if strings.Index(contentType, "video") != -1 ||
			strings.Index(contentType, "image") != -1 ||
			strings.Index(contentType, "audio") != -1 {
			return false
		}

	}
	return true
}

//checks whether response can be stored as cache
func IsRespCache(resp *fasthttp.Response) bool {
	log.Printf("http response header:%v", resp.Header)
	cacheControl := string(resp.Header.Peek("Cache-Control"))
	contentType := string(resp.Header.Peek("Content-Type"))
	if cacheControl == "" && contentType == "" {
		return true
	} else if len(cacheControl) > 0 {
		if strings.Index(cacheControl, "private") != -1 ||
			strings.Index(cacheControl, "no-cache") != -1 ||
			strings.Index(cacheControl, "no-store") != -1 ||
			strings.Index(cacheControl, "must-revalidate") != -1 ||
			(strings.Index(cacheControl, "max-age") == -1 &&
				strings.Index(cacheControl, "s-maxage") == -1 &&
				string(resp.Header.Peek("Etag")) == "" &&
				string(resp.Header.Peek("Last-Modified")) == "" &&
				(string(resp.Header.Peek("Expires")) == "" || string(resp.Header.Peek("Expires")) == "0")) {
			return false
		}

	} else if len(contentType) > 0 {
		if strings.Index(contentType, "video") != -1 ||
			strings.Index(contentType, "image") != -1 ||
			strings.Index(contentType, "audio") != -1 {
			return false
		}

	}
	return true
}
