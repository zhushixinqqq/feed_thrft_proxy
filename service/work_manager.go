// Copyright [2013-2017] <xxxxx.Inc>
//
// Author: zhushixin
// work pool
package service

import (
	commons "common"
	"modules/feed_thrft_proxy/g"
	common "protocols/thrift_src/gen-go/common"
	feed_svr "protocols/thrift_src/gen-go/feed_svr"
	"time"
	"workpool"
)

type WorkerManager struct {
	workHandler             *WorkHandler
	feedLitePool            *tunny.WorkPool
	feedExPool              *tunny.WorkPool
	marketDepthPool         *tunny.WorkPool
	stocksByTypePool        *tunny.WorkPool
	stockTradeStatusPool    *tunny.WorkPool
	kChartPool              *tunny.WorkPool
	timeChartPool           *tunny.WorkPool
	stockBasicInfoBatchPool *tunny.WorkPool
	stockTicksPool          *tunny.WorkPool
}

func NewWorkManager() *WorkerManager {
	return &WorkerManager{}
}

func (w *WorkerManager) Init() {
	w.workHandler = NewWorkHandler()
	w.workHandler.Init()

	var err error
	w.feedLitePool, err = tunny.CreatePool(int(g.Config().ThriftServer.FeedThriftPoolNum), w.workHandler.QueryFeedLite).Open()
	if err != nil {
		commons.Logger.Error("create feedPool failed:", err.Error())
	}
	w.feedExPool, err = tunny.CreatePool(int(g.Config().ThriftServer.FeedExPoolNum), w.workHandler.QueryFeedEx).Open()
	if err != nil {
		commons.Logger.Error("create feedExPool failed:", err.Error())
	}
	w.marketDepthPool, err = tunny.CreatePool(int(g.Config().ThriftServer.MarketDepthPoolNum), w.workHandler.QueryMarketDepth).Open()
	if err != nil {
		commons.Logger.Error("create marketDepthPool failed:", err.Error())
	}
	w.stocksByTypePool, err = tunny.CreatePool(int(g.Config().ThriftServer.StocksByTypePoolNum), w.workHandler.QueryStocksByType).Open()
	if err != nil {
		commons.Logger.Error("create stocksByTypePool failed:", err.Error())
	}
	w.stockTradeStatusPool, err = tunny.CreatePool(int(g.Config().ThriftServer.StockTradeStatusPoolNum), w.workHandler.QueryStockStatus).Open()
	if err != nil {
		commons.Logger.Error("create stockTradeStatus failed:", err.Error())
	}
	w.kChartPool, err = tunny.CreatePool(int(g.Config().ThriftServer.KChartPoolNum), w.workHandler.QueryKChart).Open()
	if err != nil {
		commons.Logger.Error("create kChartPool failed:", err.Error())
	}
	w.timeChartPool, err = tunny.CreatePool(int(g.Config().ThriftServer.TimeChartPoolNum), w.workHandler.QueryTimeChart).Open()
	if err != nil {
		commons.Logger.Error("create TimeChartPool failed:", err.Error())
	}
	w.stockBasicInfoBatchPool, err = tunny.CreatePool(int(g.Config().ThriftServer.StockBasicInfoBatchPoolNum), w.workHandler.QueryStockBasicInfoBatch).Open()
	if err != nil {
		commons.Logger.Error("create stockBasicInfoBatchPool failed:", err.Error())
	}
	// add by 2.3.8
	w.stockTicksPool, err = tunny.CreatePool(int(g.Config().ThriftServer.StockTicksPoolNum), w.workHandler.QueryStockTicks).Open()
	if err != nil {
		commons.Logger.Error("create stockTicksPool failed:", err.Error())
	}

}

// 有超时的workpool
func PostWork(in interface{}, wp *tunny.WorkPool, milliTimeout time.Duration) (interface{}, error) {
	out, err := wp.SendWorkTimed(milliTimeout, in)
	if err != nil {
		commons.Logger.Error("PostWork failed!")
	}
	return out, err
}

// 把接口请求，推到创建的pool中
func (w *WorkerManager) QueryFeedLite(req feed_svr.GetFeedLiteRequest) (*feed_svr.GetFeedLiteResponse, error) {
	out, err := PostWork(req, w.feedLitePool, time.Duration(g.Config().ThriftServer.FeedThriftTimeOut)*time.Millisecond)
	if err != nil {
		commons.Logger.Error("QueryFeedLite post to work failed! ", err.Error())
	}
	return out.(*feed_svr.GetFeedLiteResponse), nil
}

func (w *WorkerManager) QueryFeedEx(req feed_svr.GetFeedExRequest) (*feed_svr.GetFeedExResponse, error) {
	out, err := PostWork(req, w.feedExPool, time.Duration(g.Config().ThriftServer.FeedExTimeOut)*time.Millisecond)
	if err != nil {
		commons.Logger.Error("QueryFeedLite post to work failed! ", err.Error())
	}
	return out.(*feed_svr.GetFeedExResponse), nil
}

//市场深度
func (w *WorkerManager) QueryMarketDepth(stockTiny common.StockTiny) (*feed_svr.GetMarketDepthResponse, error) {
	out, err := PostWork(stockTiny, w.marketDepthPool, time.Duration(g.Config().ThriftServer.MarketDepthTimeOut)*time.Millisecond)
	if err != nil {
		commons.Logger.Error("QueryFeedLite post to work failed! ", err.Error())
	}
	return out.(*feed_svr.GetMarketDepthResponse), nil
}

func (w *WorkerManager) QueryStocksByType(req feed_svr.QueryExchangeStocksReq) (*feed_svr.QueryExchangeStocksResp, error) {
	out, err := PostWork(req, w.stocksByTypePool, time.Duration(g.Config().ThriftServer.StocksByTypeTimeOut)*time.Millisecond)
	if err != nil {
		commons.Logger.Error("QueryStocksByType post to work failed! ", err.Error())
	}
	return out.(*feed_svr.QueryExchangeStocksResp), nil
}

//获取市场交易状态
func (w *WorkerManager) QueryGetStockTradeStatus(req feed_svr.GetTradeStatusReq) (*feed_svr.GetTradeStatusResp, error) {
	out, err := PostWork(req, w.stockTradeStatusPool, time.Duration(g.Config().ThriftServer.StockTradeStatusTimeOut)*time.Millisecond)
	if err != nil {
		commons.Logger.Error("QueryStockStatus post to work failed! ", err.Error())
	}
	return out.(*feed_svr.GetTradeStatusResp), nil
}

func (w *WorkerManager) QueryKChart(req feed_svr.GetKChartRequest) (*feed_svr.GetKChartResponse, error) {
	out, err := PostWork(req, w.kChartPool, time.Duration(g.Config().ThriftServer.KChartTimeOut)*time.Millisecond)
	if err != nil {
		commons.Logger.Error("QueryKChart post to work failed! ", err.Error())
	}
	return out.(*feed_svr.GetKChartResponse), nil
}

func (w *WorkerManager) QueryTimeChart(req feed_svr.GetTimeChartRequest) (*feed_svr.GetTimeChartResponse, error) {
	out, err := PostWork(req, w.timeChartPool, time.Duration(g.Config().ThriftServer.TimeChartTimeOut)*time.Millisecond)
	if err != nil {
		commons.Logger.Error("QueryTimeChart post to work failed! ", err.Error())
	}
	return out.(*feed_svr.GetTimeChartResponse), nil
}

func (w *WorkerManager) QueryStockBasicInfoBatch(req feed_svr.GetStockBasicInfoBatchReq) (*feed_svr.GetStockBasicInfoBatchResp, error) {
	out, err := PostWork(req, w.stockBasicInfoBatchPool, time.Duration(g.Config().ThriftServer.StockBasicInfoBatchTimeOut)*time.Millisecond)
	if err != nil {
		commons.Logger.Error("QueryStockBasicInfoBatch post to work failed! ", err.Error())
	}
	return out.(*feed_svr.GetStockBasicInfoBatchResp), nil
}

// add by 2.3.8
func (w *WorkerManager) QueryStockTicks(req feed_svr.GetStockTicksReq) (*feed_svr.GetStockTicksResp, error) {
	out, err := PostWork(req, w.stockTicksPool, time.Duration(g.Config().ThriftServer.StockTicksTimeOut)*time.Millisecond)
	if err != nil {
		commons.Logger.Error("QueryStockTicks post to work failed! ", err.Error())
	}
	return out.(*feed_svr.GetStockTicksResp), nil
}
