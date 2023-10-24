package redisConn

import (
	"github.com/go-redis/redis/v7"
	"log"
	"redisCache/config"
)

type Client struct {
	Conn *redis.Client
}

func NewClient(conn *redis.Client) *Client {
	return &Client{Conn: conn}
}

func ConnectRedis() *redis.Client {
	conn := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	if _, err := conn.Ping().Result(); err != nil {
		log.Fatalf("Connect to redis client failed, err: %v\n", err)
	}
	return conn
}
