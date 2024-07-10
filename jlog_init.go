package jlog

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"
)

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
	path, _    = os.Executable()
	_, exec    = filepath.Split(path)
	fishLogger = newLogger(LogConfig{
		BufferSize:        2048,
		FlushInterval:     10 * time.Second,
		MaxStoreDays:      -1,
		MaxSizePerLogFile: "500MB", // 500MB
		LogCount:          -1,
		LogFullPath:       "logs/" + exec + "_" + time.Now().Format("20060102T150405") + ".log",
		LogFilePerm:       0644,
		LogDirPerm:        0755,
		Lv:                DEBUG,
		UseConsole:        true,
		Verbose:           true,
		InitCreateNewLog:  false,
		StoreToFile:       true,
		RotateEveryDay:    false,
	})
)

// LogConfig
// MaxSizePerLogFile max log file size per file,unit:B,KB,MB,GB. default 500MB. -1 will not limit file size
type LogConfig struct {
	BufferSize        int
	FlushInterval     time.Duration // unit:ms。if value is 0，then use default 10s
	MaxStoreDays      int
	MaxSizePerLogFile string // unit:B，default 500M
	LogCount          int
	LogFullPath       string
	LogFilePerm       os.FileMode
	LogDirPerm        os.FileMode
	Lv                logLevel
	UseConsole        bool
	Verbose           bool
	InitCreateNewLog  bool
	StoreToFile       bool
	RotateEveryDay    bool
}

// transform fileSizeStr to int64。support：B,KB,MB,GB
// if set to -1, will not limit log file size
// if set illegal value or value size is 0,use default value 500MB
func transformFilesizeStrToInt64(logFileSizeStr string) int64 {
	if strings.TrimSpace(logFileSizeStr) == "-1" {
		return -1
	}
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
