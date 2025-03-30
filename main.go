package main

import (
	"sync"
	datatype "wihscan/dataType"
	"wihscan/factory"
	"wihscan/global"
	"wihscan/options"
	"wihscan/scan"
	"wihscan/util"

	"github.com/gookit/color"
)

var sem chan struct{}
var writePath string

func init() {
	// // 初始化配置文件
	// util.InitRule()
}

func main() {
	color.C256(45).Printf(global.LOGO, global.Version)
	option := options.Options()
	if option == nil {
		return
	}
	writePath = option.OutputFilePath

	// 加载规则
	scan.RuleLoad()

	urls := factory.Factory(option)
	// 遍历扫描urljs
	wg := &sync.WaitGroup{}
	// 创建一个信号量，用于限制并发量
	sem = make(chan struct{}, option.Thread)
	for _, url := range urls {
		wg.Add(1)
		// 在 goroutine 启动之前获取信号量
		sem <- struct{}{}
		go thread(url, option, wg)
	}
	wg.Wait()
}

func thread(url string, option *datatype.Option, wg *sync.WaitGroup) {
	resultScan := scan.Scan(url)
	util.FormatOutput(resultScan)
	if option.OutputJson {
		util.FormatOutputWriteJson(resultScan, writePath)
	} else {
		util.FormatOutputWrite(resultScan, writePath)
	}
	<-sem
	wg.Done()
}
