// Copyright [2013-2017] <xxxxx.Inc>
//
// Author: zhushixin
// rpc interface communicate with the rpcx return to the feed_hander

package rpc

import (
	commons "common"
	"context"
	"modules/feed_thrft_proxy/g"
	"protocols"

	"github.com/luci/go-render/render"
)

type RpcManager struct {
	rpcBaseClient  *RpcBaseClient
	rpcChartClient *RpcBaseClient
	// client :=&rpcx.Client{}
}

func NewRpcManager() *RpcManager {
	return &RpcManager{}
}

func (r *RpcManager) Init() {
	commons.Logger.Info("init RpcManager!")
	r.rpcBaseClient = NewRpcBaseClient(g.Config().EtcdServer.ServiceName, g.Config().ThriftServer.RealBusAddr)
	r.rpcChartClient = NewRpcBaseClient(g.Config().EtcdServer.ChartServiceName, g.Config().ThriftServer.ChartBusAddr)
}

func (r *RpcManager) GetFeedLite(request protocols.QueryFeedLiteReq, resp *protocols.QueryFeedLiteResp) (err error) {
	defer commons.TimeCostStr("RPCQueryGetFeedLite")()
	commons.Logger.Debug("GetFeedLite request:", render.Render(request))
	client := r.rpcBaseClient.Get()
	defer r.rpcBaseClient.Put(client)
	err = client.Call(context.Background(), r.rpcBaseClient.servicename+".QueryFeedLite", request, resp)
	if err != nil {
		commons.Logger.Error("err:", err.Error())
	}
	commons.Logger.Debug("resp:", render.Render(resp))
	return err
}

func (r *RpcManager) GetFeedEx(request protocols.QueryFeedReq, resp *protocols.QueryFeedResp) (err error) {
	defer commons.TimeCostStr("RPCQueryGetFeedEx")()
	commons.Logger.Debug("GetFeedEx request:", render.Render(request))
	client := r.rpcBaseClient.Get()
	defer r.rpcBaseClient.Put(client)
	err = client.Call(context.Background(), r.rpcBaseClient.servicename+".QueryFeed", request, resp)
	if err != nil {
		commons.Logger.Error("err:", err.Error())
	}
	commons.Logger.Debug("GetFeedEx resp:", render.Render(resp))
	return err
}

func (r *RpcManager) GetMarketDepth(request protocols.QueryMarketDepthReq, resp *protocols.QueryMarketDepthResp) (err error) {
	defer commons.TimeCostStr("RPCQueryMarketDepth")()
	commons.Logger.Debug("GetMarketDepth request:", render.Render(request))
	client := r.rpcBaseClient.Get()
	defer r.rpcBaseClient.Put(client)
	err = client.Call(context.Background(), r.rpcBaseClient.servicename+".QueryMarketDepth", request, resp)
	if err != nil {
		commons.Logger.Error("err:", err.Error())
	}
	commons.Logger.Debug("GetMarketDepth resp:", render.Render(resp))
	return err
}

//获取市场上所有股票根据StockType
func (r *RpcManager) GetStocksByType(request protocols.QueryExchangeStocksReq, resp *protocols.QueryExchangeStocksResp) (err error) {
	defer commons.TimeCostStr("RPCQueryExchangeStocks")()
	commons.Logger.Debug("GetExchangeStocks request:", render.Render(request))
	client := r.rpcBaseClient.Get()
	defer r.rpcBaseClient.Put(client)
	err = client.Call(context.Background(), r.rpcBaseClient.servicename+".QueryExchangeStocks", request, resp)
	if err != nil {
		commons.Logger.Error("err:", err.Error())
	}
	commons.Logger.Debug("QueryExchangeStocks resp:", render.Render(resp))
	return err
}

func (r *RpcManager) GetKChart(request protocols.QueryKChartReq, resp *protocols.QueryKChartResp) (err error) {
	defer commons.TimeCostStr("RPCQueryKChart")()
	commons.Logger.Debug("QueryKChartReq request:", render.Render(request))
	client := r.rpcChartClient.Get()
	defer r.rpcChartClient.Put(client)
	err = client.Call(context.Background(), g.Config().EtcdServer.ChartServiceName+".QueryKChart", request, resp)
	if err != nil {
		commons.Logger.Error("err:", err.Error())
	}
	commons.Logger.Debug("RPCQueryKChart resp:", render.Render(resp))
	return err
}

func (r *RpcManager) GetTimeChart(request protocols.QueryTimeChartReq, resp *protocols.QueryTimeChartResp) (err error) {
	defer commons.TimeCostStr("RPCQueryTimeChart")()
	commons.Logger.Debug("QueryTimeChartReq request:", render.Render(request))
	client := r.rpcChartClient.Get()
	defer r.rpcChartClient.Put(client)
	err = client.Call(context.Background(), g.Config().EtcdServer.ChartServiceName+".QueryTimeChart", request, resp)
	if err != nil {
		commons.Logger.Error("err:", err.Error())
	}
	commons.Logger.Debug("RPCQueryTimeChart resp:", render.Render(resp))
	return err
}

func (r *RpcManager) GetStockTypeStatus(request protocols.GetStockTypeTradeStatusReq, resp *protocols.GetStockTypeTradeStatusResp) (err error) {
	defer commons.TimeCostStr("RPCGetStockTypeStatus:")()
	commons.Logger.Debug("GetStockTypeStatusReq request:", render.Render(request))
	client := r.rpcBaseClient.Get()
	defer r.rpcBaseClient.Put(client)
	err = client.Call(context.Background(), r.rpcBaseClient.servicename+".GetStockTypeTradeStatus", request, resp)
	if err != nil {
		commons.Logger.Error("err:", err.Error())
	}
	commons.Logger.Debug("RPCGetStockTypeStatus resp:", render.Render(resp))
	return err
}

func (r *RpcManager) GetStockBasicInfoBatch(request protocols.GetStockBasicInfoBatchReq, resp *protocols.GetStockBasicInfoBatchResp) (err error) {
	defer commons.TimeCostStr("RPCGetStockBasicInfoBatch:")()
	commons.Logger.Debug("GetStockBasicInfoBatchReq request:", render.Render(request))
	client := r.rpcBaseClient.Get()
	defer r.rpcBaseClient.Put(client)
	err = client.Call(context.Background(), r.rpcBaseClient.servicename+".GetStockBasicInfoBatch", request, resp)
	if err != nil {
		commons.Logger.Error("err:", err.Error())
	}
	commons.Logger.Debug("RPCGetStockBasicInfoBatch resp:", render.Render(resp))
	return err
}

// add by 2.3.8
func (r *RpcManager) GetStockTicks(request protocols.QueryTickerReq, resp *protocols.QueryTickerResp) (err error) {
	defer commons.TimeCostStr("RPCQueryTickerReq:")()
	commons.Logger.Debug("GetStockTicks request:", render.Render(request))
	client := r.rpcBaseClient.Get()
	defer r.rpcBaseClient.Put(client)
	//TODO
	err = client.Call(context.Background(), r.rpcBaseClient.servicename+".GetStockTicks", request, resp)
	if err != nil {
		commons.Logger.Error("err:", err.Error())
	}
	commons.Logger.Debug("RPCGetStockTicks resp:", render.Render(resp))
	return err
}
