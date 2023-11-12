package jlog

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/chroblert/jlog/jthirdutil/color"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// struct of FishLogger
type FishLogger struct {
	console           bool // displayed in console.default is false
	verbose           bool // if print log header. default is false
	iniCreateNewLog   bool
	maxStoreDays      int    // max store days
	maxSizePerLogFile int64  // max log file size per file. default 500MB
	size              int64  // all size，文件切割后，归0
	logFullPath       string //  logFullPath=logFileName+logFileExt
	logFilePerm       os.FileMode
	logDirPerm        os.FileMode
	logFileName       string        // file name
	logFileExt        string        // file suffix .log
	logCreateDate     string        // file create date
	logCount          int           // max log file count
	flushInterval     time.Duration // how long time,jlog write data to file
	bufferSize        int
	level             logLevel      // log level
	pool              sync.Pool     // Pool
	mu                sync.Mutex    // logger🔒
	writer            *bufio.Writer // buffer io
	file              *os.File      // log file object
	storeToFile       bool          // if save log data to file
	writed_size       int64         // 已经写入大小，文件切割后，不归0
	rotate_everyday   bool          // 设置是否每天分割文件。指定的文件不是本日的，则会创建新文件。default：false
}

type buffer struct {
	temp [64]byte
	bytes.Buffer
}

// log level
type logLevel int

// set log level
func (fl *FishLogger) SetLogLevel(lv logLevel) {
	if lv < DEBUG || lv > FATAL {
		panic("illegal log level")
	}
	fl.mu.Lock()
	defer fl.mu.Unlock()
	fl.level = lv
}

// set log path. eg: logs/app.log
// windows: \,/ as delimiter
// linux: / as delimiter [!!]
func (fl *FishLogger) SetLogFullPath(logFullPath string, mode ...os.FileMode) error {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	//fl.logFullPath = strings.ReplaceAll(logFullPath, "\\", "/")
	fl.logFullPath = logFullPath
	//fmt.Println("fl.logFullPath:", fl.logFullPath)
	// set log file path
	fl.logFileExt = filepath.Ext(fl.logFullPath) // .log
	//fmt.Println(fl.logFileExt)
	fl.logFileName = strings.TrimSuffix(fl.logFullPath, fl.logFileExt) // logs/app
	if fl.logFileExt == "" {
		fl.logFileExt = ".log"
	}
	var err error = nil
	if strings.HasSuffix(fl.logFileName, "/") {
		fl.logFileName = fl.logFileName + "app"
		err = fmt.Errorf("please specify correct log file path.eg: logs/app.log")
	}
	if len(mode) == 0 {
		err = os.MkdirAll(filepath.Dir(fl.logFullPath), fl.logDirPerm)
		if err != nil {
			panic(err)
		}
	} else {
		if filepath.Dir(fl.logFullPath) == "logs" {
			err = os.Chmod(filepath.Dir(fl.logFullPath), mode[0])
		} else {
			err = os.MkdirAll(filepath.Dir(fl.logFullPath), mode[0])
			if err != nil {
				panic(err)
			}
			err = os.Chmod(filepath.Dir(fl.logFullPath), mode[0])
			fl.logDirPerm = mode[0]
			if len(mode) > 1 {
				fl.logFilePerm = mode[1]
			}
		}
	}
	// 设置size
	fileInfo, err := os.Stat(fl.logFullPath)
	if err == nil {
		fmt.Println("设置fl.size:", fileInfo.Size())
		fl.size = fileInfo.Size()
	} else {
		fl.size = 0
	}
	return err
}

//	SetMaxSizePerLogFile
//
// eg. 10B,10KB,10MB,10GB. if not set correctly,will use default value 500MB.
func (fl *FishLogger) SetMaxSizePerLogFile(logFileSizeStr string) {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	fl.maxSizePerLogFile = transformFilesizeStrToInt64(logFileSizeStr)
}

// iniCreateNewLog
func (fl *FishLogger) IsIniCreateNewLog(iniCreateNewLog bool) {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	fl.iniCreateNewLog = iniCreateNewLog
}

// set max store days
// never delete if ma < 0
func (fl *FishLogger) SetMaxStoreDays(ma int) {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	fl.maxStoreDays = ma
}

// set log file count
// never delete if logCount < 0
func (fl *FishLogger) SetLogCount(logCount int) {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	fl.logCount = logCount
}

// flush to file
func (fl *FishLogger) Flush() {
	//fl.mu.Lock()
	//defer fl.mu.Unlock()
	fl.flushSync()
}

// set if print log header(caller file and line number)
func (fl *FishLogger) setVerbose(b bool) {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	fl.verbose = b
}

// set if displayed in console
func (fl *FishLogger) SetUseConsole(b bool) {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	fl.console = b
}

// set if save log data to log file
func (fl *FishLogger) SetStoreToFile(b bool) {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	fl.storeToFile = b
}

// set if rotate file every day
func (fl *FishLogger) SetRotateEveryday(b bool) {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	fl.rotate_everyday = b
}

// generate log header
func (fl *FishLogger) header(lv logLevel, depth int) *buffer {
	now := time.Now()
	buf := fl.pool.Get().(*buffer)
	year, month, day := now.Date()
	hour, minute, second := now.Clock()
	// format yyyymmdd hh:mm:ss.uuuu [DIWEF] file:line] msg
	buf.write4(0, year)
	buf.temp[4] = '/'
	buf.write2(5, int(month))
	buf.temp[7] = '/'
	buf.write2(8, day)
	buf.temp[10] = ' '
	buf.write2(11, hour)
	buf.temp[13] = ':'
	buf.write2(14, minute)
	buf.temp[16] = ':'
	buf.write2(17, second)
	buf.temp[19] = '.'
	buf.write4(20, now.Nanosecond()/1e5)
	buf.temp[24] = ' '
	copy(buf.temp[25:28], lv.Str())
	buf.temp[28] = ' '
	buf.Write(buf.temp[:29])
	// caller info
	if fl.verbose {
		_, file, line, ok := runtime.Caller(3 + depth)
		if !ok {
			file = "###"
			line = 1
		} else {
			slash := strings.LastIndex(file, "/")
			if slash >= 0 {
				file = file[slash+1:]
			}
		}
		buf.WriteString(file)
		buf.temp[0] = ':'
		n := buf.writeN(1, line)
		buf.temp[n+1] = ']'
		buf.temp[n+2] = ' '
		buf.Write(buf.temp[:n+3])
	}
	return buf
}

// print with \n
func (fl *FishLogger) println(lv logLevel, args ...interface{}) {
	if lv < fl.level {
		return
	}
	var buf *buffer
	// 11 represent Print()
	if lv == 11 {
		buf = &buffer{}
	} else {
		buf = fl.header(lv, 0)
	}
	fmt.Fprintln(buf, args...)
	// flush log data buffer to file
	fl.write(lv, buf, true)
}

// print with \n
// no log header
func (fl *FishLogger) nprintln(lv logLevel, args ...interface{}) {
	if lv < fl.level {
		return
	}
	var buf *buffer
	buf = &buffer{}
	fmt.Fprintln(buf, args...)
	// flush log data buffer to file
	fl.write(lv, buf, false)
}

// print with format,default no \n
// no log header
func (fl *FishLogger) nprintf(lv logLevel, format string, args ...interface{}) {
	if lv < fl.level {
		return
	}
	var buf *buffer
	buf = &buffer{}
	fmt.Fprintf(buf, format, args...)
	fl.write(lv, buf, false)
}

// print with format,default no \n
func (fl *FishLogger) printf(lv logLevel, format string, args ...interface{}) {
	if lv < fl.level {
		return
	}
	var buf *buffer
	if lv == 11 {
		buf = &buffer{}
	} else {
		buf = fl.header(lv, 0)
	}
	fmt.Fprintf(buf, format, args...)
	fl.write(lv, buf, true)
}

// wiret data to buffer
// isverbose: if has log header
func (fl *FishLogger) write(lv logLevel, buf *buffer, isverbose bool) {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	data := buf.Bytes()
	fl.writed_size += int64(len(data))
	if fl.console {
		switch lv {
		case DEBUG:
			// black background,blue text
			color.Blue(string(data))
		case INFO:
			// black background,white text
			color.White(string(data))
		case WARN:
			// black background,yellow text
			color.Yellow(string(data))
		case ERROR:
			// black background,red text
			color.Red(string(data))
		case FATAL:
			// black background,blue text，Display in reverse
			color.HiRed(string(data))
		default:
			color.White(string(data))
		}
	}
	if !fl.storeToFile {
		return
	}
	// first write data to file
	if fl.file == nil {
		if err := fl.rotate(); err != nil {
			os.Stderr.Write(data)
			fl.exit(err)
		}
	}

	// rotate file per day
	if fl.rotate_everyday && fl.logCreateDate != time.Now().Format("2006/01/02") {
		go fl.delete() // check old file perday
		if err := fl.rotate(); err != nil {
			fl.exit(err)
		}
	}

	// rotate file according to file size
	//log.Println("文件最大大小", fl.MaxSizePerLogFile)
	if fl.size+int64(len(data)) >= fl.maxSizePerLogFile {
		if err := fl.rotate(); err != nil {
			fl.exit(err)
		}
	}
	n, err := fl.writer.Write(data)
	fl.size += int64(n)
	if err != nil {
		fl.exit(err)
	}
	buf.Reset()
	fl.pool.Put(buf)
}

// delete old log
func (fl *FishLogger) delete() {
	if fl.maxStoreDays < 1 {
		return
	}
	dir := filepath.Dir(fl.logFullPath)
	fakeNow := time.Now().AddDate(0, 0, -fl.maxStoreDays)
	filepath.Walk(dir, func(fpath string, info os.FileInfo, err error) error {
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "logs: unable to delete old file '%s', error: %v\n", fpath, r)
			}
		}()
		if info == nil {
			return nil
		}
		// check
		if !info.IsDir() && info.ModTime().Before(fakeNow) && strings.HasSuffix(info.Name(), fl.logFileExt) && strings.HasPrefix(info.Name(), fl.logFileName+".") {
			return os.Remove(fpath)

		}
		return nil
	})
}

// flush buffer to file according to  timer
// write to file when monitor Ctrl+C
func (fl *FishLogger) daemon(stopChannel chan os.Signal) {
	tickTimer := time.NewTicker(fl.flushInterval)
	for {
		select {
		case <-tickTimer.C:
			fl.Flush()
		case <-stopChannel:
			fl.Flush()
			// 220111 bugfix
			os.Exit(-1)
		}
	}
}

// flush to file
func (fl *FishLogger) flushSync() {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	if fl.file != nil {
		fl.writer.Flush() // write data to memory
		fl.file.Sync()    // flush buffer(memory) to disk file
	}
}

func (fl *FishLogger) exit(err error) {
	fmt.Fprintf(os.Stderr, "logs: exiting because of error: %s\n", err)
	if err == nil {
		fl.flushSync()
	}
	os.Exit(0)
}

// rotate
// rotate file
// if first write data to file，
//
//	      -> check if app.log exist,if exist,then rename file
//			 -> create log file 'app.log'
//			 -> check current log file count if less than logCount config.if greater then delete old log file
//
// if not first write data to file，
//
//	-> check if current log file size less than maxLogFileSize.if not,delete old file
func (fl *FishLogger) rotate() error {
	now := time.Now()
	// rotate file
	// if file object is open,then flush data to memory to disk
	if fl.file != nil {
		// write to memory
		fl.writer.Flush()
		// flush to disk
		fl.file.Sync()
		// close file
		err := fl.file.Close()
		if err != nil {
			return err
		}
		// rename log file
		fileBackupName := filepath.Join(fl.logFileName + now.Format(".2006-01-02_150405") + fl.logFileExt)
		err = os.Rename(fl.logFullPath, fileBackupName)
		if err != nil {
			//log.Println("rename", err)
			return err
		}
		// create new log file app.log
		newLogFile, err := os.OpenFile(fl.logFullPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, fl.logFilePerm)
		if err != nil {
			return err
		}
		fl.file = newLogFile
		fl.size = 0
		// log buffer
		fl.writer = bufio.NewWriterSize(fl.file, fl.bufferSize)
	} else if fl.file == nil {
		if fl.iniCreateNewLog {
			// if is first run
			//    check if app.log exist.if exist,rename it
			_, err := os.Stat(fl.logFullPath)
			if err == nil {
				// get create date of file
				// rename log file
				fileBackupName := filepath.Join(fl.logFileName + now.Format(".2006-01-02_150405") + fl.logFileExt)
				err = os.Rename(fl.logFullPath, fileBackupName)
				if err != nil {
					return err
				}
			}
		}
		// create or open app.log
		newLogFile, err := os.OpenFile(fl.logFullPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, fl.logFilePerm)
		if err != nil {
			return err
		}
		fl.file = newLogFile
		fl.size = 0
		// log buffer
		fl.writer = bufio.NewWriterSize(fl.file, fl.bufferSize)
	}
	fileInfo, err := os.Stat(fl.logFullPath)
	fl.logCreateDate = now.Format("2006/01/02")
	if err == nil {
		// get size of current log file
		fl.size = fileInfo.Size()
		// get create date of current log file
		fl.logCreateDate = fileInfo.ModTime().Format("2006/01/02")
	}
	// fl.writer = bufio.NewWriterSize(fl.file, BufferSize)
	// check current log file count if less than logCount config.if greater then delete old log file
	if fl.logCount > 0 {
		pattern := fl.logFileName + ".*" + fl.logFileExt
		for files, _ := filepath.Glob(pattern); len(files) > fl.logCount; files, _ = filepath.Glob(pattern) {
			// delete log file
			os.Remove(files[0])
			if fl.level == -1 {
				tmpBuffer := fl.header(DEBUG, 0)
				fmt.Fprintf(tmpBuffer, "delete old log file")
				fmt.Fprintf(tmpBuffer, files[0])
				//fmt.Fprintf(tmpBuffer,"\033[0m")
				fmt.Fprintf(tmpBuffer, "\n")
				// black background, blue text
				//fmt.Fprintf(os.Stdout,"\033[1;34;40m"+string(tmpBuffer.Bytes())+"\033[0m")
				color.Blue(string(tmpBuffer.Bytes()))
				fl.writer.Write(tmpBuffer.Bytes())
			}
		}
	}
	return nil
}

// user customer instance
// 获取本次在所有文件中已经写入的大小
func (fl *FishLogger) GetAllWritedSize() int64 {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	return fl.writed_size
}

// 获取在当前文件中已经写入的大小
func (fl *FishLogger) GetCurrentFileSize() int64 {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	return fl.size
}

func (fl *FishLogger) Debug(args ...interface{}) {
	fl.println(DEBUG, args...)
}

func (fl *FishLogger) Debugf(format string, args ...interface{}) {
	fl.printf(DEBUG, format, args...)
}
func (fl *FishLogger) Info(args ...interface{}) {
	fl.println(INFO, args...)
}

func (fl *FishLogger) Infof(format string, args ...interface{}) {
	fl.printf(INFO, format, args...)
}

func (fl *FishLogger) Warn(args ...interface{}) {
	fl.println(WARN, args...)
}

func (fl *FishLogger) Warnf(format string, args ...interface{}) {
	fl.printf(WARN, format, args...)
}

func (fl *FishLogger) Error(args ...interface{}) {
	fl.println(ERROR, args...)
}

func (fl *FishLogger) Errorf(format string, args ...interface{}) {
	fl.printf(ERROR, format, args...)
}

func (fl *FishLogger) Fatal(args ...interface{}) {
	fl.println(FATAL, args...)
	os.Exit(0)
}

func (fl *FishLogger) Fatalf(format string, args ...interface{}) {
	fl.printf(FATAL, format, args...)
	os.Exit(0)
}

func (fl *FishLogger) NDebug(args ...interface{}) {
	fl.nprintln(DEBUG, args...)
}

func (fl *FishLogger) NDebugf(format string, args ...interface{}) {
	fl.nprintf(DEBUG, format, args...)
}
func (fl *FishLogger) NInfo(args ...interface{}) {
	fl.nprintln(INFO, args...)
}

func (fl *FishLogger) NInfof(format string, args ...interface{}) {
	fl.nprintf(INFO, format, args...)
}

func (fl *FishLogger) NWarn(args ...interface{}) {
	fl.nprintln(WARN, args...)
}

func (fl *FishLogger) NWarnf(format string, args ...interface{}) {
	fl.nprintf(WARN, format, args...)
}

func (fl *FishLogger) NError(args ...interface{}) {
	fl.nprintln(ERROR, args...)
}

func (fl *FishLogger) NErrorf(format string, args ...interface{}) {
	fl.nprintf(ERROR, format, args...)
}

func (fl *FishLogger) NFatal(args ...interface{}) {
	fl.nprintln(FATAL, args...)
	os.Exit(0)
}

func (fl *FishLogger) NFatalf(format string, args ...interface{}) {
	fl.nprintf(FATAL, format, args...)
	os.Exit(0)
}
