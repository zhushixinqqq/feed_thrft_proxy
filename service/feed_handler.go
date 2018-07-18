// Copyright [2013-2017] <xxxxx.Inc>
//
// Author: zhushixin
// Logical layer main function,deal with detail information through the rpcx interface
// return to FeedHandler layer

package service

import (
	commons "common"
	feedThriftCommon "modules/feed_thrft_proxy/common"
	"modules/feed_thrft_proxy/g"
	"modules/feed_thrft_proxy/rpc"
	"protocols"
	common "protocols/thrift_src/gen-go/common"
	feed_svr "protocols/thrift_src/gen-go/feed_svr"
	"runtime/debug"
	"strings"

	"github.com/luci/go-render/render"
)

type FeedHandler struct {
	rpcManager *rpc.RpcManager
}

func NewFeedHandler() *FeedHandler {
	return &FeedHandler{}
}

func (p *FeedHandler) Init() {
	commons.Logger.Info("feedHandler init...................")
	p.rpcManager = rpc.NewRpcManager()
	p.rpcManager.Init()
}
func CheckPanic() {
	if r := recover(); r != nil {
		commons.Logger.Error("catch an panic %v", r)
		debug.PrintStack()
	}
}

// get the feedlite date by the rpc interface
func (p *FeedHandler) GetFeedLite(request feed_svr.GetFeedLiteRequest) (resp *feed_svr.GetFeedLiteResponse, err error) {
	defer CheckPanic()
	queryFeedLiteResp := &protocols.QueryFeedLiteResp{}
	queryFeedLiteReq := protocols.QueryFeedLiteReq{}
	var requestId string
	commons.Logger.Debug("reuqest:", render.Render(request))
	if request.AccessType != nil {
		if request.AccessType.RequestId != nil {
			requestId = *request.AccessType.RequestId
		} else {
			requestId = commons.GenerateUniqueID()
		}
	} else {
		requestId = commons.GenerateUniqueID()
	}
	accessInfo := &protocols.AccessInfo{
		AccessType: g.Config().ThriftServer.ServiceName,
		RequestId:  requestId,
	}
	queryFeedLiteReq.AccessInfo = accessInfo

	//把thrift协议请求的request转化为protocols协议的queryFeedLiteReq
	for _, reqStockTiny := range request.StockTiny {
		var stockTiny protocols.StockTiny
		stockTiny.StockId = reqStockTiny.StockId
		stockTiny.StockType = feedThriftCommon.ThriftType2ProtType(reqStockTiny.TypeA1)
		queryFeedLiteReq.StockTinys = append(queryFeedLiteReq.StockTinys, &stockTiny)
	}
	commons.Logger.Debug("queryFeedLiteReq from rpcManager:", render.Render(queryFeedLiteReq))
	feedSvrFeedLiteResp := &feed_svr.GetFeedLiteResponse{}
	// // notice get the response from the rpcManager
	err = p.rpcManager.GetFeedLite(queryFeedLiteReq, queryFeedLiteResp)
	if err != nil {
		commons.Logger.Error("p.rpcManager.GetFeedLit faild err is:", err)
		feedSvrFeedLiteResp.Err = common.ForbagErrno_ERR_STOCK_ID_NOT_FOUND
		return feedSvrFeedLiteResp, err
	}
	commons.Logger.Debug("getfeedlite.realresp from rpcManager:", render.Render(queryFeedLiteResp))
	// feedSvrFeedLiteResp.Err = feedThriftCommon.ThriftErr(queryFeedLiteResp.ErrInfo)
	feedSvrFeedLiteResp.Err = feedThriftCommon.ThriftErr(queryFeedLiteResp.ErrInfo)
	var feedEx []*common.FeedLite
	priceTimestamp := feedThriftCommon.UnixTimeStamp()
	for feedsStockTiny, feedsFeedLite := range queryFeedLiteResp.Feeds {
		stockTiny := feedsStockTiny
		isSuspended := feedThriftCommon.IsSuspended(feedsFeedLite.RealFeed.TradeStatus)
		errCode := feedThriftCommon.ThriftErr(feedsFeedLite.RealFeed.ErrInfo)
		stockType := feedThriftCommon.ProtType2ThriftType(stockTiny.StockType)
		stockType = feedThriftCommon.SecialStockTypeConvert(stockTiny.StockId, stockType)
		feedLite := &common.FeedLite{
			Err:                errCode,
			StockId:            &(stockTiny.StockId),
			StockType:          &stockType,
			CurPrice:           &(feedsFeedLite.RealFeed.LastPrice),
			OpenPriceToday:     &(feedsFeedLite.RealFeed.OpenPrice),
			UpsAndDowns:        feedsFeedLite.RealFeed.PriceChange,
			UpsAndDownsPercent: feedsFeedLite.RealFeed.PriceChangeRate,
			HighPrice:          &(feedsFeedLite.RealFeed.HighPrice),
			LowPrice:           &(feedsFeedLite.RealFeed.LowPrice),
			LastClosePrice:     feedsFeedLite.RealFeed.PreclosePrice,
			StockName:          &(feedsFeedLite.RealFeed.StockName),
			MinQuotedPrice:     feedsFeedLite.MinQuotedPrice,
			IsTrade:            &(feedsFeedLite.IsTrade),
			ChiSpellingGrp:     &(feedsFeedLite.RealFeed.ChiSpelling),
			MarketValue:        feedsFeedLite.RealFeed.MarketValue,
			MinTradeUnit:       feedsFeedLite.RealFeed.SharesPerHand,
			IsSuspended:        &isSuspended,
			PriceTimestamp:     &priceTimestamp,
		}
		feedEx = append(feedEx, feedLite)
	}
	feedSvrFeedLiteResp.FeedEx = feedEx
	commons.Logger.Debug("feedSvrFeedLiteResp:", render.Render(feedSvrFeedLiteResp))
	return feedSvrFeedLiteResp, err
}

func (p *FeedHandler) GetFeedEx(request feed_svr.GetFeedExRequest) (*feed_svr.GetFeedExResponse, error) {
	defer CheckPanic()
	queryFeedExResp := &protocols.QueryFeedResp{}
	queryFeedExReq := protocols.QueryFeedReq{}

	var requestId string
	commons.Logger.Debug("reuqest:", render.Render(request))
	if request.AccessType != nil {
		if request.AccessType.RequestId != nil {
			requestId = *request.AccessType.RequestId
		} else {
			requestId = commons.GenerateUniqueID()
		}
	} else {
		requestId = commons.GenerateUniqueID()
	}
	accessInfo := &protocols.AccessInfo{
		AccessType: g.Config().ThriftServer.ServiceName,
		RequestId:  requestId,
	}

	queryFeedExReq.AccessInfo = accessInfo

	//把thrift协议请求的request转化为protocols协议的queryFeedReq
	for _, reqStockTiny := range request.StockTiny {
		var stockTiny protocols.StockTiny
		stockTiny.StockId = reqStockTiny.StockId
		stockTiny.StockType = feedThriftCommon.ThriftType2ProtType(reqStockTiny.TypeA1)
		queryFeedExReq.StockTinys = append(queryFeedExReq.StockTinys, &stockTiny)
	}
	commons.Logger.Debug("queryFeedExReq from rpcManager:", render.Render(queryFeedExReq))

	feedSvrFeedExResp := &feed_svr.GetFeedExResponse{}
	err := p.rpcManager.GetFeedEx(queryFeedExReq, queryFeedExResp)
	if err != nil {
		commons.Logger.Error("p.rpcManager.GetFeedEx faild err is:", err)
		return feedSvrFeedExResp, err
	}
	commons.Logger.Debug("get queryFeedExResp from rpcManager:", render.Render(queryFeedExResp))
	var feedExs []*common.FeedEx
	timeStamp := feedThriftCommon.UnixTimeStamp()
	for feedsStockTiny, feedsFeedEx := range queryFeedExResp.Feeds {
		stockTiny := feedsStockTiny
		isSuspended := feedThriftCommon.IsSuspended(feedsFeedEx.RealFeed.TradeStatus)
		errCode := feedThriftCommon.ThriftErr(feedsFeedEx.RealFeed.ErrInfo)
		delayPrice := int64(feedsFeedEx.DelayPrice)
		turnoverVolume := float64(feedsFeedEx.RealFeed.BusinessAmount)
		// 字符串截取 2017-10-10 13:00:00 -> 10-10 13:00:00
		indexSubStr := strings.IndexAny(feedsFeedEx.RealPriceTime, "-")
		realPriceTime := ""
		if indexSubStr != -1 {
			realPriceTime = feedsFeedEx.RealPriceTime[indexSubStr+1:]
		}
		amplitude := feedsFeedEx.RealFeed.Amplitude / 100
		turnoverRate := feedsFeedEx.RealFeed.TurnoverRatio / 100
		stockType := feedThriftCommon.ProtType2ThriftType(stockTiny.StockType)
		stockType = feedThriftCommon.SecialStockTypeConvert(stockTiny.StockId, stockType)

		feedEx := &common.FeedEx{
			Feed: &common.Feed{
				Err:                   errCode,
				CurPrice:              feedsFeedEx.RealFeed.LastPrice,
				MarketValue:           feedsFeedEx.RealFeed.MarketValue,
				PE:                    feedsFeedEx.RealFeed.PERate,
				StockId:               stockTiny.StockId,
				StockType:             stockType,
				OpenPriceToday:        feedsFeedEx.RealFeed.OpenPrice,
				UpsAndDowns:           feedsFeedEx.RealFeed.PriceChange,
				UpsAndDownsPercent:    feedsFeedEx.RealFeed.PriceChangeRate,
				LastClosePrice:        feedsFeedEx.RealFeed.PreclosePrice,
				MinTradeUnit:          feedsFeedEx.RealFeed.SharesPerHand,
				TurnoverVolume:        &turnoverVolume,
				TurnoverValue:         &feedsFeedEx.RealFeed.BusinessBalance,
				Amplitude:             &amplitude,
				HighPrice:             &feedsFeedEx.RealFeed.HighPrice,
				LowPrice:              &feedsFeedEx.RealFeed.LowPrice,
				FiftyTwoWeekHighPrice: &feedsFeedEx.RealFeed.Week52HighPrice,
				FiftyTwoWeekLowPrice:  &feedsFeedEx.RealFeed.Week52LowPrice,
				StockName:             &feedsFeedEx.RealFeed.StockName,
				TurnoverRate:          &turnoverRate,
				MinQuotedPrice:        &(feedsFeedEx.MinQuotedPrice),
				IsTrade:               &feedsFeedEx.IsTrade,
				IsBuy:                 &feedsFeedEx.IsBuy,
				DelayPrice:            &delayPrice,
				RealPriceTime:         &realPriceTime,
				IsSuspended:           &isSuspended,
				Timestamp:             &timeStamp,
				IsDelist:              &feedsFeedEx.IsDelist,
			},
			Stockex: &common.StockEx{},
		}
		feedExs = append(feedExs, feedEx)
	}
	feedSvrFeedExResp.FeedEx = feedExs
	commons.Logger.Debug("get queryFeedExResp from rpcManager:", render.Render(feedSvrFeedExResp))
	return feedSvrFeedExResp, nil
}

func (p *FeedHandler) GetMarketDepth(reqStockTiny common.StockTiny) (resp *feed_svr.GetMarketDepthResponse, err error) {
	defer CheckPanic()
	queryMarketDepthReq := protocols.QueryMarketDepthReq{}
	queryMarketDepthResp := &protocols.QueryMarketDepthResp{}
	requestId := commons.GenerateUniqueID()
	accessInfo := &protocols.AccessInfo{
		AccessType: g.Config().ThriftServer.ServiceName,
		RequestId:  requestId,
	}
	queryMarketDepthReq.AccessInfo = accessInfo
	var stockTiny protocols.StockTiny
	stockTiny.StockId = reqStockTiny.StockId
	stockTiny.StockType = feedThriftCommon.ThriftType2ProtType(reqStockTiny.TypeA1)
	queryMarketDepthReq.StockTinys = append(queryMarketDepthReq.StockTinys, &stockTiny)

	feedSvrMarketDepthResp := &feed_svr.GetMarketDepthResponse{}
	err = p.rpcManager.GetMarketDepth(queryMarketDepthReq, queryMarketDepthResp)
	if err != nil {
		commons.Logger.Error("p.rpcManager.GetMarketDepth faild err is:", err)
		return feedSvrMarketDepthResp, err
	}
	commons.Logger.Debug("get queryMarketDepthResp from rpcManager:", render.Render(queryMarketDepthResp))
	var marketDepths []*common.MarketDepthTuple
	for _, marketDepth := range queryMarketDepthResp.MarketDepthes {
		if len(marketDepth.BuyerMarketDepthLevels) == 10 {
			sellMarketDepths := feedThriftCommon.GetSellMarketDepths(marketDepth, reqStockTiny.TypeA1)
			//TODO@2
			BuyMarketDepths := feedThriftCommon.GetBuyMarketDepths(marketDepth, reqStockTiny.TypeA1)
			commons.Logger.Debug("sellMarketDepths----:", sellMarketDepths,
				" BuyMarketDepths----:", BuyMarketDepths)
			for i, _ := range sellMarketDepths {
				count := sellMarketDepths[i].Count
				buyCount := BuyMarketDepths[i].Count
				marketDepthTuple := &common.MarketDepthTuple{
					SellPrice: sellMarketDepths[i].Price,
					SellQty:   int32(sellMarketDepths[i].Volume),
					SellCount: &(count),
					BuyPrice:  BuyMarketDepths[i].Price,
					BuyQty:    int32(BuyMarketDepths[i].Volume),
					BuyCount:  &buyCount,
				}
				marketDepths = append(marketDepths, marketDepthTuple)
			}
		} else {
			for i := 0; i < len(marketDepth.BuyerMarketDepthLevels); i++ {
				marketDepthTuple := &common.MarketDepthTuple{
					SellPrice: marketDepth.SellerMarketDepthLevels[i].Price,
					SellQty:   int32(marketDepth.SellerMarketDepthLevels[i].Volume),
					SellCount: &(marketDepth.SellerMarketDepthLevels[i].Count),
					BuyPrice:  marketDepth.BuyerMarketDepthLevels[i].Price,
					BuyQty:    int32(marketDepth.BuyerMarketDepthLevels[i].Volume),
					BuyCount:  &(marketDepth.BuyerMarketDepthLevels[i].Count),
				}
				marketDepths = append(marketDepths, marketDepthTuple)
			}
		}
	}
	//TODO 从feedlite接口中获取昨收价
	lastClosePrice := GetLastClosePrice(stockTiny, p)
	feedSvrMarketDepthResp.LastClosePrice = &lastClosePrice
	errCode := feedThriftCommon.ThriftErr(queryMarketDepthResp.ErrInfo)
	feedSvrMarketDepthResp.Err = errCode
	feedSvrMarketDepthResp.MarketDepth = marketDepths
	return feedSvrMarketDepthResp, nil
}

//用于获取股票昨收价
func GetLastClosePrice(stockTiny protocols.StockTiny, p *FeedHandler) float64 {
	defer CheckPanic()
	queryFeedLiteResp := &protocols.QueryFeedLiteResp{}
	queryFeedLiteReq := protocols.QueryFeedLiteReq{}
	requestId := commons.GenerateUniqueID()
	accessInfo := &protocols.AccessInfo{
		AccessType: g.Config().ThriftServer.ServiceName,
		RequestId:  requestId,
	}
	lastClosePrice := 0.0
	queryFeedLiteReq.AccessInfo = accessInfo
	queryFeedLiteReq.StockTinys = append(queryFeedLiteReq.StockTinys, &stockTiny)
	commons.Logger.Debug("GetLastClosePrice queryFeedLiteReq:", render.Render(queryFeedLiteReq))
	err := p.rpcManager.GetFeedLite(queryFeedLiteReq, queryFeedLiteResp)
	if err != nil {
		commons.Logger.Error("p.rpcManager.GetFeedLit faild err is:", err)
		return lastClosePrice
	}
	commons.Logger.Debug("getfeedlite.realresp from rpcManager:", render.Render(queryFeedLiteResp))
	for _, feedsFeedLite := range queryFeedLiteResp.Feeds {
		lastClosePrice = feedsFeedLite.RealFeed.PreclosePrice
	}
	return lastClosePrice
}

func (p *FeedHandler) GetStocksByType(request feed_svr.QueryExchangeStocksReq) (*feed_svr.QueryExchangeStocksResp, error) {
	defer CheckPanic()
	queryExchangeStocksReq := protocols.QueryExchangeStocksReq{}
	queryExchangeStocksResp := &protocols.QueryExchangeStocksResp{}
	if request.AccessInfo != nil {
		accessInfo := &protocols.AccessInfo{
			AccessType: g.Config().ThriftServer.ServiceName,
			RequestId:  *request.AccessInfo.RequestId,
		}
		queryExchangeStocksReq.AccessInfo = accessInfo
	}
	queryExchangeStocksReq.StockType = feedThriftCommon.ThriftType2ProtType(request.StockType)
	commons.Logger.Debug("queryExchangeStocksReq from rpcManager:", render.Render(queryExchangeStocksReq))
	feedSvrQueryExchangeStocksResp := &feed_svr.QueryExchangeStocksResp{}
	err := p.rpcManager.GetStocksByType(queryExchangeStocksReq, queryExchangeStocksResp)
	if err != nil {
		commons.Logger.Error("p.rpcManager.GetStocksByType faild err is:", err)
		return feedSvrQueryExchangeStocksResp, err
	}
	commons.Logger.Debug("GetStocksByType from rpcManager:", render.Render(queryExchangeStocksResp))
	errCode := feedThriftCommon.ThriftErr(queryExchangeStocksResp.ErrInfo)
	feedSvrQueryExchangeStocksResp.ErrInfo = errCode
	var stocks []*common.StockTiny
	for _, protStockTiny := range queryExchangeStocksResp.StockTinys {
		stockType := protStockTiny.StockType
		stockTiny := &common.StockTiny{
			StockId: protStockTiny.StockId,
			TypeA1:  feedThriftCommon.ProtType2ThriftType(stockType),
		}
		stocks = append(stocks, stockTiny)
	}
	feedSvrQueryExchangeStocksResp.Stocks = stocks
	return feedSvrQueryExchangeStocksResp, err
}

//
func (p *FeedHandler) GetKChart(request feed_svr.GetKChartRequest) (*feed_svr.GetKChartResponse, error) {
	defer CheckPanic()
	commons.Logger.Debug("GetKChart from rpcManager:", render.Render(request))
	queryKChartResp := &protocols.QueryKChartResp{}
	var requestId string
	if request.AccessInfo != nil {
		requestId = *request.AccessInfo.RequestId
	} else {
		requestId = commons.GenerateUniqueID()
	}
	accessInfo := &protocols.AccessInfo{
		AccessType: g.Config().ThriftServer.ServiceName,
		RequestId:  requestId,
	}
	stockTiny := &protocols.StockTiny{
		StockId:   request.Stock.StockId,
		StockType: feedThriftCommon.ThriftType2ProtType(request.Stock.TypeA1),
	}
	kChartType := feedThriftCommon.KChartType2KLineType(request.KChartType)
	queryKChartReq := protocols.QueryKChartReq{
		AccessInfo: accessInfo,
		StockTiny:  stockTiny,
		KChartType: kChartType,
		StartDate:  request.StartDate,
		EndDate:    request.EndDate,
		GetType:    "range", //TODO 暂时用range
	}
	commons.Logger.Debug("queryKChartReq from rpcManager:", render.Render(queryKChartReq))
	err := p.rpcManager.GetKChart(queryKChartReq, queryKChartResp)
	commons.Logger.Debug("GetKChart from rpcManager:", render.Render(queryKChartResp))
	errCode := feedThriftCommon.ThriftErr(queryKChartResp.ErrInfo)
	feedSvrQueryKChartResp := &feed_svr.GetKChartResponse{
		Err:   &errCode,
		Inner: &common.KChart{},
	}
	if err != nil || queryKChartResp.KLines == nil {
		feedSvrQueryKChartResp.IsOk = &queryKChartResp.IsOk
		commons.Logger.Error("p.rpcManager.GetKChart faild err is:", err)
		return feedSvrQueryKChartResp, err
	}
	innerKChart := &common.KChart{}
	// stockType := queryKChartResp.StockTiny.StockType
	feedSvrStockTiny := &common.StockTiny{}
	if queryKChartResp.StockTiny != nil {
		feedSvrStockTiny.StockId = queryKChartResp.StockTiny.StockId
		feedSvrStockTiny.TypeA1 = feedThriftCommon.ProtType2ThriftType(queryKChartResp.StockTiny.StockType)
	}
	for _, kLine := range queryKChartResp.KLines {
		turnoverVolume := float64(kLine.KPoint.BusinessAmount)
		turnoverVolumeType := (common.TurnoverVolumeType)(kLine.TurnoverVolumeType)
		kPoint := &common.KPoint{
			StockTiny:  feedSvrStockTiny,
			OpenPrice:  kLine.KPoint.OpenPrice,
			ClosePrice: kLine.KPoint.ClosePrice,
			HighPrice:  kLine.KPoint.HighPrice,
			LowPrice:   kLine.KPoint.LowPrice,

			RiseFallPrice:      kLine.ChangePrice,
			RiseFallPercent:    kLine.ChangePercent,
			TurnoverVolume:     turnoverVolume,
			TurnoverValue:      &kLine.KPoint.BusinessBalance,
			TurnoverRate:       kLine.KPoint.TurnoverRate / 10000, //hs数据源放大了100倍，去掉%
			Ma5:                kLine.Ma5,
			Ma10:               kLine.Ma10,
			Ma20:               kLine.Ma20,
			Ma30:               kLine.Ma30,
			TimestampMs:        kLine.TimestampInS,
			Date:               kLine.DateTime,
			PrevClosePrice:     kLine.KPoint.PreClosePrice,
			TurnoverVolumeType: turnoverVolumeType,
		}
		innerKChart.KChart = append(innerKChart.KChart, kPoint)
	}
	feedSvrQueryKChartResp.Inner = innerKChart
	feedSvrQueryKChartResp.IsOk = &queryKChartResp.IsOk
	return feedSvrQueryKChartResp, nil
}

func (p *FeedHandler) GetTimeChart(request feed_svr.GetTimeChartRequest) (*feed_svr.GetTimeChartResponse, error) {
	defer CheckPanic()
	commons.Logger.Debug("GetTimeChart from rpcManager:", render.Render(request))
	queryTimeChartResp := &protocols.QueryTimeChartResp{}
	// 把thrift协议请求的字段转为protocols协议请求的字段
	var requestId string
	if request.AccessInfo != nil {
		requestId = *request.AccessInfo.RequestId
	} else {
		requestId = commons.GenerateUniqueID()
	}
	accessInfo := &protocols.AccessInfo{
		AccessType: g.Config().ThriftServer.ServiceName,
		RequestId:  requestId,
	}
	stockTiny := &protocols.StockTiny{
		StockId:   request.Stock.StockId,
		StockType: feedThriftCommon.ThriftType2ProtType(request.Stock.TypeA1),
	}
	timeChartPeriod := feedThriftCommon.ThriftPeriod2ProtPeriod(request.Period)
	queryTimeChartReq := protocols.QueryTimeChartReq{
		AccessInfo: accessInfo,
		StockTiny:  stockTiny,
		Period:     timeChartPeriod,
	}
	commons.Logger.Debug("queryTimeChartReq from rpcManager:", render.Render(queryTimeChartReq))
	err := p.rpcManager.GetTimeChart(queryTimeChartReq, queryTimeChartResp)
	commons.Logger.Debug("GetTimeChart from rpcManager:", render.Render(queryTimeChartResp))
	errCode := feedThriftCommon.ThriftErr(queryTimeChartResp.ErrInfo)
	// 把rpcx请求的resp转为feed_svr的GetTimeChartResponse
	feedSvrQueryTimeChartResp := &feed_svr.GetTimeChartResponse{
		Err:   &errCode,
		Inner: &common.TimeChart{},
	}
	if err != nil || queryTimeChartResp.TimeChart == nil {
		commons.Logger.Error("p.rpcManager.GetTimeChart faild err is:", err)
		return feedSvrQueryTimeChartResp, err
	}
	timeCharts := &common.TimeChart{}
	for _, timeChart := range queryTimeChartResp.TimeChart.TimeChartPoints {
		turnoverVolume := int64(timeChart.TimePoint.BusinessAmount)
		turnoverVolumeType := (common.TurnoverVolumeType)(timeChart.TurnoverVolumeType)
		dateTime := feedThriftCommon.Timestamp2String(timeChart.TimestampInS)
		timePoint := &common.TimePoint{
			Price:              timeChart.TimePoint.LatestPrice,
			TurnoverVolume:     turnoverVolume,
			TimestampMs:        timeChart.TimestampInS,
			TurnoverVolumeType: &turnoverVolumeType,
			DateTime:           &dateTime,
			AvgPrice:           &timeChart.TimePoint.AvgPrice,
			TurnoverValue:      &timeChart.TimePoint.BusinessBalance,
		}
		timeCharts.TimeChart = append(timeCharts.TimeChart, timePoint)
	}
	timeCharts.LastClosePrice = &queryTimeChartResp.TimeChart.LastClosePrice
	timeCharts.AmplitudePrice = &queryTimeChartResp.TimeChart.AmplitudePrice
	for _, fivedayChartDay := range queryTimeChartResp.TimeChart.FivedayChartDays {
		timeCharts.FivedayChartDays = append(timeCharts.FivedayChartDays, fivedayChartDay)
	}
	for _, fivedayStartIdx := range queryTimeChartResp.TimeChart.FivedayStartIdxes {
		timeCharts.FivedayStartIdx = append(timeCharts.FivedayStartIdx, fivedayStartIdx)
	}
	for _, xAxisLabel := range queryTimeChartResp.TimeChart.XAxisLabels {
		timeCharts.XAxisLabels = append(timeCharts.XAxisLabels, xAxisLabel)
	}
	for _, xAxisLabelIdx := range queryTimeChartResp.TimeChart.XAxisLabelIdxes {
		timeCharts.XAxisLabelIdx = append(timeCharts.XAxisLabelIdx, xAxisLabelIdx)
	}
	feedSvrQueryTimeChartResp.Inner = timeCharts
	feedSvrQueryTimeChartResp.TotalCountRange = &queryTimeChartResp.TotalCountRange
	feedSvrQueryTimeChartResp.IsOk = &queryTimeChartResp.IsOk
	return feedSvrQueryTimeChartResp, err
}

func (p *FeedHandler) GetStockTypeStatus(request feed_svr.GetTradeStatusReq) (*feed_svr.GetTradeStatusResp, error) {
	defer CheckPanic()
	accessInfo := &protocols.AccessInfo{
		AccessType: g.Config().ThriftServer.ServiceName,
		RequestId:  commons.GenerateUniqueID(),
	}
	queryStockTypeTradeStatusReq := protocols.GetStockTypeTradeStatusReq{
		AccessInfo: accessInfo,
		StockType:  feedThriftCommon.ThriftType2ProtType(request.StockType),
	}
	commons.Logger.Debug("queryStockTypeTradeStatusReq from rpcManager:", render.Render(queryStockTypeTradeStatusReq))
	queryStockTypeTradeStatusResp := &protocols.GetStockTypeTradeStatusResp{}
	err := p.rpcManager.GetStockTypeStatus(queryStockTypeTradeStatusReq, queryStockTypeTradeStatusResp)
	commons.Logger.Debug("queryStockTypeTradeStatusResp from rpcManager:", render.Render(queryStockTypeTradeStatusResp))
	feedSvrTradeStatusResp := &feed_svr.GetTradeStatusResp{}
	if err != nil {
		commons.Logger.Error("p.rpcManager.GetTimeChart faild err is:", err)
		return feedSvrTradeStatusResp, err
	}
	stockTypeTradeStatus := feedThriftCommon.TypeTradeStatusConvert(queryStockTypeTradeStatusResp.StockTypeTradeStatus)
	feedSvrTradeStatusResp.StockStatus = stockTypeTradeStatus
	return feedSvrTradeStatusResp, err
}

func (p *FeedHandler) GetStockBasicInfoBatch(request feed_svr.GetStockBasicInfoBatchReq) (*feed_svr.GetStockBasicInfoBatchResp, error) {
	defer CheckPanic()
	queryStockBasicInfoBatchResp := &protocols.GetStockBasicInfoBatchResp{}
	// 把thrift协议请求的字段转为protocols协议请求的字段
	var requestId string
	if request.AccessInfo != nil && request.AccessInfo.RequestId != nil {
		if request.AccessInfo.RequestId == nil {
			requestId = *request.AccessInfo.RequestId
		} else {
			requestId = commons.GenerateUniqueID()
		}
	} else {
		requestId = commons.GenerateUniqueID()
	}
	accessInfo := &protocols.AccessInfo{
		AccessType: g.Config().ThriftServer.ServiceName,
		RequestId:  requestId,
	}
	var stockTinys []*protocols.StockTiny
	for _, reqStockTiny := range request.Stocks {
		var stockTiny protocols.StockTiny
		stockTiny.StockId = reqStockTiny.StockId
		stockTiny.StockType = feedThriftCommon.ThriftType2ProtType(reqStockTiny.TypeA1)
		stockTinys = append(stockTinys, &stockTiny)
	}
	queryStockBasicInfoBatchReq := protocols.GetStockBasicInfoBatchReq{
		AccessInfo: accessInfo,
		StockTinys: stockTinys,
	}
	commons.Logger.Debug("queryStockBasicInfoBatchReq from rpcManager:", render.Render(queryStockBasicInfoBatchReq))
	err := p.rpcManager.GetStockBasicInfoBatch(queryStockBasicInfoBatchReq, queryStockBasicInfoBatchResp)
	commons.Logger.Debug("feed_svr queryStockBasicInfoBatchResp from rpcManager:", render.Render(queryStockBasicInfoBatchResp))
	feedSvrQueryStockBasicInfoBatchResp := &feed_svr.GetStockBasicInfoBatchResp{}
	feedSvrQueryStockBasicInfoBatchResp.Err = feedThriftCommon.ThriftErr(queryStockBasicInfoBatchResp.ErrInfo)
	if err != nil {
		commons.Logger.Error("p.rpcManager.GetStockBasicInfoBatch faild err is:", err)
		return feedSvrQueryStockBasicInfoBatchResp, err
	}
	var stockBasicInfos []*feed_svr.StockBasicInfo
	for stockTiny, basicInfo := range queryStockBasicInfoBatchResp.BasicInfos {
		stockType := feedThriftCommon.ProtType2ThriftType(stockTiny.StockType)
		stockTiny := &common.StockTiny{
			StockId: stockTiny.StockId,
			TypeA1:  stockType,
		}
		stockBasicInfo := &feed_svr.StockBasicInfo{
			StockTiny:                stockTiny,
			StockNameChnSimplified:   &basicInfo.StockNameChnSimp,
			StockNameEng:             &basicInfo.StockNameEng,
			StockNameChnUnsimplified: &basicInfo.StockNameChnUnSimp,
		}
		stockBasicInfos = append(stockBasicInfos, stockBasicInfo)
	}
	feedSvrQueryStockBasicInfoBatchResp.Stocks = stockBasicInfos

	return feedSvrQueryStockBasicInfoBatchResp, err
}

func (p *FeedHandler) GetStockTicks(request feed_svr.GetStockTicksReq) (*feed_svr.GetStockTicksResp, error) {
	defer CheckPanic()
	//TODO
	feedSvrQueryStockTicksResp := &feed_svr.GetStockTicksResp{}
	// 把thrift协议请求的字段转为protocols协议请求的字段
	var requestId string
	if request.AccessInfo != nil && request.AccessInfo.RequestId != nil {
		if request.AccessInfo.RequestId == nil {
			requestId = *request.AccessInfo.RequestId
		} else {
			requestId = commons.GenerateUniqueID()
		}
	} else {
		requestId = commons.GenerateUniqueID()
	}
	accessInfo := &protocols.AccessInfo{
		AccessType: g.Config().ThriftServer.ServiceName,
		RequestId:  requestId,
	}
	var stockTiny protocols.StockTiny
	stockTiny.StockId = request.Stock.StockId
	stockTiny.StockType = feedThriftCommon.ThriftType2ProtType(request.Stock.TypeA1)
	queryStockTickReq := protocols.QueryTickerReq{
		AccessInfo:      accessInfo,
		StockTiny:       &stockTiny,
		StartIndex:      request.StartIndex,
		ReqNum:          request.ReqNum,
		SearchDirection: request.SearchDirection,
	}
	commons.Logger.Debug("queryStockTickReq from rpcManager:", render.Render(queryStockTickReq))
	queryStockTicksResp := &protocols.QueryTickerResp{}
	err := p.rpcManager.GetStockTicks(queryStockTickReq, queryStockTicksResp)
	feedSvrQueryStockTicksResp.Err = feedThriftCommon.ThriftErr(queryStockTicksResp.ErrInfo)
	//TODO
	if err != nil {
		commons.Logger.Error("p.rpcManager.GetStockTicks faild err is:", err)
		return feedSvrQueryStockTicksResp, err
	}
	//

	var stockTicks []*common.StockTick
	for _, Ticks := range queryStockTicksResp.Tickers {
		// Ticks.Direction 0:卖,1:买
		var direction int32
		if Ticks.Direction == 0 {
			direction = 1
		} else if Ticks.Direction == 1 {
			direction = 0
		} else {
			direction = 2 // 未知
		}
		stockTick := &common.StockTick{
			DateTime:       Ticks.BusinessTime,
			Price:          Ticks.Price,
			TurnoverVolume: Ticks.Amount,
			TradeIndex:     int64(Ticks.Index),
			Direction:      direction,
		}
		stockTicks = append(stockTicks, stockTick)
	}
	feedSvrQueryStockTicksResp.StockTicks = stockTicks
	lastClosePrice := GetLastClosePrice(stockTiny, p)
	feedSvrQueryStockTicksResp.LastClosePrice = &lastClosePrice
	//
	commons.Logger.Debug("feed_svr feedSvrQueryStockTicksResp from rpcManager:", render.Render(feedSvrQueryStockTicksResp))
	// errCode := &protocols.ErrInfo{
	// ErrCode: protocols.ErrCode_SUCCESS,
	// ErrMsg:  "SUCCESS TEST"}
	// feedSvrQueryStockTicksResp.Err = feedThriftCommon.ThriftErr(errCode)
	return feedSvrQueryStockTicksResp, nil
}
