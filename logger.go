/*
	此模块负责日志操作相关的功能，例如日志文件，相关文件夹的创建，
	日志分为日常日志和错误日志，日志结构为
	├─dailyLog
	|	└─2024-01
	└─errorLog
    	└─2024-01
	dailyLog文件夹里面存放子文件夹，代表每个月的日志文件，
	文件夹里面就有这个月每天的日志文件
	errorLog文件夹同理

*/

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	dailySecondDir = "/dailyLog/"
	errorSecondDir = "/errorLog/"
	dataSecondDir  = "/dataLog"
	suffix         = ".log"
	dailyType      = "daily"
	errorType      = "error"
	dataType       = "data"
)

var firstDir string
var enable bool = false

// 一个自定义的logger结构体，包含了日常日志和错误日志
type MyLogger struct {
	ticker *time.Ticker
	timer  *time.Timer
	DailyLog
	ErrorLog
	DailyBuf
	ErrorBuf
}

// 日常日志
type DailyLog struct {
	dailyLogger *log.Logger
	dailyFile   *os.File
}

// 错误日志
type ErrorLog struct {
	errorLogger *log.Logger
	errorFile   *os.File
}

// 带锁的缓存
type DailyBuf struct {
	dmu     chan bool
	dbuffer *bufio.Writer
}

type ErrorBuf struct {
	emu     chan bool
	ebuffer *bufio.Writer
}

// 工厂函数，返回一个MyLogger结构体，用于记录日志
func NewMyLogger() (ml *MyLogger) {

	firstDir = setting.LoggerPath

	// 创建日志文件
	dailyLog, errorLog := createLogFile()

	// 让指针指向申请的空间
	ml = new(MyLogger)

	// 创建结构体并给其字段赋值
	ml.dbuffer = bufio.NewWriter(dailyLog)
	ml.ebuffer = bufio.NewWriter(errorLog)
	ml.dailyFile = dailyLog
	ml.errorFile = errorLog
	ml.dailyLogger = log.New(ml.dbuffer, "dailyLog:", log.Ldate|log.Ltime|log.Lshortfile)
	ml.errorLogger = log.New(ml.ebuffer, "errorLog:", log.Ldate|log.Ltime|log.Lshortfile)

	// 初始化channel并放入一个值,代表buf空闲
	ml.dmu = make(chan bool, 1)
	ml.emu = make(chan bool, 1)
	ml.dmu <- true
	ml.emu <- true

	ml.setTimer()
	ml.setTicker()

	return ml
}

// 监听函数，监听是否有timer往通道里面发送数据
// 执行对应的操作
func (myLog *MyLogger) Listener() {

	// 监听定时器的channel,当时间到了就创建一个新的定时器
	// 改进?:第一次setTimer可以使用一个一次性的timer,之后
	// 程序运行是会稳定在24h之后创建日志文件，可以使用ticker

	for {
		select {
		// 定时日志
		case <-myLog.timer.C:
			myLog.createLogFileAuto()
			nextTime := getNextCreateTime()
			myLog.timer.Reset(nextTime)
			c.reset()

		// 定时flush
		case <-myLog.ticker.C:
			myLog.flushDBuffer()
			myLog.flushEBuffer()

		// 退出信号
		case c := <-exitChan:
			fmt.Println("got signal:", c)
			myLog.doLog(dailyType, "server close! Bye Bye")
			closeSource()
			os.Exit(0)

		// 计算统计数据信号
		case <-c.cticker.C:
			ratio := c.cal()
			myLog.doLog(dailyType, "redis:"+ratio+"%")
		}

	}

}

// 设置一个timer,用于定时创建日志文件
// 这个timer会自动获取距离下次创建日志的时间
func (myLog *MyLogger) setTimer() {

	// 获取下次创建日志文件的时间段，用这个时间段创建一个定时器
	nextTime := getNextCreateTime()
	myLog.timer = time.NewTimer(nextTime)

}

// 设置一个8s的ticker,用于把数据flush到硬盘中
func (myLog *MyLogger) setTicker() {
	myLog.ticker = time.NewTicker(time.Duration(setting.FlushTime) * time.Second)
}

// 记录日志操作，由调用方传入使用的logger类型以及记录的字符串
// 这里不需要担心数据截断问题，在bufio.writer更换文件时会flush
func (ml *MyLogger) doLog(name string, content string) {
	switch name {
	case dailyType:
		if enable {
			break
		}
		<-ml.dmu
		ml.dailyLogger.Println(content)
		ml.dmu <- true
	case errorType:
		<-ml.emu
		ml.errorLogger.Println(content)
		ml.emu <- true
	}
}

// flush一次日常日志
func (ml *MyLogger) flushDBuffer() {
	<-ml.dmu
	ml.dbuffer.Flush()
	ml.dmu <- true
}

// flush一次异常日志
func (ml *MyLogger) flushEBuffer() {
	<-ml.emu
	ml.ebuffer.Flush()
	ml.emu <- true
}

// 关闭旧的日志文件，并创建新的日志文件
func (myLog *MyLogger) createLogFileAuto() {

	// 在关闭就的日志文件之前flush掉，避免丢失一部分log
	myLog.flushDBuffer()
	myLog.flushEBuffer()

	// 关闭旧的日志文件
	myLog.dailyFile.Close()
	myLog.errorFile.Close()

	// 创建新的日志文件
	newDailyFile, newErrorFile := createLogFile()

	// 将buffer指向新的日志文件
	myLog.dbuffer.Reset(newDailyFile)
	myLog.ebuffer.Reset(newErrorFile)

	// 将新的日志文件设置好
	myLog.dailyFile = newDailyFile
	myLog.errorFile = newErrorFile

	// 让日志对象指向新的buffer
	myLog.dailyLogger = log.New(myLog.dbuffer, "dailyLog:", log.Ldate|log.Ltime|log.Lshortfile|log.LstdFlags)
	myLog.errorLogger = log.New(myLog.ebuffer, "errorLog:", log.Ldate|log.Ltime|log.Lshortfile|log.LstdFlags)

}

// 关闭程序之前关闭所有资源
func closeSource() {

	// 先flush日志
	myLog.flushDBuffer()
	myLog.flushEBuffer()

	// 关闭日志文件
	myLog.dailyFile.Close()
	myLog.errorFile.Close()

	// 停止所有计时器
	myLog.ticker.Stop()
	myLog.timer.Stop()

}

// 创建对应的日志文件
func createLogFile() (dailyFile *os.File, errorFile *os.File) {

	// 子文件夹的名字
	dirName := getDateStringMonth()

	// 创建dailyLog文件夹，如果文件夹已存在则没有动作
	err := os.MkdirAll(firstDir+dailySecondDir+dirName, 0755)
	if err != nil {
		fmt.Println("err:", err)
	}
	// 同理创建errorLog文件夹
	err = os.MkdirAll(firstDir+errorSecondDir+dirName, 0755)
	if err != nil {
		fmt.Println("err:", err)
	}

	// 日志文件的名字
	fileName := getDateStringDay()
	// fileName := getDateString()
	// 拼接文件路径
	dailyFileName := firstDir + dailySecondDir + dirName + "/" + fileName + suffix

	// 以新建|追加的方式打开文件夹
	dailyFile, err = os.OpenFile(dailyFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("err:", err)
	}
	// 同理
	errorFileName := firstDir + errorSecondDir + dirName + "/" + fileName + suffix
	errorFile, err = os.OpenFile(errorFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("err:", err)
	}
	return

}

// 返回距离下次创建日志文件的时间
func getNextCreateTime() time.Duration {

	// 获取当前时间
	now := time.Now()

	// 获取第二天的零点时间的Unix时间戳
	// 这里不需要担心天数越界问题，golang会帮我们处理
	tomrrowZero := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)
	nextTime := tomrrowZero.Sub(now)

	return nextTime
}

// 返回当前年-月的字符串
func getDateStringMonth() string {
	now := time.Now()
	return now.Format("2006-01")
}

// 返回当前年-月-日的字符串
func getDateStringDay() string {
	now := time.Now()
	return now.Format("2006-01-02")
}

// 返回当前年-月-日-时-分-秒的字符串
// func getDateString() string {
// 	now := time.Now()
// 	return now.Format("2006-01-02 15-04-05")
// }
