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
	LoadCount    int	`json:"loadCount"`    //需要缓存的访问次数 
    ExtendCount  int   	`json:"extendCount"`  //需要延长存活时间的次数
}

/**
	".aac": "audio/aac",
    ".abw": "application/x-abiword",
    ".apng": "image/apng",
    ".arc": "application/x-freearc",
    ".avif": "image/avif",
    ".avi": "video/x-msvideo",
    ".azw": "application/vnd.amazon.ebook",
    ".bin": "application/octet-stream",
    ".bmp": "image/bmp",
    ".bz": "application/x-bzip",
    ".bz2": "application/x-bzip2",
    ".cda": "application/x-cdf",
    ".csh": "application/x-csh",
    ".css": "text/css",
    ".csv": "text/csv",
    ".doc": "application/msword",
    ".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    ".eot": "application/vnd.ms-fontobject",
    ".epub": "application/epub+zip",
    ".gz": "application/gzip",
    ".gif": "image/gif",
    ".htm, .html": "text/html",
    ".ico": "image/vnd.microsoft.icon",
    ".ics": "text/calendar",
    ".jar": "application/java-archive",
    ".jpeg, .jpg": "image/jpeg",
    ".js": "text/javascript",
    ".json": "application/json",
    ".jsonld": "application/ld+json",
    ".m3u8": "application/x-mpegURL",
    ".mid, .midi": "audio/midi",
    ".mjs": "text/javascript",
    ".mp3": "audio/mpeg",
    ".mp4": "video/mp4",
    ".mpeg": "video/mpeg",
    ".mpkg": "application/vnd.apple.installer+xml",
    ".odp": "application/vnd.oasis.opendocument.presentation",
    ".ods": "application/vnd.oasis.opendocument.spreadsheet",
    ".odt": "application/vnd.oasis.opendocument.text",
    ".oga": "audio/ogg",
    ".ogv": "video/ogg",
    ".ogx": "application/ogg",
    ".opus": "audio/opus",
    ".otf": "font/otf",
    ".png": "image/png",
    ".pdf": "application/pdf",
    ".php": "application/x-httpd-php",
    ".ppt": "application/vnd.ms-powerpoint",
    ".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
    ".rar": "application/vnd.rar",
    ".rtf": "application/rtf",
    ".sh": "application/x-sh",
    ".svg": "image/svg+xml",
    ".tar": "application/x-tar",
    ".tif, .tiff": "image/tiff",
    ".ts": "video/mp2t",
    ".ttf": "font/ttf",
    ".txt": "text/plain",
    ".vsd": "application/vnd.visio",
    ".wav": "audio/wav",
    ".weba": "audio/webm",
    ".webm": "video/webm",
    ".webp": "image/webp",
    ".woff": "font/woff",
    ".woff2": "font/woff2",
    ".xhtml": "application/xhtml+xml",
    ".xls": "application/vnd.ms-excel",
    ".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
    ".xml": "application/xml",
    ".xul": "application/vnd.mozilla.xul+xml",
    ".zip": "application/zip",
    ".3gp": "video/3gpp; audio/3gpp",
    ".3g2": "video/3gpp2; audio/3gpp2",
    ".7z": "application/x-7z-compressed"
**/

type MINETpye struct {
	Aac   string `json:".aac"`
	Abw   string `json:".abw"`
	Apng  string `json:".apng"`
	Arc   string `json:".arc"`
	Avif  string `json:".avif"`
	Avi   string `json:".avi"`
	Azw   string `json:".azw"`
	Bin   string `json:".bin"`
	Bmp   string `json:".bmp"`
	Bz    string `json:".bz"`
	Bz2   string `json:".bz2"`
	Cda   string `json:".cda"`
	Csh   string `json:".csh"`
	Css   string `json:".css"`
	Csv   string `json:".csv"`
	Doc   string `json:".doc"`
	Docx  string `json:".docx"`
	Eot   string `json:".eot"`
	Epub  string `json:".epub"`
	Gz    string `json:".gz"`
	Gif   string `json:".gif"`
	Htm   string `json:".htm, .html"`
	Html  string `json:".htm, .html"`
	Ico   string `json:".ico"`
	Ics   string `json:".ics"`
	Jar   string `json:".jar"`
	Jpeg  string `json:".jpeg, .jpg"`
	Jpg   string `json:".jpeg, .jpg"`
	Js    string `json:".js"`
	Json  string `json:".json"`
	Jsonld string `json:".jsonld"`
	M3u8  string `json:.m3u8`
	Mid   string `json:".mid, .midi"`
	Midi  string `json:".mid, .midi"`
	Mjs   string `json:".mjs"`
	Mp3   string `json:".mp3"`
	Mp4   string `json:".mp4"`
	Mpeg  string `json:".mpeg"`
	Mpkg  string `json:".mpkg"`
	Odp   string `json:".odp"`
	Ods   string `json:".ods"`
	Odt   string `json:".odt"`
	Oga   string `json:".oga"`
	Ogv   string `json:".ogv"`
	Ogx   string `json:".ogx"`
	Opus  string `json:".opus"`
	Otf   string `json:".otf"`
	Png   string `json:".png"`
	Pdf   string `json:".pdf"`
	Php   string `json:".php"`
	Ppt   string `json:".ppt"`
	Pptx  string `json:".pptx"`
	Rar   string `json:".rar"`
	Rtf   string `json:".rtf"`
	Sh    string `json:".sh"`
	Svg   string `json:".svg"`
	Tar   string `json:".tar"`
	Tif   string `json:".tif, .tiff"`
	Tiff  string `json:".tif, .tiff"`
	Ts    string `json:".ts"`
	Ttf   string `json:".ttf"`
	Txt   string `json:".txt"`
	Vsd   string `json:".vsd"`
	Wav   string `json:".wav"`
	Weba  string `json:".weba"`
	Webm  string `json:".webm"`
	Webp  string `json:".webp"`
	Woff  string `json:".woff"`
	Woff2 string `json:".woff2"`
	Xhtml string `json:".xhtml"`
	Xls   string `json:".xls"`
	Xlsx  string `json:".xlsx"`
	Xml   string `json:".xml"`
	Xul   string `json:".xul"`
	Zip   string `json:".zip"`
	Threegp string `json:".3gp"`
	Threeg2 string `json:".3g2"`
	Sevenz string `json:".7z"`
}


type Settings struct {
	Serverconfig
	RDBConfig
	MINETpye
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

	// 使用json的decoder,将config和rdbconfig文件解析成结构体
	decoder := json.NewDecoder(config)
	decoder.Decode(&setting.Serverconfig)

	decoder = json.NewDecoder(RDBConfig)
	decoder.Decode(&setting.RDBConfig)
}
