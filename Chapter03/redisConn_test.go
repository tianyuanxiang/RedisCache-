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
	t.Run("Test String", func(t *testing.T) {
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

		res4 := Client.Conn.Append("new-string-key", "hello ").String()
		utils.AssertStringResult(t, "6L", res4)

		res5 := Client.Conn.Append("new-string-key", "world!").String()
		utils.AssertStringResult(t, "6L", res5)

		// Setrange命令用指定的字符串覆盖给定key所储存的字符串值，覆盖的位置从偏移量 offset 开始(主要用来覆盖的)。
		res6 := Client.Conn.SetRange("new-string-key", 0, "fuck").String()
		utils.AssertStringResult(t, "6L", res6)
		// 得到指定位置的元素
		res7 := Client.Conn.GetRange("new-string-key", 6, 30).String()
		utils.AssertStringResult(t, "6L", res7)

		res8 := Client.Conn.SetBit("another-key", 7, 1).Val()
		utils.AssertnumResult(t, 100, res8)

		res9 := Client.Conn.Get("another-key").Val()
		utils.AssertStringResult(t, "100", res9)
	})
	// 列表
	t.Run("Test List Push and Pop", func(t *testing.T) {
		l1 := Client.Conn.RPush("list-key", "last").Val()
		utils.AssertnumResult(t, 101, l1)

		l2 := Client.Conn.LPush("list-key", "first").Val()
		utils.AssertnumResult(t, 101, l2)

		l3 := Client.Conn.LRange("list-key", 0, -1).Val()
		t.Log("the list is: ", l3)

		l4 := Client.Conn.RPush("list-key", "a", "b", "c").Val()
		utils.AssertnumResult(t, 101, l4)
	})
}
