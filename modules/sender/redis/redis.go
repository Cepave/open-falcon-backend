package redis

import (
	"github.com/Cepave/open-falcon-backend/modules/sender/g"
	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"time"
)

var ConnPool *redis.Pool

func InitConnPool() {
	redisConfig := g.Config().Redis

	ConnPool = &redis.Pool{
		MaxIdle:     redisConfig.MaxIdle,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisConfig.Addr)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: PingRedis,
	}
}

func PingRedis(c redis.Conn, t time.Time) error {
	_, err := c.Do("ping")
	if err != nil {
		log.Println("[ERROR] ping redis fail", err)
	}
	return err
}
