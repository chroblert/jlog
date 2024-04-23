package jlog

import (
	"os"
)

// set if print log header
func SetVerbose(b bool) {
	fishLogger.setVerbose(b)
}

// set if displayed in console
func SetUseConsole(b bool) {
	fishLogger.SetUseConsole(b)
}

// set log level
func SetLevel(lv logLevel) {
	fishLogger.SetLogLevel(lv)
}

// set max store days
// never delete if ma < 0
func SetMaxStoreDays(ma int) {
	fishLogger.SetMaxStoreDays(ma)
}

// set max log file count
// never delete if logCount < 0
func SetLogCount(logCount int) {
	fishLogger.SetLogCount(logCount)
}

// if create new log file when first run
func IsIniCreateNewLog(iniCreateNewLog bool) {
	fishLogger.IsIniCreateNewLog(iniCreateNewLog)
}

// set log file path. eg: logs/app.log
// windows: \,/ as path delimter
// linux: only / as path delimter [!!]
func SetLogFullPath(logFullPath string, mode ...os.FileMode) error {
	return fishLogger.SetLogFullPath(logFullPath, mode...)
}

// set max log file size
// eg. 10B,10KB,10MB,10GB. if not set correctly,will use default value 500MB.
func SetMaxSizePerLogFile(logfileSize string) {
	fishLogger.SetMaxSizePerLogFile(logfileSize)
}

// set if save log to log file
func SetStoreToFile(b bool) {
	fishLogger.SetStoreToFile(b)
}

// set if rotate file every day
func SetRotateEveryday(b bool) {
	fishLogger.SetRotateEveryday(b)
}

// -------- instance fishLogger
func Println(args ...interface{}) {
	fishLogger.nprintln(DEBUG, args...)
}
func Printf(format string, args ...interface{}) {
	fishLogger.nprintf(DEBUG, format, args...)
}

func Debug(args ...interface{}) {
	fishLogger.println(DEBUG, args...)
}

func Debugf(format string, args ...interface{}) {
	fishLogger.printf(DEBUG, format, args...)
}
func Info(args ...interface{}) {
	fishLogger.println(INFO, args...)
}

func Infof(format string, args ...interface{}) {
	fishLogger.printf(INFO, format, args...)
}

func Warn(args ...interface{}) {
	fishLogger.println(WARN, args...)
}

func Warnf(format string, args ...interface{}) {
	fishLogger.printf(WARN, format, args...)
}

func Error(args ...interface{}) {
	fishLogger.println(ERROR, args...)
}

func Errorf(format string, args ...interface{}) {
	fishLogger.printf(ERROR, format, args...)
}

func Fatal(args ...interface{}) {
	fishLogger.println(FATAL, args...)
	fishLogger.Flush()
	os.Exit(0)
}
func Fatalf(format string, args ...interface{}) {
	fishLogger.printf(FATAL, format, args...)
	fishLogger.Flush()
	os.Exit(0)
}

// flush to file
func Flush() {
	//fmt.Println("size1:",fishLogger.writer.Buffered())
	fishLogger.Flush()
	//fmt.Println("size2:",fishLogger.writer.Size())

}

func CloseAfterFlush() {
	fishLogger.CloseAfterFlush()
}

func NDebug(args ...interface{}) {
	fishLogger.nprintln(DEBUG, args...)
}

func NDebugf(format string, args ...interface{}) {
	fishLogger.nprintf(DEBUG, format, args...)
}
func NInfo(args ...interface{}) {
	fishLogger.nprintln(INFO, args...)
}

func NInfof(format string, args ...interface{}) {
	fishLogger.nprintf(INFO, format, args...)
}

func NWarn(args ...interface{}) {
	fishLogger.nprintln(WARN, args...)
}

func NWarnf(format string, args ...interface{}) {
	fishLogger.nprintf(WARN, format, args...)
}

func NError(args ...interface{}) {
	fishLogger.nprintln(ERROR, args...)
}

func NErrorf(format string, args ...interface{}) {
	fishLogger.nprintf(ERROR, format, args...)
}

func NFatal(args ...interface{}) {
	fishLogger.nprintln(FATAL, args...)
	fishLogger.Flush()
	os.Exit(0)
}
func NFatalf(format string, args ...interface{}) {
	fishLogger.nprintf(FATAL, format, args...)
	fishLogger.Flush()
	os.Exit(0)
}
