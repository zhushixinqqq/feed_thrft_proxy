// Copyright [2013-2017] <xxxxx.Inc>

//
// Author: zhushixin
package rpc

import (
	commons "common"
	"flag"
	"modules/feed_thrft_proxy/g"
	"time"

	"github.com/smallnest/pool"
	"github.com/smallnest/rpcx"
	"github.com/smallnest/rpcx/clientselector"
	"github.com/smallnest/rpcx/codec"
)

type RpcBaseClient struct {
	clientPool  *pool.Pool
	servicename string
	etcdurl     string
	serviceAddr string
}

// serviceName: if not etcd serviceUrl
func NewRpcBaseClient(serviceName, serviceAddr string) *RpcBaseClient {
	flag.Parse()
	etcdurl := g.Config().EtcdServer.EtcdUrl
	if etcdurl == "" {
		commons.Logger.Critical("get the g.Config().EtcdServer.EtcdUrl failed:")
	}
	if serviceName == "" {
		commons.Logger.Critical("get the g.Config().EtcdServer.ServiceName failed:")
	}
	basePath := g.Config().EtcdServer.BasePath
	commons.Logger.Info(etcdurl, basePath, serviceName)
	var s rpcx.ClientSelector
	if g.Config().EtcdServer.Used {
		s = clientselector.NewEtcdV3ClientSelector([]string{etcdurl}, basePath+serviceName, time.Minute, rpcx.RandomSelect, time.Duration(10)*time.Second)
	} else {
		s = &rpcx.DirectClientSelector{
			Network:     "tcp",
			Address:     serviceAddr, // url
			DialTimeout: time.Duration(10) * time.Second,
		}
	}
	clientPool := &pool.Pool{
		New: func() interface{} {
			return rpcx.NewClient(s)
		},
	}
	return &RpcBaseClient{
		clientPool:  clientPool,
		servicename: serviceName,
		etcdurl:     etcdurl,
		serviceAddr: serviceAddr,
	}
}

func (c *RpcBaseClient) Get() *rpcx.Client {
	client := c.clientPool.Get().(*rpcx.Client)
	client.ClientCodecFunc = codec.NewGobClientCodec
	client.FailMode = rpcx.Failover
	client.Retries = 3
	return client
}

func (c *RpcBaseClient) Put(client *rpcx.Client) {
	c.clientPool.Put(client)
}
