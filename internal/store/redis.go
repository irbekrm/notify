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

// func (r Redis) ApplyStartTimeForRepo(ctx context.Context, t time.Time, r repo.RepositoriesList) error {
// 	rdb := dbConnPool.Get().(*redis.Client)
// 	defer dbConnPool.Put(rdb)
// 	err := rdb.Set(ctx)
// }

func (r Redis) ReadIssues() ([]repo.Issue, error) {
	return []repo.Issue{}, nil
}

func (r Redis) WriteIssue(repo.Issue) error {
	return nil
}

func (r Redis) ReadTime(s string) (string, error) {
	return "", nil
}

func (r Redis) WriteTime(t, repo string) error {
	return nil
}
