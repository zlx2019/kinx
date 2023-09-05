// @Title config.go
// @Description 服务配置
// @Author Zero - 2023/9/5 15:15:32

package knet

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

const (
	// 默认加载的配置文件目录
	defaultConfigFile = "config/kinx.json"
	// 默认服务名称
	defaultName = "kinx"
	// 默认服务Host
	defaultHost = "0.0.0.0"
	// 默认服务端口
	defaultPort = 9780
)

// 配置文件路径
var confFilePath string

// configs 服务端全局配置
var configs *serverConfig

// serverConfig 服务配置属性实体
type serverConfig struct {
	// 服务名
	Name string `json:"name"`
	// 服务IP
	Host string `json:"host"`
	// 服务端口
	Port int `json:"port"`
}

// 解析命令行参数，获取服务的配置文件
func init() {
	flag.StringVar(&confFilePath, "f", defaultConfigFile, "服务端配置文件")
}

// 创建默认的服务配置
func newDefaultConfig() *serverConfig {
	return &serverConfig{
		Name: defaultName,
		Host: defaultHost,
		Port: defaultPort,
	}
}

// 加载服务的配置文件
func loadConfigs() {
	// 读取配置文件
	bytes, err := os.ReadFile(confFilePath)
	if err != nil {
		fmt.Printf("load config file: %s failed cause: %s \n", confFilePath, err.Error())
		// 读取配置文件失败，使用默认配置
		configs = newDefaultConfig()
		return
	}
	err = json.Unmarshal(bytes, &configs)
	if err != nil {
		fmt.Printf("unmarshal config file failed cause: %s \n", err.Error())
		// 解析配置文件失败，使用默认配置
		configs = newDefaultConfig()
		return
	}
	// 处理必填属性
	if len(configs.Name) == 0 {
		configs.Name = defaultName
	}
	if len(configs.Host) == 0 {
		configs.Host = defaultHost
	}
	if configs.Port == 0 {
		configs.Port = defaultPort
	}
}
