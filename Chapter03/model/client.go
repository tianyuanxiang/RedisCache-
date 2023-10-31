package model

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"log"
	"sync/atomic"
	"time"
)

type Client struct {
	Conn *redis.Client
}

func NewClient(conn *redis.Client) *Client {
	return &Client{Conn: conn}
}
func (r *Client) Reset() {
	r.Conn.FlushDB()
}

// 发布与订阅

func (r *Client) Publisher(n int) {
	time.Sleep(time.Second * 3)
	for n > 0 {
		r.Conn.Publish("channel", n)
		// 发布消息后进行短暂的休眠，让消息可以一条接一条地出现
		n--
	}
}

func (r *Client) RunPubsub1() {
	// 创建发布与订阅对象，并让它订阅给定的频道
	pubsub := r.Conn.Subscribe("channel")
	defer pubsub.Close()

	var count int32 = 0
	// 通过遍历pubsub的监听Channel的执行结果来监听订阅信息
	for item := range pubsub.Channel() {
		fmt.Println(item.String())
		atomic.AddInt32(&count, 1)
		fmt.Printf("我是1的%d\n", count)
		switch count {
		case 4:
			if err := pubsub.Unsubscribe("channel"); err != nil {
				log.Println("unsubscribe fail, err: ", err)
			} else {
				fmt.Println("unsubscribe success")
			}
		case 5:
			break
		default:
		}
	}
}
func (r *Client) RunPubsub2() {
	// 创建发布与订阅对象，并让它订阅给定的频道
	pubsub := r.Conn.Subscribe("channel")
	defer pubsub.Close()

	var count int32 = 0
	// 通过遍历pubsub的监听Channel的执行结果来监听订阅信息
	for item := range pubsub.Channel() {
		fmt.Println(item.String())
		atomic.AddInt32(&count, 1)
		fmt.Printf("我是2的%d\n", count)
		switch count {
		case 4:
			if err := pubsub.Unsubscribe("channel"); err != nil {
				log.Println("unsubscribe fail, err: ", err)
			} else {
				fmt.Println("unsubscribe success")
			}
		case 5:
			break
		default:
		}
	}
}
