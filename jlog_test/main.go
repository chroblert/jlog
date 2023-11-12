package main

import (
	"github.com/chroblert/jlog"
)

var (
// nlog *jlog.FishLogger
)

func main() {
	//nlog = jlog.New(jlog.LogConfig{
	//	BufferSize:        0,
	//	FlushInterval:     0,
	//	MaxStoreDays:      0,
	//	MaxSizePerLogFile: "1000B",
	//	LogCount:          0,
	//	LogFullPath:       "",
	//	Lv:                0,
	//	UseConsole:        true,
	//	Verbose:           true,
	//	InitCreateNewLog:  false,
	//	StoreToFile:       true,
	//	LogFilePerm:       0,
	//})
	////jlog.SetStoreToFile(true)
	////jlog.SetMaxStoreDays(4)
	//nlog.IsIniCreateNewLog(false)
	//nlog.SetLogFullPath("logs/1111.test", 0777)
	////jlog.SetLevel(jlog.DEBUG)
	//nlog.SetMaxSizePerLogFile("10MB")
	////jlog.SetVerbose(true)
	////jlog.Info("info1")
	////jlog.Warn("warn1")
	//nlog.Error("error")
	//nlog.Info("info2")
	//for i := 0; i < 10; i++ {
	//	nlog.NErrorf("a")
	//}
	defer jlog.Flush()
	jlog.SetLogFullPath("logs/testd.log", 0777, 0644)
	jlog.NInfo("jlog test")
	jlog.Flush()
	jlog.NInfo(jlog.GetCurrentFileSize())

	jlog.NInfo(jlog.GetAllWritedSize())
}
