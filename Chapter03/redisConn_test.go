package main

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"redisCache/Chapter03/model"
	"redisCache/redisConn"
	"redisCache/utils"
	"testing"
	"time"
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

	// 有序集合-范围型
	t.Run("ZSets-Range", func(t *testing.T) {
		z1 := Client.Conn.ZAdd("zset-1",
			&redis.Z{Member: "a", Score: 1},
			&redis.Z{Member: "b", Score: 2},
			&redis.Z{Member: "c", Score: 3},
		)
		t.Log("the hash is: ", z1)

		z2 := Client.Conn.ZAdd("zset-2",
			&redis.Z{Member: "b", Score: 4},
			&redis.Z{Member: "c", Score: 1},
			&redis.Z{Member: "d", Score: 0},
		)
		t.Log("the hash is: ", z2)
		// 类似于集合的交集运算
		z3 := Client.Conn.ZInterStore("zset-i", &redis.ZStore{Keys: []string{"zset-1", "zset-2"}})
		t.Log("the hash is: ", z3)

		z4 := Client.Conn.ZRange("zset-i", 0, -1)
		t.Log("the hash is: ", z4)

		// 并集运算
		z5 := Client.Conn.ZUnionStore("zset-u", &redis.ZStore{Keys: []string{"zset-1", "zset-2"}})
		t.Log("the hash is: ", z5)

		z6 := Client.Conn.ZRange("zset-u", 0, -1)
		t.Log("the hash is: ", z6)

		// 可以把集合作为输入传给ZINTERSTORE和ZUNIONSTORE，命令会将集合看作是成团分值全为1的有序集合来处理
		z7 := Client.Conn.SAdd("set-1", "a", "d")
		t.Log("the hash is: ", z7)

		z8 := Client.Conn.ZUnionStore("zset-u2", &redis.ZStore{
			Keys: []string{"zset-1", "zset-2", "set-1"},
		})
		t.Log("the hash is: ", z8)

		z9 := Client.Conn.ZRange("zset-u2", 0, -1)
		t.Log("the hash is: ", z9)
	})

	// 发布与订阅
	t.Run("Publish and subscribe", func(t *testing.T) {
		go Client.RunPubsub1()
		go Client.RunPubsub2()
		Client.Publisher(6)
		time.Sleep(1 * time.Second)
		defer Client.Reset()
	})

	// stream练习--生产消息
	t.Run("Stream Test", func(t *testing.T) {
		s1 := Client.Conn.XAdd(&redis.XAddArgs{
			Stream: "my_streams_topic",
			MaxLen: 500,
			Values: map[string]interface{}{
				"location_id":          222,
				"detection":            "{\"det_prob\":0.857006907463074,\"x\":987,\"y\":673,\"w\":73,\"h\":87}",
				"feature_b64":          "IzGSsLeXj6y1qhotLKuTLrascZU3r6mrbqw/IpesJ623MNAmvi/IpSuoq7FisA6uyy9zMXMpEqpBqC2wFKVEsFgrIrBho5YtOqxDK1AwmCNQnv4bVSi2oCqriKlgqt+psqumKkQq363+p7QkLikQrB0coSq7odUnyi2NKHwutqi7LV2tEKjiLOCqoaB8sH+tHa0eIkuwzS9YLWAlfbAfrnCpOS2/qPYkyyAcK70iXKbKME0yHC6vpt6wHq44pQeqW7I6LUOplK5JqIOgjCmWMJYyBK/8qmCySC7wKRUt+q22KJYkeiQCs2EocDAqIp2eoTDVsPerNyQAl2EjSSziLw==",
				"image_url":            "http://192.168.7.222:9081/153,030cd1d23b8dd1.jpg",
				"image_url2":           "",
				"fit":                  222,
				"source_id":            222,
				"parent_info_id":       "00000000-0000-0000-0000-000000000000",
				"similarity_threshold": 0.54435,
				"video_url":            "rtsp://admin:yisa123456@192.168.7.211",
				"quality_int":          44862,
				"rgb_liveness":         2,
				"info_id":              "6530f424-4873-bea5-69f2-f1e368657c07",
				"mask_id":              0,
				"mask_id_prob":         100,
			},
		})
		t.Log("the hash is: ", s1)
	})
	// stream练习--消费消息
	t.Run("Stream Test", func(t *testing.T) {
		datas := make([]map[string]interface{}, 0)
		s2 := Client.Conn.XRead(&redis.XReadArgs{
			Streams: []string{"my_streams_topic", "0"},
			// 要获取的消息数量
			Count: 0,
			Block: time.Second * 5,
		})
		if len(s2.Val()) == 0 {
			// 处理错误
			panic(s2)
		}
		// 处理读取到的消息
		for _, message := range s2.Val() {
			// 处理每个流的消息
			// streamName := message.Stream
			for _, xMessage := range message.Messages {
				// 处理每条消息
				messageID := xMessage.ID
				messageData := xMessage.Values
				temp := map[string]interface{}{
					"messageID": messageID,
					"send_time": messageData["send_time"],
					"fit":       messageData["fit"],
				}
				datas = append(datas, temp)
				// 在这里处理消息的数据
				fmt.Println(messageID)
				fmt.Println(messageData)
			}
		}
		fmt.Println(datas)
		t.Log("the hash is: ", s2)
	})

	// stream练习--创建消费者组
	t.Run("Stream Group Test", func(t *testing.T) {
		// start：指定消费者组从 Stream 的哪个位置开始消费消息。这通常是一个消息的 ID。0表示从Stream开始的位置
		s3 := Client.Conn.XGroupCreate("my_streams_topic", "my_group1", "0")
		if s3.Err() != nil {
			// 处理错误
			panic(s3)
		}
	})

	// stream练习--基于消费者组消费消息
	t.Run("Stream Group Consumer Test", func(t *testing.T) {
		for {
			s4 := Client.Conn.XReadGroup(&redis.XReadGroupArgs{
				Group:    "my_group1",
				Consumer: "my_consumer1",
				Streams:  []string{"my_streams_topic", "0"},
				Count:    1,
				Block:    0,
				// 表示是否要在接收消息后发送确认（ACK）给 Redis，这里设置为 false，表示消费者会发送 ACK。
				NoAck: false,
			})
			fmt.Println(s4.Val())

			// 发送 ACK,之后从streams中移除该消息
			ackIDs := make([]string, len(s4.Val()))
			for i, stream := range s4.Val() {
				ackIDs[i] = stream.Messages[0].ID
			}
			Client.Conn.XAck("my_streams_topic", "my_group1", ackIDs...)
		}

	})
}
