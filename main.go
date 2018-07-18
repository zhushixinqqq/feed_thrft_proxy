// Copyright [2013-2017] <xxxxx.Inc>
//
// Author: zhushixin
// feed_proxy server main function

package main

import (
	commons "common"
	"flag"
	"fmt"
	"modules/feed_thrft_proxy/g"
	"modules/feed_thrft_proxy/service"
	"os"
)

type Arith int

// config
var (
	flagNetworkAddr = flag.String("svrConfFile", "cfg/cfg.json", "config file contains thrift and etcd")
)

func main() {
	version := flag.Bool("version", false, "show version")
	logConfigFile := flag.String("lcfg", "seelog.xml", "seelog config file")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	commons.LoggerInit(*logConfigFile)
	//读取配置
	g.ParseConfig(*flagNetworkAddr)
	if len(g.Config().ThriftServer.Address) > 0 {
		commons.Logger.Info("global config address about thriftserver is:", g.Config().ThriftServer.Address)
	}
	//启动thrift server
	go service.Thrift_Start()
	select {}
}
