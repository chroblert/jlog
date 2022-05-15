package main

import (
	"github.com/chroblert/jlog"
)

var (
	nlog *jlog.FishLogger
)

func main() {
	nlog = jlog.New(jlog.LogConfig{
		BufferSize:        0,
		FlushInterval:     0,
		MaxStoreDays:      0,
		MaxSizePerLogFile: 1024,
		LogCount:          0,
		LogFullPath:       "",
		Lv:                0,
		UseConsole:        true,
		Verbose:           false,
		InitCreateNewLog:  false,
		StoreToFile:       true,
	})
	jlog.SetStoreToFile(true)
	jlog.SetMaxStoreDays(4)
	jlog.IsIniCreateNewLog(false)
	jlog.SetLogFullPath("logs\\test.log")
	jlog.SetLevel(jlog.DEBUG)
	jlog.SetMaxSizePerLogFile(1024)
	jlog.SetVerbose(false)
	jlog.Info("info")
	nlog.Error("error")
	nlog.Info("info")
	for i := 0; i < 10; i++ {
		nlog.NErrorf("a")
	}

}
