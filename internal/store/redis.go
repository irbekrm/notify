package store

import (
	"context"
	"fmt"
	"log"
	"sync"

	redis "github.com/go-redis/redis/v8"
)

type Redis struct{}

func NewRedisStore(addr, passwd string) (Storer, error) {
	dbConnPool = sync.Pool{
		New: func() interface{} {
			return redis.NewClient(&redis.Options{Addr: addr, Password: passwd})
		},
	}
	ctx := context.TODO()
	rdb := dbConnPool.Get().(*redis.Client)
	defer dbConnPool.Put(rdb)
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed connecting to redis: %v", err)
	}
	log.Printf("connected to redis at %s", addr)
	return &Redis{}, nil
}
