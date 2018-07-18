// Copyright [2013-2017] <xxxxx.Inc>
//
// Author: zhushixin

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	common "protocols/thrift_src/gen-go/common"
	feed_svr "protocols/thrift_src/gen-go/feed_svr"
	"runtime"
	"strconv"
	"sync"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
)

var (
	flagHost  = flag.String("h", "127.0.0.1", "IP")
	flagPort  = flag.String("p", "8777", "Port")
	flagCount = flag.Int("nums", 1, "nums")
)

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
func GetFeedLite(wg *sync.WaitGroup) {
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
	end := time.Now()
	result := end.Sub(start).Nanoseconds() / 1000000
	getid := 1
	log.Print("GID:", getid, " getfeedlite cost:", result, " ms")
	wg.Done()
}

func main() {
	flag.Parse()
	fmt.Println("start")
	fmt.Println(*flagHost, " ", *flagPort, " ", *flagCount)
	runtime.GOMAXPROCS(runtime.NumCPU())
	wg := new(sync.WaitGroup)
	wg.Add(*flagCount)
	for i := 0; i < *flagCount; i++ {
		go GetFeedLite(wg)
	}
	wg.Wait()
}
