package store

import (
	"github.com/go-redis/redis"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
)

type Redis struct {
	Self *redis.Client
}

var RedisDB *Redis

func (redisDB *Redis) Init() {
	RedisDB = &Redis{
		Self: GetRedisDB(),
	}
}

func (redisDB *Redis) Close() {
	RedisDB.Self.Close()
}

func GetRedisDB() *redis.Client {
	return openRedis(
		viper.GetString("redis.host"),
		viper.GetString("redis.password"),
		viper.GetInt("redis.database"),
	)
}

func openRedis(host, password string, database int) *redis.Client {
	redisDB := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       database,
	})
	//_, err := redisDB.Ping().Result()
	//if err != nil {
	//	log.Errorf(err, "redis连接超时. URL: %s", host)
	//	panic(err)
	//}
	log.Infof("redis连接成功. URL: %s", host)
	return redisDB
}
