/*
 这是一个简单的配置参数的文件，
 需要手动设置的参数可以在这里设置
*/

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// 全局的配置文件参数
type Serverconfig struct {
	Suffix     string `json:"suffix"`     // 文件的后缀
	Prefix     string `json:"prefix"`     //文件的路径，这里命名意思为前缀
	ServerIp   string `json:"serverIp"`   // 服务器的地址
	ServerPort string `json:"serverPort"` //服务器的端口
	LoggerPath string `json:"loggerPath"` //日志文件的路径
}

// redis数据库配置文件仓库
type RDBConfig struct {
	RdbIp        string `json:"rdpIp"`        // redis数据库ip地址
	RdpPort      string `json:"rdpPort"`      // redis数据库端口
	PassWord     string `json:"password"`     //redis数据库密码
	DB           int    `json:"db"`           //启用的数据库(桶)
	PoolSize     int    `json:"poolSize"`     //连接池最大连接数
	MinIdleConns int    `json:"minIdleConns"` //最小空闲连接数
	MaxIdleConns int    `json:"maxIdleConns"` //最大空闲连接数
	PoolTimeout  int    `json:"poolTimeOut"`  //等待连接最大时间
	LoadCount    int    `json:"loadCount"`    //需要缓存的访问次数
	ExtendCount  int    `json:"extendCount"`  //需要延长存活时间的次数
}

type Settings struct {
	Serverconfig
	RDBConfig
	MineType map[string]interface{}
}

func parseConfig() {
	// 打开配置文件
	config, err := os.Open("./setting/Serverconfig.json")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	// 延迟关闭，避免内存泄露
	defer config.Close()

	// 打开数据库配置文件
	RDBConfig, err := os.Open("./setting/RDBConfig.json")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	// 延迟关闭，避免内存泄露
	defer RDBConfig.Close()

	mineType, err := os.Open("./setting/MINEType.json")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	// 延迟关闭，避免内存泄露
	defer mineType.Close()

	// 使用json的decoder,将config和rdbconfig文件解析成结构体
	decoder := json.NewDecoder(config)
	decoder.Decode(&setting.Serverconfig)

	decoder = json.NewDecoder(RDBConfig)
	decoder.Decode(&setting.RDBConfig)

	decoder = json.NewDecoder(mineType)
	decoder.Decode(&setting.MineType)

}
