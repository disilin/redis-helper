package main

import (
	"common"
	"common/sdkredis"
	"log"
	"os"
	"path/filepath"
	"redis_helper/models"
	"strings"
	"time"
)

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func main() {
	configPath := getParentDirectory(getCurrentDirectory()) + "/config/redis_helper.json"

	// 读取配置文件
	config := common.NewConfig(configPath)

	log.Println("task rule config:", config)

	// redis 初始化
	sdkredis.RedisInit(config)

	// 初始化task规则
	models.InitRule(config.Rule)

	// 开始时间
	startTime := time.Now().Unix()
	// redis操作任务
	models.SartTask()
	// 结束时间
	endTime := time.Now().Unix()
	// 总耗时（单位：秒）
	log.Printf("this task used time, %d seconds.\n", endTime-startTime)
}
