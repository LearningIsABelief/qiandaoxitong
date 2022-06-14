package util

import (
	"github.com/spf13/viper"
	"qiandao/store"
	"time"
)

// RedisSet
// @Description: 设置key到redis
// @Author YangXuZheng
// @Date: 2022-06-13 13:46
func RedisSet(key, value string) error {
	return store.RedisDB.Self.Set(key, value, time.Duration(viper.GetInt("jwt.token-validity-second"))).Err()
}

// RedisGet
// @Description: 从redis获取一个key
// @Author YangXuZheng
// @Date: 2022-06-13 15:31
func RedisGet(key string) (string, error) {
	return store.RedisDB.Self.Get(key).Result()
}

// RedisDel
// @Description: 从redis删除一个key
// @Author YangXuZheng
// @Date: 2022-06-13 15:50
func RedisDel(key string) {
	store.RedisDB.Self.Del(key)
}

// RedisScanKeys
// @Description: 正则匹配key
// @Author YangXuZheng
// @Date: 2022-06-13 16:54
func RedisScanKeys(regular string) ([]string, error) {
	return store.RedisDB.Self.Keys(regular).Result()
}
