package main

//
// Author: zhushixin

import (
	"flag"
	"fmt"
	"protocols"

	"github.com/luci/go-render/render"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx"
	"github.com/smallnest/rpcx/codec"
	"github.com/smallnest/rpcx/plugin"
)

// type Args struct {
// Req protocols.QueryFeedLiteReq
// }

// type Reply struct {
// Resp protocols.QueryFeedLiteResp
// }

type Arith int

// func (t *Arith) QueryRealFeed(req protocols.QueryRealFeedReq, resp *protocols.QueryFeedLiteResp) error {
// defer glog.Flush()
// // TODO(nathan) rpc访问次数上报监控
// glog.Infoln("QueryRealFeed enter.", render.Render(req))

// resp.RealFeeds = make(map[protocols.StockTiny]*protocols.RealFeed)
// errInfo := protocols.ContructErrInfo(protocols.ErrCode_StockNotFound, "error test")
// resp.ErrInfo = &errInfo

// // realDataManager.QueryRealFeed(req, resp)

// glog.Infoln("QueryRealFeed resp.", render.Render(resp))
// return nil
// }

func (t *Arith) QueryFeedLite(req protocols.QueryFeedLiteReq, resp *protocols.QueryFeedLiteResp) error {
	errinfo := &protocols.ErrInfo{
		ErrCode: protocols.ErrCode_SUCCESS,
		ErrMsg:  "ErrTest",
	}
	resp.ErrInfo = errinfo
	feedlites := make(map[protocols.StockTiny]*(protocols.FeedLite))
	feedlite_ := &protocols.FeedLite{
		RealFeed: &protocols.RealFeed{
			StockName:   "00701",
			LastPrice:   11.22,
			PriceChange: 33.33,
		},
		IsTrade: true,
	}
	var stocktiny protocols.StockTiny
	stocktiny.StockId = "00700name"
	stocktiny.StockType = protocols.StockType_HK
	feedlites[stocktiny] = feedlite_

	feedlites_ := &protocols.FeedLite{
		RealFeed: &protocols.RealFeed{
			StockName:   "00702name",
			LastPrice:   99.22,
			PriceChange: 102.11,
		},
	}
	var stocktinys protocols.StockTiny
	stocktinys.StockId = "00702"
	stocktinys.StockType = protocols.StockType_HK
	feedlites[stocktinys] = feedlites_

	resp.Feeds = feedlites
	fmt.Println("req:", render.Render(req))
	fmt.Println("resp:", render.Render(resp))
	return nil
}

// func (t *Arith1) Error(request *protocols.QueryFeedLiteReq, resp *protocols.QueryFeedLiteResp) error {
// panic("ERROR")
// }

var addr = flag.String("s", "0.0.0.0:8975", "service address")
var e = flag.String("e", "http://10.1.2.151:8974", "etcd URL")
var n = flag.String("n", "Arith", "Service Name")

func main() {
	flag.Parse()

	server := rpcx.NewServer()
	server.ServerCodecFunc = codec.NewGobServerCodec
	rplugin := &plugin.EtcdV3RegisterPlugin{
		ServiceAddress:      "tcp@" + *addr,
		EtcdServers:         []string{*e},
		BasePath:            "/rpcx",
		Metrics:             metrics.NewRegistry(),
		UpdateIntervalInSec: 60,
	}
	rplugin.Start()
	server.PluginContainer.Add(rplugin)
	server.RegisterName(*n, new(Arith), "weight=1&m=devops")
	server.Serve("tcp", *addr)

}
