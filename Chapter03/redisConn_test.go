package main

import (
	"github.com/go-redis/redis/v7"
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

		l5 := Client.Conn.LRange("list-key", 0, -1).Val()
		// utils.AssertStringResult(t, "101", l5)
		t.Log("the list is: ", l5)

		// TODO：保留start到stop的元素(-1表示最后一个),start和stop会被保留
		l6 := Client.Conn.LTrim("list-key", 4, -1)
		t.Log("the list is: ", l6)

		l7 := Client.Conn.LRange("list-key", 0, -1).Val()
		// utils.AssertStringResult(t, "101", l5)
		t.Log("the list is: ", l7)
	})

	// 阻塞式列表弹出和移动
	t.Run("Test List BPush and BPop", func(t *testing.T) {
		l1 := Client.Conn.RPush("list", "item1").Val()
		utils.AssertnumResult(t, 101, l1)
		//
		l2 := Client.Conn.RPush("list", "item2").Val()
		utils.AssertnumResult(t, 101, l2)
		//
		l3 := Client.Conn.LPush("list2", "item3").Val()
		utils.AssertnumResult(t, 101, l3)
		// 从list2列表中弹出最右端的元素，将这个元素Push到list列表的最左边
		l4 := Client.Conn.BRPopLPush("list2", "list", 1).Val()
		utils.AssertStringResult(t, "101", l4)

		l5 := Client.Conn.BRPopLPush("list2", "list", 1).Val()
		utils.AssertStringResult(t, "101", l5)
		//t.Log("the list is: ", l5)

		l6 := Client.Conn.LRange("list", 0, -1)
		t.Log("the list is: ", l6)
		//
		l7 := Client.Conn.BLPop(1, "list", "list2")
		// utils.AssertStringResult(t, "101", l7)
		t.Log("the list is: ", l7)
	})

	// 集合常用命令
	t.Run("Test Set SADD and SREM", func(t *testing.T) {
		s1 := Client.Conn.SAdd("set-key", "a", "b", "c")
		t.Log("the set is: ", s1)

		s2 := Client.Conn.SRem("set-key", "c", "d")
		t.Log("the set is: ", s2)

		s3 := Client.Conn.SCard("set-key")
		t.Log("the set is: ", s3)
		////
		s4 := Client.Conn.SMembers("set-key")
		t.Log("the set is: ", s4)

		// 将一个元素从一个集合移到另一个集合中
		s5 := Client.Conn.SMove("set-key", "set-key2", "a")
		t.Log("the set is: ", s5)

		s6 := Client.Conn.SMembers("set-key2")
		t.Log("the set is: ", s6)
	})

	// TODO:组合处理多个集合，并集、交集和差集
	t.Run("Combined Processing of Multiple Collections", func(t *testing.T) {
		s1 := Client.Conn.SAdd("skey1", "a", "b", "c", "d")
		t.Log("the set is: ", s1)

		s2 := Client.Conn.SAdd("skey2", "c", "d", "e", "f")
		t.Log("the set is: ", s2)

		// 返回差集(我有你无)元素
		s3 := Client.Conn.SDiff("skey1", "skey2")
		t.Log("the set is: ", s3)

		// 返回交集(都有)元素
		s4 := Client.Conn.SInter("skey1", "skey2")
		t.Log("the set is: ", s4)

		// 返回并集(至少有一个有)元素，两个集合里面的所有不重复的
		s5 := Client.Conn.SUnion("skey1", "skey2")
		t.Log("the set is: ", s5)
	})

	// 散列哈希-基础
	t.Run("Hash-Basics", func(t *testing.T) {
		h1 := Client.Conn.HMSet("hash-key", map[string]interface{}{
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
		})
		t.Log("the hash is: ", h1)

		h2 := Client.Conn.HMGet("hash-key", "k2", "k3")
		t.Log("the hash is: ", h2)

		h3 := Client.Conn.HLen("hash-key")
		t.Log("the hash is: ", h3)

		h4 := Client.Conn.HDel("hash-key", "k1", "k3")
		t.Log("the hash is: ", h4)
	})

	// 散列哈希-批量操作
	t.Run("Hash-batch operation", func(t *testing.T) {
		h1 := Client.Conn.HMSet("hash-key2", map[string]interface{}{
			"short": "hello",
			"long":  1000 * 1,
		})
		t.Log("the hash is: ", h1)

		// 为了避免一次取出对执行效率的影响，可以先把所有的key取出，再挨个取出值
		h2 := Client.Conn.HKeys("hash-key2")
		t.Log("the hash is: ", h2)

		// 检查给定的键是否存在于散列中
		h3 := Client.Conn.HExists("hash-key2", "num")
		t.Log("the hash is: ", h3)

		h4 := Client.Conn.HIncrBy("hash-key2", "num", 0)
		t.Log("the hash is: ", h4)

		h5 := Client.Conn.HExists("hash-key2", "num")
		t.Log("the hash is: ", h5)
	})

	// 有序集合-基础
	t.Run("Ordered Sets-Basic", func(t *testing.T) {
		z1 := Client.Conn.ZAdd("zset-key",
			&redis.Z{Member: "a", Score: 3},
			&redis.Z{Member: "b", Score: 2},
			&redis.Z{Member: "c", Score: 1},
		)
		t.Log("the hash is: ", z1)

		h2 := Client.Conn.ZCard("zset-key")
		t.Log("the hash is: ", h2)

		h3 := Client.Conn.ZIncrBy("zset-key", 3, "c")
		t.Log("the hash is: ", h3)

		h4 := Client.Conn.ZScore("zset-key", "b")
		t.Log("the hash is: ", h4)

		h5 := Client.Conn.ZRank("zset-key", "c")
		t.Log("the hash is: ", h5)

		h6 := Client.Conn.ZCount("zset-key", "0", "3")
		t.Log("the hash is: ", h6)

		h7 := Client.Conn.ZRem("zset-key", "b")
		t.Log("the hash is: ", h7)

		h8 := Client.Conn.ZRange("zset-key", 0, -1)
		t.Log("the hash is: ", h8)
	})
}
