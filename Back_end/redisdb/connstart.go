package redisdb

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var RedisClient *redis.Client


func OpenRedis () *redis.Client {

	conn := redis.NewClient (&redis.Options{
		Addr: viper.GetString("REDIS_CLIENT_NAME"),
		Password: viper.GetString("REDIS_PASSWORD"),
		DB : 0 , 
	})

	ping , err := conn.Ping(context.Background()).Result()

	if err != nil {
		log.Fatal("could not connect to redis : " , err)
	}

	log.Println("connection established :" , ping)

	RedisClient = conn

  return RedisClient
} 
