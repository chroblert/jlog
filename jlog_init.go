package jlog

import (
	"os"
	"strconv"
	"time"
	"unicode"
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
		MaxSizePerLogFile: "500MB", // 500MB
		LogCount:          -1,
		LogFullPath:       "logs/app.log",
		LogFilePerm:       0644,
		Lv:                DEBUG,
		UseConsole:        true,
		Verbose:           true,
		InitCreateNewLog:  false,
		StoreToFile:       true,
	})
)

type LogConfig struct {
	BufferSize    int
	FlushInterval time.Duration // 单位ms。若为0，则默认设置10s
	MaxStoreDays  int
	//MaxSizePerLogFile int64 // 单位B，默认500M
	MaxSizePerLogFile string // 单位B，默认500M
	LogCount          int
	LogFullPath       string
	LogFilePerm       os.FileMode
	LogDirPerm        os.FileMode
	Lv                logLevel
	UseConsole        bool
	Verbose           bool
	InitCreateNewLog  bool
	StoreToFile       bool
}

// 将描述性文件大小转为int64。支持：B,KB,MB,GB
// 若不符合格式或大小为0，则设置默认值 500MB
func transformFilesizeStrToInt64(logFileSizeStr string) int64 {
	var number int64 = 0
	var logfileSize int64 = 0
	for i, c := range logFileSizeStr {
		if unicode.IsDigit(c) {
			tmpNum, _ := strconv.Atoi(string(c))
			number = number*10 + int64(tmpNum)
		} else {
			switch logFileSizeStr[i:] {
			case "B":
				logfileSize = number
			case "KB":
				logfileSize = number * 1024
			case "MB":
				logfileSize = number * 1024 * 1024
			case "GB":
				logfileSize = number * 1024 * 1024 * 1024
			default:
				logfileSize = 500 * 1024 * 1024
			}
			break
		}
	}
	if logfileSize == 0 {
		logfileSize = 500 * 1024 * 1024
	}
	return logfileSize
}
