package main

import (
	"redisCache/redisConn"
	"testing"
)

func TestLoginCookies(t *testing.T) {
	// 获取redis配置 && 连接redis
	conn := redisConn.ConnectRedis()
	conn.Get("key")
}
