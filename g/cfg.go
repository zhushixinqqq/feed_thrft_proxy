// Copyright [2013-2017] <xxxxx.Inc>

//
// Author: zhushixin
package g

import (
	commons "common"
	"encoding/json"
	"sync"

	"github.com/toolkits/file"
)

type ThriftServerPortalConfig struct {
	Address                    string `json:"address"`
	Enabled                    bool   `json:"enabled"`
	ServiceName                string `json:"serviceName"`
	FeedThriftPoolNum          int32  `json:"feedThriftPoolNum"`
	FeedThriftTimeOut          int32  `json:"feedThriftTimeOut"`
	FeedExPoolNum              int32  `json:"feedExPoolNum"`
	FeedExTimeOut              int32  `json:"feedExTimeOut"`
	MarketDepthPoolNum         int32  `json:"marketDepthPoolNum"`
	MarketDepthTimeOut         int32  `json:"marketTimeOut"`
	StocksByTypePoolNum        int32  `json:"stocksByTypePoolNum"`
	StocksByTypeTimeOut        int32  `json:"stocksByTypeTimeOut"`
	StockTradeStatusPoolNum    int32  `json:"stockTradeStatusPoolNum"`
	StockTradeStatusTimeOut    int32  `json:"stockTradeStatusTimeOut"`
	KChartPoolNum              int32  `json:"kChartPoolNum"`
	KChartTimeOut              int32  `json:"kChartTimeOut"`
	TimeChartPoolNum           int32  `json:"timeChartPoolNum"`
	TimeChartTimeOut           int32  `json:"timeChartTimeOut"`
	StockBasicInfoBatchPoolNum int32  `json:"stockBasicInfoBatchPoolNum"`
	StockBasicInfoBatchTimeOut int32  `json:"stockBasicInfoBatchTimeOut"`
	StockTicksPoolNum          int32  `json:"stockTicksPoolNum"`
	StockTicksTimeOut          int32  `json:"stockTicksTimeOut"`
	RealBusAddr                string `json:"realBusAddr"`
	ChartBusAddr               string `json:"chartBusAddr"`
}

type EtcdServerPortalConfig struct {
	EtcdUrl          string `json:"etcdUrl"`
	ServiceName      string `json:"serviceName"`
	BasePath         string `json:"basePath"`
	ChartServiceName string `json:"chartserviceName"`
	Used             bool   `json:"used"`
}

type GlobalConfig struct {
	ThriftServer *ThriftServerPortalConfig `json:"thriftServer"`
	EtcdServer   *EtcdServerPortalConfig   `json:"etcdServer"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		commons.Logger.Critical("use -c to specify configuration file")
	}
	if !file.IsExist(cfg) {
		commons.Logger.Critical("config file:", cfg, "is not existent")
	}
	ConfigFile = cfg
	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		commons.Logger.Critical("read config file:", cfg, "fail:", err)
	}
	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		commons.Logger.Critical("parse config file:", cfg, "fail:", err)
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &c
	commons.Logger.Info("read config file:", cfg, "successfully")
}
