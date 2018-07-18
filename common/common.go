// Copyright [2013-2017] <xxxxx.Inc>
//
// Author: zhushixin

package feedThriftCommon

import (
	commons "common"
	"math"
	"protocols"
	common "protocols/thrift_src/gen-go/common"
	"time"

	"github.com/luci/go-render/render"
)

//thrift 协议stocktype转化为protocols协议 stocktype
func ThriftType2ProtType(stockType common.StockType) protocols.StockType {
	protStockType := (protocols.StockType)(stockType)
	if stockType == common.StockType_US_NASDAQ || stockType == common.StockType_US_NYSE {
		protStockType = protocols.StockType_US
	}
	return protStockType
}

//protocols协议 stocktype转化thrift 协议stocktype
func ProtType2ThriftType(stockType protocols.StockType) common.StockType {
	thriftStockType := (common.StockType)(stockType)
	return thriftStockType
}

func SecialStockTypeConvert(stockId string, stockType common.StockType) common.StockType {
	thriftStockType := stockType
	commons.Logger.Debug("thriftStockTypeorigin StockType:", render.Render(stockId), ":", render.Render(thriftStockType))
	if stockId == "HSI" || stockId == "HSCE" || stockId == "HSCC" {
		thriftStockType = common.StockType_HK_INDX
		commons.Logger.Debug("thriftStockType StockType:", render.Render(stockId), ":", render.Render(thriftStockType))
	}
	return thriftStockType
}

//获取时间戳单位是秒(s)
func UnixTimeStamp() int64 {
	timesTamp := (time.Now().Unix())
	return timesTamp
}

//判断是否停牌
func IsSuspended(tradeStatus protocols.TradeStatus) bool {
	if tradeStatus == protocols.TradeStatus_Halt ||
		tradeStatus == protocols.TradeStatus_Stopt ||
		tradeStatus == protocols.TradeStatus_Susp {
		return true
	} else {
		return false
	}
}

//TODO protocols ErrCode 转为thrift 协议的forbagErrno
//ForbagErrno_UNKNOWN_ERR:2 ForbagErrno_SUCCESS:0 ForbagErrno_INNER_ERROR:3
func ThriftErr(errInfo *protocols.ErrInfo) common.ForbagErrno {
	var forbagErrno common.ForbagErrno
	if errInfo == nil {
		forbagErrno = common.ForbagErrno_INNER_ERROR
		commons.Logger.Error("errInfo is nil forbagErrno:", render.Render(forbagErrno))
	} else if errInfo.ErrCode == protocols.ErrCode_SUCCESS {
		forbagErrno = common.ForbagErrno_SUCCESS
	} else if errInfo.ErrCode == protocols.ErrCode_StockNotFound {
		forbagErrno = common.ForbagErrno_ERR_STOCK_IS_DELIST
	} else {
		forbagErrno = common.ForbagErrno_UNKNOWN_ERR
	}
	return forbagErrno
}

// thrift协议中的k线类型，转为protocols协议k线类型
// TODO 暂时为日,周,月K,后续根据需求添加K线类型
func KChartType2KLineType(kchartType common.KChartType) protocols.KLineType {
	switch kchartType {
	case common.KChartType_DAY:
		return protocols.KLineType_Day // 日K
	case common.KChartType_WEEK:
		return protocols.KLineType_Week // 周K
	case common.KChartType_MONTH:
		return protocols.KLineType_Month // 月k
	}
	return protocols.KLineType_Unknown // 未知类型
}

//thrift协议TimeChartPeriod转化为protocol的TimeChartPeriod
func ThriftPeriod2ProtPeriod(timeChartPeriod common.TimeChartPeriod) protocols.TimeChartPeriod {
	switch timeChartPeriod {
	case common.TimeChartPeriod_ONE_DAY:
		return protocols.TimeChartPeriod_OneDay // 1天分时
	case common.TimeChartPeriod_FIVE_DAY:
		return protocols.TimeChartPeriod_FiveDay // 5日分时
	}
	return protocols.TimeChartPeriod_Unknown // 未知类型分时
}

// 时间戳转化为string类型字符串 例:1508381400->2017-10-19 10:50:00
func Timestamp2String(timestampIns int64) string {
	unixTimeStamp := time.Unix(timestampIns+60*60*8, 0)
	strTime := unixTimeStamp.Format("2006-01-02 15:04:05") //2006-01-02 15:04:05这个每组数字都是有独特的含义
	return strTime
}

// 市场交易状态，protocols协议转化为thrift协议
func TypeTradeStatusConvert(tradeStatus protocols.StockTypeTradeStatus) common.TradeStatus {
	switch tradeStatus {
	case protocols.IsTrading:
		return common.TradeStatus_IS_TRADING // 交易状态
	case protocols.NotTrading:
		return common.TradeStatus_IS_NON_TRADING // 非交易状态
	}
	return common.TradeStatus_TRAD_UNKNOWN // 未知状态
}

// add by 2.3.8
func SetMinQuotedPrice(price float64, stockType protocols.StockType) float64 {
	if stockType == protocols.StockType_HK {
		if price >= 0.01 && price <= 0.25 {
			return 0.001
		} else if price > 0.25 && price <= 0.50 {
			return 0.005
		} else if price > 0.5 && price <= 10.00 {
			return 0.01
		} else if price > 10.00 && price <= 20.00 {
			return 0.02
		} else if price > 20.00 && price <= 100.00 {
			return 0.05
		} else if price > 100.00 && price <= 200.00 {
			return 0.1
		} else if price > 200.00 && price <= 500.00 {
			return 0.2
		} else if price > 500.00 && price <= 1000.00 {
			return 0.5
		} else if price > 1000.00 && price <= 2000.00 {
			return 1.00
		} else if price > 2000.00 && price <= 5000.00 {
			return 2.00
		} else if price > 5000.00 && price <= 9995.00 {
			return 5.00
		} else {
			commons.Logger.Error("invalid price ", price)
		}
		return 0.001
	} else {
		return 0.01
	}
}

// add by 2.3.8 卖盘最小变动单位
func SetMinQuotedSellPrice(price float64, stockType protocols.StockType) float64 {
	if stockType == protocols.StockType_HK {
		if price >= 0.01 && price < 0.25 {
			return 0.001
		} else if price >= 0.25 && price < 0.50 {
			return 0.005
		} else if price >= 0.5 && price < 10.00 {
			return 0.01
		} else if price >= 10.00 && price < 20.00 {
			return 0.02
		} else if price >= 20.00 && price < 100.00 {
			return 0.05
		} else if price >= 100.00 && price < 200.00 {
			return 0.1
		} else if price >= 200.00 && price < 500.00 {
			return 0.2
		} else if price >= 500.00 && price < 1000.00 {
			return 0.5
		} else if price >= 1000.00 && price < 2000.00 {
			return 1.00
		} else if price >= 2000.00 && price < 5000.00 {
			return 2.00
		} else if price >= 5000.00 && price < 9995.00 {
			return 5.00
		} else {
			commons.Logger.Error("invalid price ", price)
		}
		return 0.001
	} else {
		return 0.01
	}
}

// GetSellMarketDepths
func GetSellMarketDepths(marketDepth *protocols.MarketDepth, stockType common.StockType) []protocols.MarketDepthLevel {
	minQuotedPrice := 0.0000
	sellMarketDepths := make([]protocols.MarketDepthLevel, 0)
	commons.Logger.Debug("get marketDepth len:", len(marketDepth.BuyerMarketDepthLevels))
	lastSellPrice := marketDepth.SellerMarketDepthLevels[0].Price
	if marketDepth.SellerMarketDepthLevels[0].Price > 0.00001 {
		minQuotedPrice = SetMinQuotedSellPrice(marketDepth.SellerMarketDepthLevels[0].Price,
			ThriftType2ProtType(stockType))
	}
	for i := 0; i < len(marketDepth.BuyerMarketDepthLevels); i++ {
		if marketDepth.SellerMarketDepthLevels[i].Price > 0.00001 {
			minQuotedPrice = SetMinQuotedSellPrice(marketDepth.SellerMarketDepthLevels[i].Price,
				ThriftType2ProtType(stockType))
		}
		var tmpMarketDepth protocols.MarketDepthLevel
		// commons.Logger.Debug("Price:", marketDepth.SellerMarketDepthLevels[i].Price,
		// "lastPrice:", lastSellPrice, "Price-lastSellPrice:", marketDepth.SellerMarketDepthLevels[i].Price-lastSellPrice)
		if marketDepth.SellerMarketDepthLevels[i].Price-lastSellPrice < 0.00001 &&
			marketDepth.SellerMarketDepthLevels[i].Price-lastSellPrice > -0.00001 {
			tmpMarketDepth.Volume = marketDepth.SellerMarketDepthLevels[i].Volume
			tmpMarketDepth.Count = marketDepth.SellerMarketDepthLevels[i].Count
			tmpMarketDepth.Price = marketDepth.SellerMarketDepthLevels[i].Price
			sellMarketDepths = append(sellMarketDepths, tmpMarketDepth)
			lastSellPrice += minQuotedPrice
		} else if marketDepth.SellerMarketDepthLevels[i].Price-lastSellPrice > 0.00001 {
			count := int32(math.Floor(((marketDepth.SellerMarketDepthLevels[i].Price - lastSellPrice) / minQuotedPrice) + 0.5))
			for j := count; j >= 0; j-- {
				if lastSellPrice-marketDepth.SellerMarketDepthLevels[i].Price < -0.00001 {
					tmpMarketDepth.Volume = 0
					tmpMarketDepth.Count = 0
					tmpMarketDepth.Price = lastSellPrice
					sellMarketDepths = append(sellMarketDepths, tmpMarketDepth)
					lastSellPrice += minQuotedPrice
				} else if lastSellPrice-marketDepth.SellerMarketDepthLevels[i].Price > -0.00001 &&
					lastSellPrice-marketDepth.SellerMarketDepthLevels[i].Price < 0.00001 {
					tmpMarketDepth.Volume = marketDepth.SellerMarketDepthLevels[i].Volume
					tmpMarketDepth.Count = marketDepth.SellerMarketDepthLevels[i].Count
					tmpMarketDepth.Price = marketDepth.SellerMarketDepthLevels[i].Price
					sellMarketDepths = append(sellMarketDepths, tmpMarketDepth)
					lastSellPrice += minQuotedPrice
				}
			}
		} else if marketDepth.SellerMarketDepthLevels[i].Price-lastSellPrice < -0.00001 {
			for k := len(sellMarketDepths); k < 10; k++ {
				tmpMarketDepth.Volume = 0
				tmpMarketDepth.Count = 0
				tmpMarketDepth.Price = lastSellPrice
				sellMarketDepths = append(sellMarketDepths, tmpMarketDepth)
				lastSellPrice += minQuotedPrice
			}
		}
	}
	return sellMarketDepths
}

//
func GetBuyMarketDepths(marketDepth *protocols.MarketDepth, stockType common.StockType) []protocols.MarketDepthLevel {
	minQuotedPrice := 0.0000
	buyMarketDepths := make([]protocols.MarketDepthLevel, 0)
	commons.Logger.Debug("get buyermarketDepth len:", len(marketDepth.BuyerMarketDepthLevels))
	lastBuyerPrice := marketDepth.BuyerMarketDepthLevels[0].Price
	if marketDepth.BuyerMarketDepthLevels[0].Price > 0.00001 {
		minQuotedPrice = SetMinQuotedPrice(marketDepth.BuyerMarketDepthLevels[0].Price,
			ThriftType2ProtType(stockType))
	}
	for i := 0; i < len(marketDepth.BuyerMarketDepthLevels); i++ {
		if marketDepth.BuyerMarketDepthLevels[i].Price > 0.00001 {
			minQuotedPrice = SetMinQuotedPrice(marketDepth.BuyerMarketDepthLevels[i].Price,
				ThriftType2ProtType(stockType))
		}
		var tmpMarketDepth protocols.MarketDepthLevel
		// commons.Logger.Debug("buy Price:", marketDepth.BuyerMarketDepthLevels[i].Price,
		// " lastPrice:", lastBuyerPrice, " Price-lastBuyerPrice:", marketDepth.BuyerMarketDepthLevels[i].Price-lastBuyerPrice,
		// "minQuotedPrice", minQuotedPrice)
		if marketDepth.BuyerMarketDepthLevels[i].Price-lastBuyerPrice < 0.00001 &&
			marketDepth.BuyerMarketDepthLevels[i].Price-lastBuyerPrice > -0.00001 {
			tmpMarketDepth.Volume = marketDepth.BuyerMarketDepthLevels[i].Volume
			tmpMarketDepth.Count = marketDepth.BuyerMarketDepthLevels[i].Count
			tmpMarketDepth.Price = marketDepth.BuyerMarketDepthLevels[i].Price
			buyMarketDepths = append(buyMarketDepths, tmpMarketDepth)
			lastBuyerPrice -= minQuotedPrice
		} else if marketDepth.BuyerMarketDepthLevels[i].Price-lastBuyerPrice < -0.00001 &&
			marketDepth.BuyerMarketDepthLevels[i].Price > 0.00001 {
			// 精度
			count := int32(math.Floor(((lastBuyerPrice - marketDepth.BuyerMarketDepthLevels[i].Price) / minQuotedPrice) + 0.5))
			for j := count; j >= 0; j-- {
				if lastBuyerPrice-marketDepth.BuyerMarketDepthLevels[i].Price > 0.00001 {
					tmpMarketDepth.Volume = 0
					tmpMarketDepth.Count = 0
					tmpMarketDepth.Price = lastBuyerPrice
					buyMarketDepths = append(buyMarketDepths, tmpMarketDepth)
					lastBuyerPrice -= minQuotedPrice
				} else if lastBuyerPrice-marketDepth.BuyerMarketDepthLevels[i].Price > -0.00001 &&
					lastBuyerPrice-marketDepth.BuyerMarketDepthLevels[i].Price < 0.00001 {
					tmpMarketDepth.Volume = marketDepth.BuyerMarketDepthLevels[i].Volume
					tmpMarketDepth.Count = marketDepth.BuyerMarketDepthLevels[i].Count
					tmpMarketDepth.Price = marketDepth.BuyerMarketDepthLevels[i].Price
					buyMarketDepths = append(buyMarketDepths, tmpMarketDepth)
					lastBuyerPrice -= minQuotedPrice
				}
			}
		} else if marketDepth.BuyerMarketDepthLevels[i].Price < 0.00001 && marketDepth.BuyerMarketDepthLevels[i].Price > -0.00001 {
			for k := len(buyMarketDepths); k < 10; k++ {
				tmpMarketDepth.Volume = 0
				tmpMarketDepth.Count = 0
				tmpMarketDepth.Price = lastBuyerPrice
				buyMarketDepths = append(buyMarketDepths, tmpMarketDepth)
				lastBuyerPrice -= minQuotedPrice
			}
		}
	}
	return buyMarketDepths
}
