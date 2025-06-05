package cache

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hhr0815hhr/gint/internal/config"
	"github.com/hhr0815hhr/gint/internal/log"
)

var (
	Client redis.Cmdable
)

const (
	RedisTypeSingle  = "single"
	RedisTypeCluster = "cluster"
)

func InitializeCache() redis.Cmdable {
	if config.Conf.Redis.Type == RedisTypeSingle {
		InitClient()
	} else {
		InitCluster()
	}
	log.Logger.Println("初始化redis...success")
	return Client
}

func InitClient() {
	Client = redis.NewClient(&redis.Options{
		Addr:         config.Conf.Redis.Host,
		Password:     config.Conf.Redis.Password,
		DB:           config.Conf.Redis.Database,
		PoolSize:     config.Conf.Redis.PoolSize,
		PoolTimeout:  5 * time.Second,
		MinIdleConns: config.Conf.Redis.MinIdleConns,
		IdleTimeout:  time.Duration(config.Conf.Redis.MaxIdleTime) * time.Minute,
	})
}

func InitCluster() {
	Client = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:           config.Conf.Redis.ClusterHosts,
		Password:        config.Conf.Redis.Password, // 如果你的集群有密码
		DialTimeout:     5 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		IdleTimeout:     time.Duration(config.Conf.Redis.MaxIdleTime) * time.Minute,
		PoolSize:        config.Conf.Redis.PoolSize,
		MinIdleConns:    config.Conf.Redis.MinIdleConns,
		PoolTimeout:     5 * time.Second,
		MaxRetries:      3,
		MaxRetryBackoff: 100 * time.Millisecond,
		RouteByLatency:  false, // 可以设置为 true 以基于延迟路由命令
		RouteRandomly:   false, // 可以设置为 true 以随机路由只读命令
	})

}
