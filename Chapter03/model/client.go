package model

import (
	"github.com/go-redis/redis/v7"
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
