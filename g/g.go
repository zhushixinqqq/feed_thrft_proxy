// Copyright [2013-2017] <xxxxx.Inc>
//
// Author: zhushixin
//
package g

import (
	"runtime"
)

const (
	VERSION = "0.1.0"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
