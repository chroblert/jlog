package jlog

import (
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

// 字符串等级
func (lv logLevel) Str() string {
	if lv >= DEBUG && lv <= FATAL {
		return logShort[lv*3 : lv*3+3]
	}
	return "[N]"
}

// newLogger 实例化logger
// path 日志完整路径 eg:logs/app.log
func newLogger(logConf LogConfig) *FishLogger {
	fl := new(FishLogger)
	// 日志配置
	fl.bufferSize = logConf.BufferSize
	fl.flushInterval = logConf.FlushInterval
	fl.maxStoreDays = logConf.MaxStoreDays
	fl.maxSizePerLogFile = transformFilesizeStrToInt64(logConf.MaxSizePerLogFile)
	fl.logCount = logConf.LogCount
	fl.logFullPath = logConf.LogFullPath // logs/app.log
	fl.level = logConf.Lv
	fl.console = logConf.UseConsole
	fl.verbose = logConf.Verbose
	fl.iniCreateNewLog = logConf.InitCreateNewLog
	fl.storeToFile = logConf.StoreToFile
	fl.logFilePerm = logConf.LogFilePerm
	fl.logDirPerm = logConf.LogDirPerm

	fl.pool = sync.Pool{
		New: func() interface{} {
			return new(buffer)
		},
	}
	// 220509: 设置不将日志保存到文件
	if !fl.storeToFile {
		return fl
	}
	//日志文件路径设置
	// 若未传入明文路径，则使用logs\app.log作为默认路径
	if len(strings.TrimSpace(fl.logFullPath)) == 0 {
		fl.logFullPath = "logs/app.log"
	}
	fl.logFileExt = filepath.Ext(fl.logFullPath)                       // .log
	fl.logFileName = strings.TrimSuffix(fl.logFullPath, fl.logFileExt) // logs/app
	if fl.logFileExt == "" {
		fl.logFileExt = ".log"
	}
	os.MkdirAll(filepath.Dir(fl.logFullPath), fl.logDirPerm)
	if fl.flushInterval < 1 {
		fl.flushInterval = 10 * time.Second
	}
	signalChannel := make(chan os.Signal, 1)
	go fl.daemon(signalChannel)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	return fl
}

// 新建一个jlog示例
// 若不传入LogConfig，则使用默认的只进行创建。每十秒将日志写入文件，不限制存储天数和文件个数，单日志文件大小500MB，在控制台显示，每次运行不新建日志文件
// 否则，根据传入的LogConfig创建jlog示例
func New(logConfs ...LogConfig) *FishLogger {
	if len(logConfs) == 1 {
		logConf := logConfs[0]
		//if logConf.MaxSizePerLogFile < 1 {
		//	logConf.MaxSizePerLogFile = 524288000
		//}
		if logConf.LogFilePerm == 0 {
			logConf.LogFilePerm = 0644
		}
		if logConf.LogDirPerm == 0 {
			logConf.LogFilePerm = 0775
		}
		if logConf.BufferSize == 0 {
			logConf.BufferSize = 2048
		}
		if logConf.FlushInterval == 0 {
			logConf.FlushInterval = 10 * time.Second
		}
		return newLogger(logConf)
	}
	return newLogger(LogConfig{
		BufferSize:        2048,
		FlushInterval:     10 * time.Second,
		MaxStoreDays:      -1,
		MaxSizePerLogFile: "500MB", // 500MB
		LogCount:          -1,
		LogFullPath:       "logs/app.log",
		LogFilePerm:       0644,
		LogDirPerm:        0775,
		Lv:                DEBUG,
		UseConsole:        true,
		Verbose:           true,
		InitCreateNewLog:  false,
		StoreToFile:       true,
	})
}
