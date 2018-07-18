// Copyright [2013-2017] <xxxxx.Inc>
//
// Author: zhushixin
package service

import (
	commons "common"
	common "protocols/thrift_src/gen-go/common"
	feed_svr "protocols/thrift_src/gen-go/feed_svr"

	"github.com/luci/go-render/render"
)

type WorkHandler struct {
	feedHandler *FeedHandler
}

func NewWorkHandler() *WorkHandler {
	return &WorkHandler{}
}

func (h *WorkHandler) Init() {
	commons.Logger.Info("WrokHandler Init ok................")
	h.feedHandler = NewFeedHandler()
	h.feedHandler.Init()
}

func (h *WorkHandler) QueryFeedLite(in interface{}) interface{} {
	req := in.(feed_svr.GetFeedLiteRequest)
	queryFeedLiteResp, err := h.feedHandler.GetFeedLite(req)
	if err != nil {
		commons.Logger.Error("feedHandler.GetFeedLite is faild:", err)
	}
	commons.Logger.Debug("workhandler queryFeedLiteResp", render.Render(queryFeedLiteResp))
	return queryFeedLiteResp
}

// 调用feed_hander接口中的GetFeedEx函数
func (h *WorkHandler) QueryFeedEx(in interface{}) interface{} {
	req := in.(feed_svr.GetFeedExRequest)
	queryFeedExResp, err := h.feedHandler.GetFeedEx(req)
	if err != nil {
		commons.Logger.Error("feedHandler.GetFeedEx is faild:", err)
	}
	commons.Logger.Debug("workhandler queryFeedExResp", render.Render(queryFeedExResp))
	return queryFeedExResp
}

// 调用feed_hander接口中的GetMarketDepth函数
func (h *WorkHandler) QueryMarketDepth(in interface{}) interface{} {
	stockTiny := in.(common.StockTiny)
	queryMarketDepth, err := h.feedHandler.GetMarketDepth(stockTiny)
	if err != nil {
		commons.Logger.Error("feedHandler.GetMarketDepth is faild:", err)
	}
	commons.Logger.Debug("workhandler queryMarketDepth ", render.Render(queryMarketDepth))
	return queryMarketDepth
}

func (h *WorkHandler) QueryStocksByType(in interface{}) interface{} {
	queryExchangeStocksReq := in.(feed_svr.QueryExchangeStocksReq)
	queryExchangeStocksResp, err := h.feedHandler.GetStocksByType(queryExchangeStocksReq)
	if err != nil {
		commons.Logger.Error("feedHandler.GetStocksByType is faild:", err)
	}
	commons.Logger.Debug("workhandler  ", render.Render(queryExchangeStocksResp))
	return queryExchangeStocksResp
}

//获取股票市场交易状态，根据stocktype类型判断
func (h *WorkHandler) QueryStockStatus(in interface{}) interface{} {
	//TODO
	queryGetTradeStatusReq := in.(feed_svr.GetTradeStatusReq)
	queryExchangeStocksResp, err := h.feedHandler.GetStockTypeStatus(queryGetTradeStatusReq)
	if err != nil {
		commons.Logger.Error("feedHandler.GetStocksTypeTradeStatus is faild:", err)
	}
	commons.Logger.Debug("workhandler  ", render.Render(queryExchangeStocksResp))
	return queryExchangeStocksResp
}

func (h *WorkHandler) QueryKChart(in interface{}) interface{} {
	queryKChartReq := in.(feed_svr.GetKChartRequest)
	queryKChartResp, err := h.feedHandler.GetKChart(queryKChartReq)
	if err != nil {
		commons.Logger.Error("feedHandler.GetKChart is faild:", err)
	}
	commons.Logger.Debug("workhandler QueryKchart: ", render.Render(queryKChartResp))
	return queryKChartResp
}

func (h *WorkHandler) QueryTimeChart(in interface{}) interface{} {
	queryTimeChartReq := in.(feed_svr.GetTimeChartRequest)
	queryTimeChartResp, err := h.feedHandler.GetTimeChart(queryTimeChartReq)
	if err != nil {
		commons.Logger.Error("feedHandler.GetTimeChart is faild:", err)
	}
	commons.Logger.Debug("workhandler QueryTimechart: ", render.Render(queryTimeChartResp))
	//TODO
	return queryTimeChartResp
}

func (h *WorkHandler) QueryStockBasicInfoBatch(in interface{}) interface{} {
	queryStockBasicInfoBatchReq := in.(feed_svr.GetStockBasicInfoBatchReq)
	queryStockBasicInfoBatchResp, err := h.feedHandler.GetStockBasicInfoBatch(queryStockBasicInfoBatchReq)
	if err != nil {
		commons.Logger.Error("feedHandler.GetStockBasicInfoBatch is faild:", err)
	}
	commons.Logger.Debug("workhandler QueryStockBasicInfoBatch resp: ", render.Render(queryStockBasicInfoBatchResp))
	return queryStockBasicInfoBatchResp
}

// add by 2.3.8
func (h *WorkHandler) QueryStockTicks(in interface{}) interface{} {
	queryStockTicksReq := in.(feed_svr.GetStockTicksReq)
	queryStockTicksResp, err := h.feedHandler.GetStockTicks(queryStockTicksReq)
	if err != nil {
		commons.Logger.Error("feedHandler.GetStockTicks is faild:", err)
	}
	commons.Logger.Debug("workhandler GetStockTicks resp: ", render.Render(queryStockTicksResp))
	return queryStockTicksResp
}
