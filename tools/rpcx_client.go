// Copyright [2013-2017] <xxxxx.Inc>
//
// Author: zhushixin

package main

import (
	"context"
	"flag"
	"fmt"
	"protocols"
	"time"

	"github.com/smallnest/rpcx"
	"github.com/smallnest/rpcx/clientselector"
	"github.com/smallnest/rpcx/codec"
)

type Args struct {
	Req protocols.QueryFeedLiteReq
}

type Reply struct {
	Resp protocols.QueryFeedLiteResp
}

// var e = flag.String("e", "http://10.1.2.151:8974", "etcd URL")

var e = flag.String("e", "http://10.1.2.151:8974", "etcd URL")
var n = flag.String("ns", "Arith", "Service Name")

func CreateClientSelector() *rpcx.Client {
	var s rpcx.ClientSelector
	s = clientselector.NewEtcdV3ClientSelector([]string{*e}, "/rpcx/"+*n, time.Minute, rpcx.RandomSelect, time.Minute)
	client := rpcx.NewClient(s)
	client.ClientCodecFunc = codec.NewGobClientCodec

	fmt.Println("CreateClientSelector")
	return client
}
func Make() {
	var client *rpcx.Client
	client = CreateClientSelector()
	args := &Args{}
	req := &protocols.QueryFeedLiteReq{}
	accessionfo := &protocols.AccessInfo{AccessType: "accesstype", RequestId: "0000"}
	// stocktiny := protocols.StockTiny{StockId: "00700", StockType: protocols.StockType_HK}
	req.AccessInfo = accessionfo
	args.Req = *req
	var reply Reply
	err := client.Call(context.Background(), *n+".GetFeedLite", args, &reply)
	if err != nil {
		fmt.Println("Call is err", err)
	}
	fmt.Println(reply)

	client.Close()

}

func main() {
	Make()
}
