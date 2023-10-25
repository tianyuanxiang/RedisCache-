package main

import (
	"redisCache/Chapter03/model"
	"redisCache/redisConn"
	"redisCache/utils"
	"testing"
)

func TestLoginCookies(t *testing.T) {
	// 获取redis配置 && 连接redis
	conn := redisConn.ConnectRedis()
	Client := model.NewClient(conn)
	t.Run("Test INCR and DECR", func(t *testing.T) {
		Client.Conn.Get("key")
		// res := Client.Conn.Incr("key").Val()
		// utils.AssertnumResult(t, 3, res)

		res := Client.Conn.IncrBy("key", 15).Val()
		utils.AssertnumResult(t, 18, res)

		res = Client.Conn.DecrBy("key", 5).Val()
		utils.AssertnumResult(t, 13, res)

		res2, _ := Client.Conn.Get("key").Int64()
		utils.AssertnumResult(t, 28, res2)

		res3 := Client.Conn.Set("key", "13", 0).Val()
		utils.AssertStringResult(t, "13", res3)
	})

}
