package common

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var RedisCli *redis.Client

func InitRedis() (err error) {
	host := Conf.Redis.Host
	port := Conf.Redis.Port
	password := Conf.Redis.Password

	RedisCli = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	})
	logger.Info(fmt.Sprintf("Ping Redis: %s:%s", host, port))
	_, err = RedisCli.Ping().Result()
	if err != nil {
		logger.Error(fmt.Sprintf("init redis failed:%v", err))
		return
	}
	logger.Info(fmt.Sprintf("connect to redis: %s:%s", host, port))
	return
}

// CacheOptions 获取配置
type CacheOptions struct {
	Key      string                      //缓存key
	Duration time.Duration               //缓存过期时间
	Fun      func() (interface{}, error) //自定义获取缓存结果的函数
	Receiver interface{}                 //存放获取结果
}

// GetSet 利用接口抽象获取缓存的流程
func (c *CacheOptions) GetSet() (interface{}, error) {
	_, err := GetSetCache(c)
	if err != nil {
		return nil, err
	}
	return c.Receiver, nil
}

// GetSetCache 获取缓存，不存在则调用Fun函数获取对应数据加入缓存
func GetSetCache(c *CacheOptions) (using bool, err error) {
	if c == nil || c.Receiver == nil || c.Key == "" {
		err = fmt.Errorf("illegal arguments")
		logger.Error(err.Error())
		return
	}
	//查询缓存
	val, err := RedisCli.Get(c.Key).Result()
	if err != nil && err != redis.Nil {
		logger.Error(err.Error())
		return
	}
	if err == redis.Nil {
		//调用对应函数设置并获取缓存
		c.Receiver, err = c.Fun()
		if err != nil {
			return
		}
		if fmt.Sprint(c.Receiver) == "<nil>" {
			return false, nil
		}
		logger.Debug(fmt.Sprintf("Set cache %s", c.Key))
		var buf []byte
		if data, ok := c.Receiver.([]byte); ok {
			buf = data
		} else {
			buf, err = json.Marshal(&c.Receiver)
			if err != nil {
				logger.Error(err.Error())
				return
			}
		}
		err = RedisCli.Set(c.Key, buf, c.Duration).Err()
		if err != nil {
			logger.Error(err.Error())
			return
		}
	} else {
		//如果存在则解析缓存
		using = true
		logger.Debug(fmt.Sprintf("Hit cache %s", c.Key))
		if _, ok := c.Receiver.([]byte); ok {
			c.Receiver = []byte(val)
			return
		}
		err = json.Unmarshal([]byte(val), &c.Receiver)
		if err != nil {
			logger.Error(fmt.Sprintf("解析缓存失败 key:%s value:%v", c.Key, val))
			return
		}
	}
	return
}
