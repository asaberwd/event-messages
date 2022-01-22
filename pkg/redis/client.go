package elasticcache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sort"
	"time"
)

var ctx = context.Background()
type Client struct {
	RedisProvider Provider
}

type Provider interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

func NewClient(redisProvider Provider) *Client{
	return &Client{RedisProvider: redisProvider}
}

func (c *Client) CreateKey(key string, value string) error{
	keyWithTime := fmt.Sprintf("%s-%v", key, time.Now().UTC().UnixNano())
	err := c.RedisProvider.Set(ctx, keyWithTime, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetKey(key string) (string, error){
	s, err := c.RedisProvider.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("%s is not exist", key)
	}
	go func() {
		c.RedisProvider.Del(ctx, s)
	}()
	return s, nil
}

func (c *Client) GetAllKeys(prefix string) ([]string, error){
	pattern := fmt.Sprintf("%s*", prefix)
	keys := make([]string, 0)
	values := make([]string, 0)
	iter := c.RedisProvider.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return values, err
	}
	sort.Strings(keys)
	for _,k := range keys {
		val, err := c.GetKey(k)
		if err != nil {
			continue
		}
		values = append(values, val)
	}
	go func() {
		c.RedisProvider.Del(ctx, keys...)
	}()
	return values, nil
}
