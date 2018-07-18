// Copyright [2013-2017] <xxxxx.Inc>
//
// Author: zhushixin

package main

import (
	"bytes"
	"flag"
	"log"
	"net"
	common "protocols/thrift_src/gen-go/common"
	feed_svr "protocols/thrift_src/gen-go/feed_svr"
	"runtime"
	"strconv"
	"testing"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/luci/go-render/render"
)

var (
	flagHost = flag.String("h", "127.0.0.1", "IP")
	flagPort = flag.String("p", "8777", "Port")
)

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
func GetFeedLite() {
	start := time.Now()
	tSocket, err := thrift.NewTSocket(net.JoinHostPort(*flagHost, *flagPort))
	if err != nil {
		log.Fatalln("tSocket error:", err)
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	transport := transportFactory.GetTransport(tSocket)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	client := feed_svr.NewFeedSvrServiceClientFactory(transport, protocolFactory)

	if err := transport.Open(); err != nil {
		log.Fatalln("Error opening:", *flagHost+":"+*flagPort)
	}
	defer transport.Close()
	var stockTinys []*common.StockTiny
	stockTiny := &common.StockTiny{
		StockId: "00700",
		TypeA1:  common.StockType_HK,
	}
	stockTinys = append(stockTinys, stockTiny)
	requestId := "GetFeedLiteStressTestId"
	accessInfo := &common.AccessInfo{
		AccessType: 0,
		AccessName: "GetFeedLiteStressTest",
		RequestId:  &requestId,
	}
	request := &feed_svr.GetFeedLiteRequest{}
	request.StockTiny = stockTinys
	request.AccessType = accessInfo
	resp, err := client.GetFeedLite(request)
	_ = resp
	guid := GetGID()
	time.Sleep(time.Second)
	log.Println("guid:", guid, " getfeedlite cost:", time.Since(start))
}

func Test_GetStockTypeTradeStatus(t *testing.T) {
	flag.Parse()
	start := time.Now()
	tSocket, err := thrift.NewTSocket(net.JoinHostPort(*flagHost, *flagPort))
	if err != nil {
		log.Fatalln("tSocket error:", err)
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	transport := transportFactory.GetTransport(tSocket)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	client := feed_svr.NewFeedSvrServiceClientFactory(transport, protocolFactory)

	if err := transport.Open(); err != nil {
		log.Fatalln("Error opening:", *flagHost+":"+*flagPort)
	}
	request := &feed_svr.GetTradeStatusReq{
		StockType: common.StockType_HK,
	}
	resp, _ := client.GetStockTypeTradeStatus(request)
	log.Println(" status cost:", time.Since(start))
	log.Println(" status ", render.Render(resp))
}

func TimeCost(funcName string) func() {
	start := time.Now()
	return func() {
		log.Println("getfeedlite cost:", time.Since(start))
		// glog.Info(funcName, " cost:", time.Since(start))
	}

}

func Benchmark_GetFeedLite(b *testing.B) {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.N = 100
	// productId := tproductId
	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			GetFeedLite()
		}
	})
}

func GetFeedEx() {
	start := time.Now()
	tSocket, err := thrift.NewTSocket(net.JoinHostPort(*flagHost, *flagPort))
	if err != nil {
		log.Fatalln("tSocket error:", err)
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	transport := transportFactory.GetTransport(tSocket)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	client := feed_svr.NewFeedSvrServiceClientFactory(transport, protocolFactory)

	if err := transport.Open(); err != nil {
		log.Fatalln("Error opening:", *flagHost+":"+*flagPort)
	}
	defer transport.Close()
	var stockTinys []*common.StockTiny
	stockTiny := &common.StockTiny{
		StockId: "00700",
		TypeA1:  common.StockType_HK,
	}
	stockTinys = append(stockTinys, stockTiny)
	requestId := "GetFeedExStressTestId"

	accessInfo := &common.AccessInfo{
		AccessType: 0,
		AccessName: "GetFeedExStressTest",
		RequestId:  &requestId,
	}
	request := &feed_svr.GetFeedExRequest{}
	request.StockTiny = stockTinys
	request.AccessType = accessInfo
	resp, err := client.GetFeedEx(request)
	_ = resp
	log.Println("getfeedex cost:", time.Since(start))
}

func Benchmark_GetFeedEx(b *testing.B) {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.N = 100
	// productId := tproductId
	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			GetFeedEx()
		}
	})
}
