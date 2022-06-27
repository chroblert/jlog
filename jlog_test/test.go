package main

import (
	"fmt"
	"github.com/chroblert/jlog"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync"
	"time"
)

//func init(){
//	log.Println("test")
//}

func mainA() {
	go func() {
		http.ListenAndServe("0.0.0.0:8899", nil)
	}()
	jlog.SetLevel(jlog.DEBUG)
	//jlog.SetVerbose(false)
	jlog.Warn("warn: main")
	jlog.Println("xxx")
	jlog.Printf("%s\n", "testlll")
	fmt.Fprintln(os.Stderr, "xxxxxxx")
	jlog.NDebug("ndebug")
	var wg = &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(t int) {
			//jlog.Debug(t)
			wg.Done()
		}(i)
	}
	wg.Wait()
	jlog2 := jlog.New(jlog.LogConfig{
		BufferSize:        2048,
		FlushInterval:     3 * time.Second,
		MaxStoreDays:      5,
		MaxSizePerLogFile: "1000B",
		LogCount:          5,
		LogFullPath:       "logs/app2.log",
		Lv:                jlog.DEBUG,
		UseConsole:        false,
	})
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func(t int) {
			//jlog2.Debug("jlog2", t)
			//fmt.Println("fmt1",i)
			log.Println(t)
			wg.Done()
		}(i)
	}
	wg.Wait()
	jlog.Flush()
	jlog2.Flush()
	time.Sleep(5 * time.Second)
	fmt.Println("After 5s")
	time.Sleep(5 * time.Second)
	for {
	}

}
