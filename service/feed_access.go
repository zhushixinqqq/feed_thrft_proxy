// Copyright [2013-2017] <xxxxx.Inc>
//
// Author: zhushixin

package service

import (
	commons "common"
	"errors"
	common "protocols/thrift_src/gen-go/common"
	feed_svr "protocols/thrift_src/gen-go/feed_svr"

	"github.com/luci/go-render/render"
)

type AccessLayer struct {
	workerManager *WorkerManager
}

func NewAccessLayer() *AccessLayer {
	return &AccessLayer{}
}

func (p *AccessLayer) Init() {
	commons.Logger.Info("AccessLayer Init..............")
	p.workerManager = NewWorkManager()
	p.workerManager.Init()
}

func (p *AccessLayer) Ping() error {
	return nil
}

//Parameters:
//request:StockTiny,AccessType; resp:Err,FeedEx
func (p *AccessLayer) GetFeedLite(request *feed_svr.GetFeedLiteRequest) (resp *feed_svr.GetFeedLiteResponse, err error) {
	defer commons.TimeCostStr("feed access GetFeedLite")()
	commons.Logger.Info("GetFeedLite req", render.Render(request))
	queryFeedLiteResponse := &feed_svr.GetFeedLiteResponse{}
	if request.StockTiny == nil {
		commons.Logger.Error("feedsvr getfeedlite request(stocktiny)is nil", render.Render(request))
		queryFeedLiteResponse.Err = common.ForbagErrno_UNKNOWN_ERR
		return queryFeedLiteResponse, err
	}
	queryFeedLiteResponse, err = p.workerManager.QueryFeedLite(*request)
	commons.Logger.Info("GetFeedLite get the resp from feed_handler:", render.Render(queryFeedLiteResponse))
	if err != nil {
		commons.Logger.Error("get GetFeedLite failed ", err)
		err = errors.New("get GetFeedLite  err may be no available service")
		return queryFeedLiteResponse, err
	}
	return queryFeedLiteResponse, err
}

func (p *AccessLayer) GetFeed() (r map[common.StockType]map[string]*common.Feed, err error) {
	return
}

// Parameters:
//  - _request
func (p *AccessLayer) GetStocksByType(request *feed_svr.QueryExchangeStocksReq) (resp *feed_svr.QueryExchangeStocksResp, err error) {
	commons.Logger.Info("GetStocksByType req", render.Render(request))
	defer commons.TimeCostStr("feed access GetStocksByType")()
	// queryExchangeStocksResponse := &feed_svr.QueryExchangeStocksResp{}
	if request == nil {
		commons.Logger.Error("feedsvr getfeedEx request(stocktiny)is nil", render.Render(request))
		err := errors.New("request is nill check it")
		return resp, err
	}
	resp, err = p.workerManager.QueryStocksByType(*request)
	commons.Logger.Info("QueryStocksByType get the resp from feed_handler:", render.Render(resp))
	if err != nil {
		commons.Logger.Error("get stocks failed ", err)
	}
	return resp, err
}

// Parameters:
// request
func (p *AccessLayer) GetFeedEx(request *feed_svr.GetFeedExRequest) (r *feed_svr.GetFeedExResponse, err error) {
	commons.Logger.Info("GetFeedEx req", render.Render(request))
	defer commons.TimeCostStr("feed access GetFeedEx")()
	queryFeedExResponse := &feed_svr.GetFeedExResponse{}
	if request.StockTiny == nil {
		commons.Logger.Error("feedsvr getfeedEx request(stocktiny)is nil", render.Render(request))
		err = errors.New("request is nill check the StockTiny")
		return queryFeedExResponse, err
	}
	queryFeedExResponse, err = p.workerManager.QueryFeedEx(*request)
	commons.Logger.Info("GetFeedEx get the resp from feed_handler:", render.Render(queryFeedExResponse))
	if err != nil {
		commons.Logger.Error("get GetFeedEx failed ", err)
	}
	return queryFeedExResponse, err
}

// Parameters:
// request
func (p *AccessLayer) GetTimeChart(request *feed_svr.GetTimeChartRequest) (*feed_svr.GetTimeChartResponse, error) {
	commons.Logger.Info("feedaccess GetTimeChart req:", render.Render(request))
	defer commons.TimeCostStr("feed access GetTimeChart")()
	if request == nil {
		commons.Logger.Error("feedsvr GetTimeChartRequest requestis nil", render.Render(request))
		err := errors.New("request nill check the GetTimeChartRequest")
		return nil, err
	}
	queryTimeChartResponse, err := p.workerManager.QueryTimeChart(*request)
	commons.Logger.Info(" get the queryTimeChartResponse from work_manager:", render.Render(queryTimeChartResponse))
	if err != nil {
		commons.Logger.Error("get queryTimeChartResponse failed ", err)
	}
	return queryTimeChartResponse, err
}

// Parameters:
// request
func (p *AccessLayer) GetKChart(request *feed_svr.GetKChartRequest) (*feed_svr.GetKChartResponse, error) {
	commons.Logger.Info("feedaccess GetKChart req:", render.Render(request))
	defer commons.TimeCostStr("feed access GetKChart")()
	if request == nil {
		commons.Logger.Error("feedsvr GetKChartRequest request is nil", render.Render(request))
		err := errors.New("request is nil check the GetKChartRequest")
		return nil, err
	}
	queryKChartResponse, err := p.workerManager.QueryKChart(*request)
	commons.Logger.Info(" get the queryKChartResponse from work_manager:", render.Render(queryKChartResponse))
	if err != nil {
		commons.Logger.Error("get queryKChartResponse failed ", err)
	}
	return queryKChartResponse, nil
}

// Parameters:
//  - Stock
func (p *AccessLayer) GetStockEx(stock *common.StockTiny) (r *common.StockEx, err error) {
	return
}
func (p *AccessLayer) GetSuspendedInfo() (r *common.SuspendedInfo, err error) {
	return
}

// Parameters:
//  - StockType
func (p *AccessLayer) GetSuspendedInfoEx(stock_type common.StockType) (r *common.SuspendedInfo, err error) {
	return
}

// Parameters:
//  - Stocks
func (p *AccessLayer) GetRealTimePrice(stocks []string) (r map[string]*common.StockMulti, err error) {
	return
}

// Parameters:
//  - Stocks
func (p *AccessLayer) GetUSHKRealTimePrice(stocks []*common.StockTiny) (r []*common.StockMulti, err error) {
	return
}

// Parameters:
//  - Stocks
//  - IsNeedExchangeRate
func (p *AccessLayer) GetAllRealTimePrice(stocks []*common.StockTiny, is_need_exchange_rate bool) (r *common.GetAllRealTimePriceResult_, err error) {
	return
}

// Parameters:
//  - Stocks
func (p *AccessLayer) GetStockName(stocks []string) (r map[string]string, err error) {
	return
}

// Parameters:
//  - CurrencyType
func (p *AccessLayer) GetExchangeRate(currency_type []common.CurrencyType) (r []float64, err error) {
	return
}

// Parameters:
//  - From
//  - To
func (p *AccessLayer) GetExchangeRateStandard(from common.CurrencyType, to common.CurrencyType) (r *feed_svr.FundInOutExchangeRate, err error) {
	return
}

// Parameters:
//  - Para
func (p *AccessLayer) Test(para int32) (r int32, err error) {
	return
}

// func (p *AccessLayer) Ping() (err error) {
// return nil
// }

// Parameters:
//  - StockIds
//  - StartDate
//  - EndDate
func (p *AccessLayer) GetEventList(stock_ids []*common.StockTiny, start_date string, end_date string) (r *common.GetEventResult_, err error) {
	return
}
func (p *AccessLayer) GetTodayEventList() (r map[common.StockType]map[string][]*common.Event, err error) {
	return
}

// Parameters:
//  - StockTiny
func (p *AccessLayer) GetMarketDepth(stockTiny *common.StockTiny) (*feed_svr.GetMarketDepthResponse, error) {
	commons.Logger.Info("GetMarketDepth req", render.Render(stockTiny))
	defer commons.TimeCostStr("feed access GetMarketDepth")()
	if stockTiny == nil {
		commons.Logger.Error("feedsvr getfeedEx request(stocktiny)is nil", render.Render(stockTiny))
		err := errors.New("stockTiny nill check the StockTiny")
		return nil, err
	}
	queryMarketDepthResponse, err := p.workerManager.QueryMarketDepth(*stockTiny)
	commons.Logger.Info("QueryMarKetDepth get the resp from work_manager:", render.Render(queryMarketDepthResponse))
	if err != nil {
		commons.Logger.Error("get queryMarketDepthResponse failed ", err)
	}
	return queryMarketDepthResponse, err
}

// Parameters:
//  - StockId
func (p *AccessLayer) DeBugOneStockPdf(stock_id *common.StockTiny) (r *common.StockPdfInfo, err error) {
	return
}

// Parameters:
//  - Stock
func (p *AccessLayer) DeBugCrawlStockInfo(stock *common.StockTiny) (r *common.CrawlInfoItem, err error) {
	return
}

//获取股票市场交易状态，根据stocktype类型判断
// Parameters:
// tradeReq
func (p *AccessLayer) GetStockTypeTradeStatus(tradereq *feed_svr.GetTradeStatusReq) (*feed_svr.GetTradeStatusResp, error) {
	commons.Logger.Info("GetStockTypeTradeStatus req", render.Render(tradereq))
	defer commons.TimeCostStr("feed access GetStockTypeTradeStatus")()
	queryTradeStatusResp := &feed_svr.GetTradeStatusResp{}
	if tradereq == nil {
		commons.Logger.Error("feedsvr GetTradeStatusResp request is nil", render.Render(tradereq))
		err := errors.New("tradereq nill check the tradereq")
		return nil, err
	}
	queryTradeStatusResp, err := p.workerManager.QueryGetStockTradeStatus(*tradereq)
	commons.Logger.Info("queryTradeStatusResp get the resp from work_manager:", render.Render(queryTradeStatusResp))
	if err != nil {
		commons.Logger.Error("get queryTradeStatusResp failed ", err)
	}
	return queryTradeStatusResp, err
}

// Parameters:
//  - Req
func (p *AccessLayer) GetStockBasicInfoBatch(req *feed_svr.GetStockBasicInfoBatchReq) (*feed_svr.GetStockBasicInfoBatchResp, error) {
	commons.Logger.Info("GetStockBasicInfoBatch req", render.Render(req))
	defer commons.TimeCostStr("feed access GetStockBasicInfoBatch")()
	queryStockBasicInfoBatchResp := &feed_svr.GetStockBasicInfoBatchResp{}
	if req == nil {
		commons.Logger.Error("feedsvr GetStockBasicInfoBatch request is nil", render.Render(req))
		err := errors.New("req nill check the GetStockBasicInfoBatchReq")
		return nil, err
	}
	queryStockBasicInfoBatchResp, err := p.workerManager.QueryStockBasicInfoBatch(*req)
	commons.Logger.Info("queryStockBasicInfoBatchResp get the resp from work_manager:", render.Render(queryStockBasicInfoBatchResp))
	if err != nil {
		commons.Logger.Error("get queryStockBasicInfoBatchResp failed ", err)
	}
	return queryStockBasicInfoBatchResp, err
}

// Parameters:
//  - Req
func (p *AccessLayer) GetStockBrokers(req *feed_svr.GetStockBrokersReq) (*feed_svr.GetStockBrokdersResp, error) {
	commons.Logger.Info("GetStockBrokers req", render.Render(req))
	defer commons.TimeCostStr("feed access GetStockBrokers")()
	queryStockBrokersResp := &feed_svr.GetStockBrokdersResp{}
	if req == nil {
		commons.Logger.Error("feedsvr GetStockBrokers request is nil", render.Render(req))
		err := errors.New("req nill check the GetStockBrokersReq")
		return nil, err
	}
	// queryStockBrokersResp, err := p.workerManager.QueryStockBrokers(*req)
	// commons.Logger.Info("queryStockBrokersResp get the resp from work_manager:", render.Render(queryStockBasicInfoBatchResp))
	// if err != nil {
	// commons.Logger.Error("get queryStockBrokersResp failed ", err)
	// }

	return queryStockBrokersResp, nil
}

// Parameters:
//  - Req
func (p *AccessLayer) GetStockTicks(req *feed_svr.GetStockTicksReq) (*feed_svr.GetStockTicksResp, error) {
	commons.Logger.Info("GetStockTicks req", render.Render(req))
	defer commons.TimeCostStr("feed access GetStockTicks")()
	queryStockTicksResp := &feed_svr.GetStockTicksResp{}
	if req == nil {
		commons.Logger.Error("feedsvr GetStockTicks request is nil", render.Render(req))
		err := errors.New("req nill check the GetStockTicksReq")
		return nil, err
	}
	queryStockTicksResp, err := p.workerManager.QueryStockTicks(*req)
	commons.Logger.Info("queryStockTicksResp get the resp from work_manager:", render.Render(queryStockTicksResp))
	if err != nil {
		commons.Logger.Error("get queryStockTicksResp failed ", err)
	}
	return queryStockTicksResp, nil
}

// Parameters:
//  - _request
// func (p *AccessLayer) GetFeedLite(_request *feed_svr.GetFeedLiteRequest) (r *feed_svr.GetFeedLiteResponse, err error) {
// return
// }
