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

// log level
func (lv logLevel) Str() string {
	if lv >= DEBUG && lv <= FATAL {
		return logShort[lv*3 : lv*3+3]
	}
	return "[N]"
}

// newLogger instance logger
// path log file full path. eg:logs/app.log
func newLogger(logConf LogConfig) *FishLogger {
	fl := new(FishLogger)
	// set log config
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
	fl.rotate_everyday = logConf.RotateEveryDay

	fl.pool = sync.Pool{
		New: func() interface{} {
			return new(buffer)
		},
	}
	// 220509: not save log to log file
	if !fl.storeToFile {
		return fl
	}
	// set log file path
	// if not specify logfile path, use logs/app.log as default log file path
	if len(strings.TrimSpace(fl.logFullPath)) == 0 {
		fl.logFullPath = "logs/app.log"
	}
	fl.logFileExt = filepath.Ext(fl.logFullPath)                       // .log
	fl.logFileName = strings.TrimSuffix(fl.logFullPath, fl.logFileExt) // logs/app
	if fl.logFileExt == "" {
		fl.logFileExt = ".log"
	}
	// 设置size
	fileInfo, err := os.Stat(fl.logFullPath)
	if err == nil {
		fl.size = fileInfo.Size()
	} else {
		fl.size = 0
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

// create a new instance
// if not specify LogConfig, use default configuration to create.
//
//	write data to log file per 10 second
//	no limit store days,no limit log file count, max file size 500MB per file,displayed in console,create new log file only when first run or log file size > 500MB
//
// or create new instance with specified LogConfig
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
		RotateEveryDay:    false,
	})
}
