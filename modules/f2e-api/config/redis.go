package config

import (
	"errors"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

type RedisConn struct {
	Conn *redis.Client
	// default extend bucket key of external alarms
	DefualtBucketKey string
	Enable           bool
}

var myReids RedisConn

func GetRedisConn() RedisConn {
	return myReids
}

func SetBucketKey(key string) {
	myReids.DefualtBucketKey = key
}

func StartRedis(enable bool, addr string, password string, bucketkey string) {
	myReids = RedisConn{}
	myReids.Enable = enable
	if enable {
		myReids.Conn = redis.NewClient(&redis.Options{
			Addr: addr,
			// if keep black meams no password set
			Password: password,
			// use default DB
			DB: 0,
		})
		myReids.DefualtBucketKey = bucketkey
	}
}

func PutAlarmToRedis(body string) error {
	client := GetRedisConn()
	if client.Enable {
		conn := client.Conn
		bucket := client.DefualtBucketKey
		size, err := conn.LPush(client.DefualtBucketKey, body).Result()
		if err != nil {
			log.Errorf("got error during lpush data to redis: %v", err.Error())
		} else {
			log.Debugf("redis push to %s, send size: %v", bucket, size)
		}
		return err
	} else {
		err := errors.New("redis is not enabled. please ask system admin for help")
		log.Error(err.Error())
		return err
	}
}
