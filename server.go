package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/panjf2000/goproxy/config"
	"github.com/panjf2000/goproxy/handlers"
	_ "net/http"
)

func main() {
	goProxy := handlers.NewProxyServer()

	log.Infof("start proxy server in port%s", config.RuntimeViper.GetString("server.port"))
	log.Fatal(goProxy.ListenAndServe())
}
