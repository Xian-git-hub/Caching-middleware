package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client // 全局的go-redis里的redis客户端，通过这个访问redis
var myLog *MyLogger   // 全局的日志对象，用来记录日志
var setting Settings  // 全局的参数对象，使用参数

type helloHandlerStruct struct {
	content string
}

func init() {

	parseConfig()

	rdb = redis.NewClient(&redis.Options{
		Addr:     setting.RdbIp + setting.RdpPort,
		Password: setting.RdpPort, // 没有密码，默认值
		DB:       setting.DB,      // 默认DB 0

		PoolSize:     setting.PoolSize,     //最大连接数
		MinIdleConns: setting.MinIdleConns, //最小空闲连接数
		MaxIdleConns: setting.MaxIdleConns, //最大空闲连接数

		PoolTimeout: time.Duration(setting.PoolTimeout) * time.Second, //等待连接最长时间，这里设置为5s
	})

}

func main() {
	// syscall.Umask(0)

	// 创建mylog结构体变量
	myLog = NewMyLogger()
	// 第一次运行，设置第一个timer
	// myLog.setTimer()
	// 启动一个协程，让其监听timer对channel的操作
	go myLog.Listener()

	mux := http.NewServeMux()

	mux.HandleFunc("/greet", greetingHandler)
	mux.HandleFunc("/download", handleRequestFile)
	mux.Handle("/", &helloHandlerStruct{"hello world handled by struct!"})
	mux.HandleFunc("/test", handledTest)

	myLog.dailyLogger.Println("server start! welcome")

	http.ListenAndServe(setting.ServerIp+setting.ServerPort, mux)

	fmt.Println("server close! Bye Bye")
	myLog.dailyLogger.Println("server close! Bye Bye")
}

func (handler *helloHandlerStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, handler.content)
}

func handledTest(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	fileName := query.Get("file")
	fileSuffix := getFileSuffix(fileName)

	fmt.Fprint(w, setting.MineType[fileSuffix])
}

func greetingHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	fileName := query.Get("file")
	fmt.Fprint(w, fileName)

}

// 处理预加载逻辑
// func handlePreload(w http.ResponseWriter, r *http.Response) {

// }

// 处理请求文件逻辑
func handleRequestFile(w http.ResponseWriter, r *http.Request) {

	// 获取参数,拼接出文件名及文件的磁盘路径
	query := r.URL.Query()
	fileName := query.Get("file") + setting.Suffix
	filePath := setting.Prefix + fileName

	// 获取文件,以[]byte形式
	data, err := getFile(fileName, filePath)
	if err != nil {
		myLog.errorLogger.Printf("getFile err:%v\n", err)
		fmt.Println("err1:", err)
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprint(w, "您请求的数据服务器中不存在，请联系管理员")
		}
		if data != nil {
			myLog.dailyLogger.Println("get from disk:", filePath)
			fileSuffix := getFileSuffix(fileName)
			w.Header().Set("Content-Type", setting.MineType[fileSuffix].(string))
			w.Write(data)
		}
		return
	}

	// 指定返回头中的disposition-content,让浏览器以附件的形式下载文件

	fileSuffix := getFileSuffix(fileName)
	w.Header().Set("Content-Type", setting.MineType[fileSuffix].(string))
	// w.Header().Set("content-disposition", "attachment;filename="+fileName)
	w.Write(data)
}

/*
获取文件，并将文件发送给浏览器
具体功能:
redis中存在就从redis中加载,redis中不存在就从硬盘加载，并将内容加载到redis中
*/
func getFile(fileName string, filePath string) (data []byte, err error) {

	// 尝试从redis中获取数据
	// 判断key是否存在于redis中
	exist, err := rdb.Exists(context.Background(), fileName).Result()
	if err != nil {
		myLog.errorLogger.Println("getFile() err:", err)
	}

	// 文件key存在
	if exist == 1 {
		// 先给文件的access++
		err = increseAccess(fileName)
		if err != nil {
			myLog.errorLogger.Panicln("getFile() err:", err)
		}

		// 判断文件是否是热点数据，是否需要延长其存活时间
		if isHotkey(fileName) {
			// 是热点数据，延长其存活时间并返回数据
			data, err = getFileFromRedis(fileName)
			if err != nil {
				myLog.errorLogger.Println("getFile() err:", err)
				return data, err
			}

			err = setTTL(fileName, time.Duration(setting.HotTTL)*time.Minute)
			if err != nil {
				myLog.errorLogger.Println("getFile() err:", err)
				return data, err
			}
			myLog.dailyLogger.Printf("file %v has extended its ttl\n", fileName)
			myLog.dailyLogger.Println("get from redis:", filePath)
			return
		}

		// 判断文件是否需要缓存
		if isLoadToRedis(fileName) { //>5
			// 文件访问数达到6，说明还没缓存但是需要缓存
			if getFileAccess(fileName) == int64(setting.LoadCount+1) {
				data, err = getFileStream(filePath)
				// 从硬盘获取文件错误，没有数据返回
				if err != nil {
					myLog.errorLogger.Println("getFile() err:", err)
					return nil, err
				}
				// 将获得的字节流加载到redis中
				err = loadFileToRedis(fileName, data)
				if err != nil {
					myLog.errorLogger.Printf("loadFileToRedis err:%v\n", err)
					myLog.dailyLogger.Println("get from disk:" /*, filePath*/)
					return data, err
				} else if opErr, ok := err.(*net.OpError); ok {
					if opErr.Timeout() {
						myLog.errorLogger.Printf("getFile() tiomeout operation:%v\n", opErr.Op)
					} else {
						myLog.errorLogger.Printf("getFile() err operation:%v\n", opErr.Op)
					}
					myLog.dailyLogger.Println("get from disk:" /*filePath*/)
					return data, err
				}
				return
			} else { // 这些是已经缓存了的但是还没被延长ttl的文件
				data, err = getFileFromRedis(fileName)
				if err != nil {
					myLog.errorLogger.Println("getFile() err:", err)
					return
				}
				myLog.dailyLogger.Println("get from redis:", filePath)
				return
			}
		}

		// 不需要缓存，从硬盘加载后返回即可
		data, err = getFileStream(filePath)
		if err != nil {
			myLog.errorLogger.Println("getFile() err:", err)
			return nil, err
		}
		myLog.dailyLogger.Println("get from disk:" /*filePath*/)
		return
	}

	// 如果程序运行到这里，说明内存没有命中，那么从硬盘中加载

	// 创建文件的key,设置文件的access

	// 获得文件的字节流
	data, err = getFileStream(filePath)
	if err != nil {
		myLog.errorLogger.Printf("getFile() err:%v\n", err)
		return
	}

	// 将数据返回并且创建这个key的access
	err = loadAccessToRedis(fileName, 1)
	if err != nil {
		myLog.errorLogger.Printf("getFile() err:%v\n", err)
		return data, err
	}
	myLog.dailyLogger.Println("get from disk:" /*filePath*/)
	return
}

// 获取文件的字节流,相当于从硬盘加载数据
func getFileStream(filePath string) (fileStream []byte, err error) {

	// 打开一个文件
	file, err := os.Open(filePath)
	if err != nil {
		// fmt.Println("err:", err)
		myLog.errorLogger.Printf("getFileStream() open file err:%v", err)
		return
	}
	// 延迟关闭，避免内存泄露
	defer file.Close()

	// 采用bufio读取文件，提升效率
	reader := bufio.NewReader(file)
	buf := make([]byte, 50)
	fileStream = make([]byte, 0)

	// 循环读取文件，以字节流的形式
	for {
		size, err := reader.Read(buf)
		// 当读到文件末尾时
		if size == 0 || err == io.EOF {
			break
		}
		fileStream = append(fileStream, buf[:size]...)
	}

	return fileStream, nil
}

// 从redis获取文件
func getFileFromRedis(key string) (result []byte, err error) {
	str, err := rdb.HGet(context.Background(), key, "data").Result()
	if err != nil {
		// 结果为空，redis中不存在数据
		if err == redis.Nil {
			// fmt.Println("文件不在内存中")
			myLog.errorLogger.Printf("getFileFromRedis() err:%v\n", err)
			return
		} else {
			myLog.errorLogger.Printf("getFileFromRedis() err:%v\n", err)
			return
		}
	}
	return []byte(str), err
}

// 获取文件的访问次数
func getFileAccess(key string) int64 {
	access, err := rdb.HGet(context.Background(), key, "access").Result()
	if err != nil {
		myLog.errorLogger.Println("getFileAccess() err:", err)
		return -1
	}
	accessNum, _ := strconv.ParseInt(access, 10, 64)
	return accessNum
}

// 将文件加载至redis中
func loadFileToRedis(key string, fileStream []byte) (err error) {
	err = rdb.HSet(context.Background(), key, "data", fileStream).Err()
	if err != nil {
		myLog.errorLogger.Printf("loadFileToRedis() err:%v\n", err)
		return
	}
	// 设置其ttl
	err = setTTL(key, time.Duration(setting.TTL)*time.Minute)
	if err != nil {
		myLog.errorLogger.Printf("loadFileToRedis() err:%v\n", err)
		return
	}
	return
}

// 将文件的访问次数加载至redis中
func loadAccessToRedis(key string, accessNum int) (err error) {
	err = rdb.HSet(context.Background(), key, "access", accessNum).Err()
	if err != nil {
		myLog.errorLogger.Printf("loadAccessToRedis() err:%v\n", err)
		return
	}
	// 设置其ttl
	err = setTTL(key, time.Duration(setting.TTL)*time.Minute)
	if err != nil {
		myLog.errorLogger.Printf("loadAccessToRedis() err:%v\n", err)
		return
	}
	return
}

// 缓存策略，通过策略判断这个数据是否需要加入到缓存
func isLoadToRedis(key string) bool {
	accessNum := getFileAccess(key)
	// 如果大于阈值，则加载到redis中
	if accessNum > int64(setting.LoadCount) {
		return true
	} else if accessNum == -1 {
		myLog.errorLogger.Println("isHotKey() err:key don't exist in redis")
		return false
	} else {
		return false
	}

}

// 缓存策略，判断这个key是否要延长其ttl
func isHotkey(key string) bool {
	accessNum := getFileAccess(key)
	return accessNum > int64(setting.ExtendCount)
}

// 设置key的ttl
func setTTL(key string, time time.Duration) error {
	_, err := rdb.Expire(context.Background(), key, time).Result()
	if err != nil {
		myLog.errorLogger.Println("extendTTL() err:", err)
		return err
	}
	return err
}

// 使文件访问次数自增一
func increseAccess(key string) (err error) {
	_, err = rdb.HIncrBy(context.Background(), key, "access", 1).Result()
	if err != nil {
		myLog.errorLogger.Println("increseAccess() err:", err)
		return
	}
	return nil
}

// 预加载，手动选择一些热数据加载至redis中4
// func preload() {

// }

// 获取文件的后缀返回
func getFileSuffix(file string) (suffix string) {
	fileSlice := strings.Split(file, ".")
	return fileSlice[len(fileSlice)-1]
}
