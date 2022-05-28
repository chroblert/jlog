package jlog

import (
	"time"
)

//
//var applog = new(FishLogger)
//
//// 20201011: 使用Logs
//func InitLogs(logpath string, amaxSize int64, amaxAge, alogCount int) {
//	maxSize = amaxSize // 单个文件最大大小
//	maxAge = amaxAge   // 单个文件保存2天
//	LogCount = alogCount
//	applog = newLogger(logpath)
//	defer applog.Flush()
//	applog.SetLogLevel(DEBUG)
//	applog.setVerbose(true)
//	applog.SetUseConsole(true)
//	//applog.info("test")
//}
//func Println(args ...interface{}) {
//	// applog.info(args)
//	applog.println(INFO, args...)
//}
//
//func Printf(format string, args ...interface{}) {
//	// applog.infof(format, args...)
//	applog.printf(INFO, format, args...)
//}

const (
	DEBUG logLevel = iota
	INFO
	WARN
	ERROR
	FATAL

	digits   = "0123456789"
	logShort = "[D][I][W][E][F]"
)

var (
	fishLogger = newLogger(LogConfig{
		BufferSize:        2048,
		FlushInterval:     10 * time.Second,
		MaxStoreDays:      -1,
		MaxSizePerLogFile: 512000000,
		LogCount:          -1,
		LogFullPath:       "logs/app.log",
		Lv:                DEBUG,
		UseConsole:        true,
		Verbose:           true,
		InitCreateNewLog:  false,
	})
)

type LogConfig struct {
	BufferSize        int
	FlushInterval     time.Duration
	MaxStoreDays      int
	MaxSizePerLogFile int64 // 单位B，默认500M
	LogCount          int
	LogFullPath       string
	Lv                logLevel
	UseConsole        bool
	Verbose           bool
	InitCreateNewLog  bool
	StoreToFile       bool
}
