package sdkredis

import (
	"common"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Redis struct {
	Config    *common.Config // 配置文件
	RedisPool *redis.Pool    // 连接池
}

var RedisClient *Redis

func RedisInit(config *common.Config) {
	RedisClient = NewRedis(config)
}

func NewRedis(config *common.Config) (r *Redis) {
	r = new(Redis)
	r.Config = config
	fmt.Printf("redis host[%v]port[%v]passwd[%v]dbindex[%d]\n", r.Config.Redis.Address, r.Config.Redis.Port, r.Config.Redis.Password, r.Config.Redis.DBIndex)

	var REDIS_HOST string = r.Config.Redis.Address + ":" + r.Config.Redis.Port
	var REDIS_PWD string = r.Config.Redis.Password
	var REDIS_DB int = r.Config.Redis.DBIndex

	r.RedisPool = &redis.Pool{
		MaxIdle:     50,                // 最大的空闲连接数
		MaxActive:   100,               // 最大的激活连接数
		IdleTimeout: 180 * time.Second, // 最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
		//Wait:        true,              //是否在超过最大连接数的时候等待
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", REDIS_HOST)
			if err != nil {
				return nil, err
			}

			// 验证密码
			if _, err := conn.Do("AUTH", REDIS_PWD); err != nil {
				conn.Close()
				return nil, err
			}

			// 选择db
			if _, err := conn.Do("SELECT", REDIS_DB); err != nil {
				return nil, err
			}

			return conn, nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}

	return
}

/**
 *  从连接池获得连接
 */
func (self *Redis) GetConn() redis.Conn {
	return self.RedisPool.Get()
}

/**
 *  释放redis连接
 */
func RedisClose(conn redis.Conn) {
	conn.Close()
}

/**
 *  执行redis命令
 */
func RedisCommand(conn redis.Conn, commandName string, args ...interface{}) (reply interface{}, err error) {
	reply, err = conn.Do(commandName, args...)

	return
}

/**
*  选择数据库
 */
func RedisDBSelect(conn redis.Conn, db int) (reply interface{}, err error) {
	reply, err = RedisCommand(conn, "select", db)

	return
}

/**
 *  set key value
 */
func RedisSetKV(conn redis.Conn, key string, value string) (reply interface{}, err error) {
	reply, err = RedisCommand(conn, "set", key, value)

	return
}

/**
 *  del key
 */
func RedisDelKey(conn redis.Conn, key string) (reply interface{}, err error) {
	reply, err = RedisCommand(conn, "del", key)

	return
}

/**
 *  keys命令
 */
func RedisGetKeys(conn redis.Conn, pattern string) (strlist []string, err error) {
	reply, err := RedisCommand(conn, "keys", pattern)
	strlist, _ = redis.Strings(reply, nil)
	return
}

/**
 *	scan 命令
 */

func RedisScanKeys(conn redis.Conn, start int, pattern string, count int) (index int, keys []string, err error) {
	reply, err := RedisCommand(conn, "scan", start, "match", pattern, "count", count)
	values, err := redis.Values(reply, err)
	index, err = redis.Int(values[0], err)
	keys, err = redis.Strings(values[1], err)
	return
}

/**
 * 批量设置有效期
 */
func RedisMultiExpire(conn redis.Conn, keys []string, seconds int64) (reply interface{}, err error) {
	conn.Send("MULTI")
	for _, key := range keys {
		conn.Send("EXPIRE", key, seconds)
	}
	return conn.Do("EXEC")
}

/**
 * 批量删除
 */
func RedisMultiDelete(conn redis.Conn, keys []string) (reply interface{}, err error) {
	conn.Send("MULTI")
	for _, key := range keys {
		conn.Send("DEL", key)
	}
	return conn.Do("EXEC")
}
