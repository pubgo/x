package xconfig_oss

import (
	"github.com/pubgo/x/logs"
	"github.com/pubgo/x/xconfig/xconfig_log"
	"github.com/pubgo/x/xdi"
)

var logger = logs.DebugLog("pkg", "oss")

func init() {
	xdi.InitInvoke(func(log xconfig_log.Log) {
		logger = log.With().Str("pkg", "oss").Logger()
	})
}
