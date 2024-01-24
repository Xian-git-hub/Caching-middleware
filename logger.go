/*
	此模块负责日志操作相关的功能，例如日志文件，相关文件夹的创建，
	日志分为日常日志和错误日志，日志结构为
	├─dailyLog
	|	└─2024-01
	└─errorLog
    	└─2024-01
	dialyLog文件夹里面存放子文件夹，代表每个月的日志文件，
	文件夹里面就有这个月每天的日志文件
	errorLog文件夹同理

*/

package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	dailySecondDir = "/dailyLog/"
	errorSecondDir = "/errorLog/"
	suffix         = ".log"
)

var firstDir string

// 一个自定义的logger结构体，包含了日常日志和错误日志
type MyLogger struct {
	timer       *time.Timer
	dailyLogger *log.Logger
	errorLogger *log.Logger
	dailyFile   *os.File
	errorFile   *os.File
}

// 工厂函数，返回一个MyLogger结构体，用于记录日志
func NewMyLogger() (ml *MyLogger) {

	firstDir = setting.LoggerPath

	dailyLog, errorLog := createLogFile()

	myLog := &MyLogger{
		dailyLogger: log.New(dailyLog, "dailyLog:", log.Ldate|log.Ltime|log.Lshortfile|log.LstdFlags),
		errorLogger: log.New(errorLog, "errorLog:", log.Ldate|log.Ltime|log.Lshortfile|log.LstdFlags),
		dailyFile:   dailyLog,
		errorFile:   errorLog,
	}
	myLog.setTimer()

	return myLog
}

// 监听函数，监听是否有timer往通道里面发送数据
// 执行对应的操作
func (myLog *MyLogger) Listener() {

	// 监听定时器的channel,当时间到了就创建一个新的定时器
	// 改进?:第一次setTimer可以使用一个一次性的timer,之后
	// 程序运行是会稳定在24h之后创建日志文件，可以使用ticker
	for range myLog.timer.C {
		fmt.Println("time is up:", time.Now().Format("2006-01-02 15:04:05"))
		myLog.createLogFileAuto()
		nextTime := getNextCreateTime()
		myLog.timer.Reset(nextTime)
	}

}

// 设置一个timer,用于定时创建日志文件
// 这个timer会自动获取距离下次创建日志的时间
func (myLog *MyLogger) setTimer() {

	// 获取下次创建日志文件的时间段，用这个时间段创建一个定时器
	nextTime := getNextCreateTime()
	myLog.timer = time.NewTimer(nextTime)

}

// 关闭旧的日志文件，并创建新的日志文件
func (myLog *MyLogger) createLogFileAuto() {

	// fmt.Println("hahaha")
	// 关闭旧的日志文件
	myLog.dailyFile.Close()
	myLog.errorFile.Close()

	// 创建新的日志文件
	newDailyFile, newErrorFile := createLogFile()

	// 将新的日志文件设置好
	myLog.dailyFile = newDailyFile
	myLog.errorFile = newErrorFile
	myLog.dailyLogger = log.New(newDailyFile, "dailyLog:", log.Ldate|log.Ltime|log.Lshortfile|log.LstdFlags)
	myLog.errorLogger = log.New(newErrorFile, "errorLog:", log.Ldate|log.Ltime|log.Lshortfile|log.LstdFlags)

	// fmt.Println("lalala")
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

	// fmt.Println("nextTime:", tomrrowZero.Format("2006-01-02 15:04:05"))
	// fmt.Println("subTime:", int64(nextTime.Seconds()))

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
