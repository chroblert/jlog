基于：https://github.com/zxysilent/logs 进行开发
## 0x01 介绍
### 功能
- 日志等级DEBUG、INFO、WARN、ERROR、FATAL
- 日志分级输出
- 不同等级不同颜色
- 日志文件个数限制
- 日志保存时间限制
- 日志文件自动分割
- 日志调用信息显示
- 可直接使用
- 可新建实例后使用

### 函数
带时间及调用信息输出：
- Debug(),Debugf()
- Info(),Infof()
- Warn(),Warnf()
- Error(),Errorf()
- Fatal(),Fatalf()
 
不带时间及调用信息输出：
- NDebug(),NDebugf()
- NInfo(),NInfof()
- NWarn(),NWarnf()
- NError(),NErrorf()
- NFatal(),NFatalf()

设置选项:
```go
jlog.SetUseConsole(true)                // 设置是否将日志输出到控制台
jlog.Flush()                            // 写入文件 // 主程序结束前调用
jlog.SetStoreToFile(true)               // 设置是否将日志存储到本地文件中
jlog.SetMaxStoreDays(4)                 // 是指日志文件最多存储多少天
jlog.IsIniCreateNewLog(false)           // 设置每次使用jlog的时候是否新建一个日志文件
jlog.SetLogFullPath("logs\\test.log")   // 设置日志保存的全路径
jlog.SetLevel(jlog.DEBUG)               // 设置日志等级，只有不低于该等级的日志才会显示
jlog.SetVerbose(false)                  // 设置是否显示调用jlog的文件名及行号
jlog.SetMaxSizePerLogFile(1024)         // 设置纯日志内容最大大小(jlog.Nxxx函数输出日志)，单位B，默认500MB
```
## 0x02 Use
### install
`go get -u github.com/chroblert/jlog`
### 直接使用
```
默认情况下，10s写入一次文件，不限制存储文件个数，不限制文件存储天数，日志路径为logs\app.log
默认输出日志到控制台，默认保存日志到文件，默认重复启动应用只创建一次文件
默认纯日志内容最大大小为256M
```
打印日志有如下方法:
```go
jlog.Debug()
jlog.Debugf()
jlog.Info()
jlog.Infof()
jlog.Warn()
jlog.Warnf()
jlog.Error()
jlog.Errorf()
jlog.Fatal()
jlog.Fatalf()
...
```
### 新建实例使用
```go
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
```