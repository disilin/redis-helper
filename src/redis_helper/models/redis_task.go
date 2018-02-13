package models

import (
	"common"
	"common/sdkredis"
	"log"

	"github.com/garyburd/redigo/redis"
)

var (
	TaskRule common.RuleConf
	total    int64
)

type State int

// 任务类型枚举
const (
	TASK_EXPIRE State = iota // value --> 0
	TASK_DELETE              // value --> 1
)

func InitRule(rule common.RuleConf) {
	TaskRule = rule
}

// 启动任务
func SartTask() {
	if TaskRule.Method == "expire" {
		log.Printf("=============redis task method is %s=============", TaskRule.Method)
		SetExpireTask()
		log.Printf("=============%s task method is done=============", TaskRule.Method)
	} else if TaskRule.Method == "delete" {
		log.Printf("=============redis task method is %s=============", TaskRule.Method)
		DeleteTask()
		log.Printf("=============%s task method is done=============", TaskRule.Method)
	} else {
		log.Fatalf("=============error!!! unknown redis method: %s============\n",
			TaskRule.Method)
	}
}

// 设置有效期任务
func SetExpireTask() {
	handleTask(TASK_EXPIRE)
}

// 删除任务
func DeleteTask() {
	handleTask(TASK_DELETE)
}

// 处理任务
func handleTask(taskType State) {
	redisConn := sdkredis.RedisClient.GetConn()
	defer redisConn.Close()
	for _, pattern := range TaskRule.Pattern {
		dispatchTask(redisConn, 0, pattern, TaskRule.ScanCount, taskType)
	}
}

// 任务分发
func dispatchTask(redisConn redis.Conn, start int, pattern string, count int, taskType State) (index int, keys []string, err error) {
	log.Println("--> task scan keys start:", start)
	start, keys, err = sdkredis.RedisScanKeys(redisConn, start, pattern, count)
	if err != nil {
		log.Println(err)
	}
	total += int64(len(keys))
	log.Printf("--> handle keys count: %d, handled keys total: %d\n", len(keys), total)
	// 根据task类型，迭代分发任务
	if taskType == TASK_EXPIRE {
		doExpireTask(redisConn, keys, TaskRule.ExpireTime)
	} else if taskType == TASK_DELETE {
		doDeleteTask(redisConn, keys)
	}
	if start != 0 {
		dispatchTask(redisConn, start, pattern, count, taskType)
	}
	return
}

// 执行设置有效期任务
func doExpireTask(redisConn redis.Conn, keys []string, seconds int64) {
	//todo
	sdkredis.RedisMultiExpire(redisConn, keys, seconds)
}

// 执行删除任务
func doDeleteTask(redisConn redis.Conn, keys []string) {
	//todo
	sdkredis.RedisMultiDelete(redisConn, keys)
}
