package store

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/irbekrm/notify/internal/repo"

	redis "github.com/go-redis/redis/v8"
)

type Redis struct{}

func NewRedisStore(addr, passwd string) (RWIssuerTimer, error) {
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

func (r Redis) ReadIssues(ctx context.Context, rp repo.Repository) ([]repo.Issue, error) {
	return []repo.Issue{}, nil
}

func (r Redis) WriteIssue(ctx context.Context, issue repo.Issue) error {
	return nil
}

func (r Redis) ReadTime(ctx context.Context, rp repo.Repository) (string, bool, error) {
	rdb := dbConnPool.Get().(*redis.Client)
	defer dbConnPool.Put(rdb)
	val, err := rdb.Get(ctx, rp.Name).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	return val, err == nil, err
}

func (r Redis) WriteTime(ctx context.Context, t string, rp repo.Repository) error {
	rdb := dbConnPool.Get().(*redis.Client)
	defer dbConnPool.Put(rdb)
	return rdb.Set(ctx, rp.Name, t, 0).Err()
}
