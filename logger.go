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
	suffix         = ".txt"
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

	return &MyLogger{
		dailyLogger: log.New(dailyLog, "dailyLog:", log.Ldate|log.Ltime|log.Lshortfile|log.LstdFlags),
		errorLogger: log.New(errorLog, "errorLog:", log.Ldate|log.Ltime|log.Lshortfile|log.LstdFlags),
		dailyFile:   dailyLog,
		errorFile:   errorLog,
	}
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

// 监听函数，监听是否有timer往通道里面发送数据
func (myLog *MyLogger) Listener() {
	// timer发送了数据，说明创建下一个日志文件的时间到了
	<-myLog.timer.C

	// 先关闭之前的文件并创建今天的文件
	// 避免误差，先沉睡2s
	time.Sleep(2 * time.Second)
	myLog.createLogFileAuto()

	// 设置下一个timer
	myLog.setTimer()

}

// 关闭旧的日志文件，并创建新的日志文件
func (myLog *MyLogger) createLogFileAuto() {
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
}

// 设置一个timer,用于定时创建日志文件
// 这个timer会自动获取距离下次创建日志的时间
func (myLog *MyLogger) setTimer() {
	// 获得距离下次创建日志文件的秒数
	duration := getNextCreateTime()
	// 用这个秒数设置一个timer
	myLog.timer = time.NewTimer(time.Duration(duration))
}

// 返回距离下次创建日志文件的时间
func getNextCreateTime() int64 {

	// 获取当前时间
	now := time.Now()

	// 获取第二天的零点时间的Unix时间戳
	// 这里不需要担心天数越界问题，golang会帮我们处理
	tomrrowZero := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)
	nextTime := tomrrowZero.Sub(now)

	return int64(nextTime.Seconds())

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
