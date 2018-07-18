// Copyright [2013-2017] <xxxxx.Inc>
//
// Author: zhushixin

package service

import (
	commons "common"
	"modules/feed_thrft_proxy/g"
	feed_svr "protocols/thrift_src/gen-go/feed_svr"

	"git.apache.org/thrift.git/lib/go/thrift"
)

func runServer(transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory, addr string) error {
	transport, err := thrift.NewTServerSocket(addr)
	if err != nil {
		commons.Logger.Critical("runServer transport failed:%s", err)
		return err
	}
	handler := NewAccessLayer()
	handler.Init()
	processor := feed_svr.NewFeedSvrServiceProcessor(handler)
	thriftServer := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
	commons.Logger.Info("feed thrift proxy runServer..........................")
	return thriftServer.Serve()
}
func Thrift_Start() {
	if !g.Config().ThriftServer.Enabled {
		commons.Logger.Critical("Thrift_Start function g.Config().ThriftServer.Enabledfailed")
		return
	}
	if g.Config().ThriftServer.Address == "" {
		commons.Logger.Critical("get the ThriftServer.Address failed:")
		return
	}

	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	if err := runServer(transportFactory, protocolFactory, g.Config().ThriftServer.Address); err != nil {
		commons.Logger.Critical("runServer failed: %s", err)
	}
}
