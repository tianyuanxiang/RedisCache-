package model

import (
	"crypto"
	"encoding/hex"
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"log"
	"net/url"
	"redisCache/Chapter02/common"
	"redisCache/Chapter02/repository"
	"redisCache/utils"
	"strings"
	"sync/atomic"
	"time"
)

type Client struct {
	Conn *redis.Client
}

func NewClient(conn *redis.Client) *Client {
	return &Client{Conn: conn}
}

/*
* @Desc:登录和cookie缓存
* @Author: tianyx
* @Date:   2023/10/23 17:20
 */

func (r *Client) UpdateToken(token, user, item string) {
	// 获取当前时间戳
	timestamp := time.Now().Unix()
	// 记录令牌与登录用户的关系
	r.Conn.HSet("login:", token, user)
	// 记录令牌最后一次出现的时间(最后一条浏览信息)
	r.Conn.ZAdd("recent:", &redis.Z{Member: token, Score: float64(timestamp)})
	if item != "" {
		// 记录用户浏览过的商品
		r.Conn.ZAdd("viewed:"+token, &redis.Z{Score: float64(timestamp), Member: token})
		// 只保留最近25条记录，-26 表示倒数第26个成员
		r.Conn.ZRemRangeByRank("viewed"+token, 0, -26)
	}
}

func (r *Client) CleanSessions() {
	for !common.QUIT {
		size := r.Conn.ZCard("recent").Val()
		if size <= common.LIMIT {
			time.Sleep(1 * time.Second)
			continue
		}
		// 获取需要移除的令牌ID(最多100个)
		endIndex := utils.Min(size-common.LIMIT, 100)
		tokens := r.Conn.ZRange("recent", 0, endIndex-1).Val()

		var sessionKey []string
		// 给被删除的令牌构建键名
		for _, token := range tokens {
			sessionKey = append(sessionKey, "viewed:"+token)
		}
		r.Conn.Del(sessionKey...)
		// 删除令牌对应的用户信息
		r.Conn.HDel("login", tokens...)
		// 删除最近浏览记录
		r.Conn.ZRem("recent", tokens)
	}
}

/*
* @Desc:使用Redis实现购物车
* @Author: tianyx
* @Date:   2023/10/23 17:20
 */

func (r *Client) AddToCart(session, item string, count int) {
	if count <= 0 {
		// 移除指定商品
		r.Conn.HDel("cart:"+session, item)
	} else {
		r.Conn.HSet("cart:"+session, item, count)
	}
}

func (r *Client) CleanFullSessions() {
	for !common.QUIT {
		size := r.Conn.ZCard("recent").Val()
		if size <= common.LIMIT {
			time.Sleep(1 * time.Second)
			continue
		}
		endIndex := utils.Min(size-common.LIMIT, 100)
		tokens := r.Conn.ZRange("recent", 0, endIndex-1).Val()

		var sessionKey []string
		for _, token := range tokens {
			sessionKey = append(sessionKey, "viewed:"+token)
			// 删除旧会话对应用户的购物车
			sessionKey = append(sessionKey, "cart:"+token)
		}
		r.Conn.Del(sessionKey...)
		r.Conn.HDel("login", tokens...)
		r.Conn.ZRem("recent", tokens)
	}
	defer atomic.AddInt32(&common.FLAG, -1)
}

/*
* @Title:网页缓存
* @Desc:把网页上不经常加载(刷新)的内容缓存下来
* @Author: tianyx
* @Date:   2023/10/23 17:20
 */
func (r *Client) CacheRequest(request string, callback func(string) string) string {
	// 不能被缓存的请求，直接调用回调函数(itemId 为空或请求是动态的，则不应该被缓存)
	// 回调函数可以被用来判断是否存在缓存内容，如果存在缓存，可以直接返回缓存结果；
	if r.CanCache(request) {
		return callback(request)
	}
	// 如果不存在，则进行相应的处理
	page_key := "cache:" + hashRequest(request)
	content := r.Conn.Get(page_key).Val()
	// 如果页面还没有缓存
	if content == "" {
		content = callback(request)
		r.Conn.Set(page_key, content, time.Second*300)
	}
	return content
}

func (r *Client) CanCache(request string) bool {
	itemId := extractItemId(request)
	if itemId == "" || isDynamic(request) {
		return false
	}
	rank := r.Conn.ZRank("viewed:", itemId).Val()
	return rank != 0 && rank < 10000
}

func extractItemId(request string) string {
	parsed, _ := url.Parse(request)
	queryValue, _ := url.ParseQuery(parsed.RawQuery)
	query := queryValue.Get("item")
	return query
}

func isDynamic(request string) bool {
	parsed, _ := url.Parse(request)
	queryValue, _ := url.ParseQuery(parsed.RawQuery)
	for _, v := range queryValue {
		for _, j := range v {
			if strings.Contains(j, "_") {
				return false
			}
		}
	}
	return true
}

func hashRequest(request string) string {
	hash := crypto.MD5.New()
	hash.Write([]byte(request))
	res := hash.Sum(nil)
	return hex.EncodeToString(res)
}

/*
* @Title:数据行缓存
* @Author: tianyx
* @Date:   2023/10/23 17:20
 */
func (r *Client) ScheduleRowCache(row_id string, delay int64) {
	// 数据行的延迟值
	r.Conn.ZAdd("delay", &redis.Z{Member: row_id, Score: float64(delay)})
	// 对需要缓存的数据进行调度
	r.Conn.ZAdd("schedule", &redis.Z{Member: row_id, Score: float64(time.Now().Unix())})
}

func (r *Client) CacheRows() {
	for !common.QUIT {
		next := r.Conn.ZRangeWithScores("schedule:", 0, 0).Val()
		now := time.Now().Unix()
		if len(next) == 0 || next[0].Score > float64(now) {
			time.Sleep(50 * time.Millisecond)
			continue
		}
		row_id := next[0].Member.(string)
		delay := r.Conn.ZScore("delay", row_id).Val()
		if delay > 0 {
			r.Conn.ZRem("delay:", row_id)
			r.Conn.ZRem("schedule:", row_id)
			r.Conn.Del("inv:" + row_id)
			continue
		}
		row := repository.Get(row_id)
		r.Conn.ZAdd("schedule:", &redis.Z{Member: row_id, Score: float64(now) + delay})
		jsonRow, err := json.Marshal(row)
		if err != nil {
			log.Fatalf("marshal json failed, data is: %v, err is: %v\n", row, err)
		}
		r.Conn.Set("inv:"+row_id, jsonRow, 0)
	}
	defer atomic.AddInt32(&common.FLAG, -1)
}

/*
* @Title:网页分析
* @Author: tianyx
* @Date:   2023/10/24 17:20
 */

func (r *Client) RescaleViewed() {
	for !common.QUIT {
		r.Conn.ZRemRangeByRank("viewed:", 0, -20001)
		// 将viewed:集合中的元素的分值乘以0.5
		r.Conn.ZInterStore("viewed:", &redis.ZStore{Weights: []float64{0.5}, Keys: []string{"viewed:"}})
		time.Sleep(5 * time.Minute)
	}
}
