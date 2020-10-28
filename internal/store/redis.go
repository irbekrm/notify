package store

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/irbekrm/notify/internal/repo"

	redis "github.com/go-redis/redis/v8"
)

const (
	ISSUESKEYTEMPLATE    string = `%s-issues`
	STARTTIMEKEYTEMPLATE string = `%s-starttime`
)

type Redis struct{}

func NewRedisStore(ctx context.Context, addr, passwd string) (WriterFinder, error) {
	dbConnPool = sync.Pool{
		New: func() interface{} {
			return redis.NewClient(&redis.Options{Addr: addr, Password: passwd})
		},
	}
	rdb := dbConnPool.Get().(*redis.Client)
	defer dbConnPool.Put(rdb)
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed connecting to redis: %v", err)
	}
	log.Printf("connected to redis at %s", addr)
	return &Redis{}, nil
}

func (r Redis) FindIssue(ctx context.Context, issue repo.Issue, rp repo.Repository) (bool, error) {
	rdb := dbConnPool.Get().(*redis.Client)
	defer dbConnPool.Put(rdb)
	key := repoIssuesKey(rp)
	return rdb.SIsMember(ctx, key, issue.Number()).Result()
}

func (r Redis) WriteIssue(ctx context.Context, issue repo.Issue, rp repo.Repository) error {
	rdb := dbConnPool.Get().(*redis.Client)
	defer dbConnPool.Put(rdb)
	key := repoIssuesKey(rp)
	_, err := rdb.SAdd(ctx, key, issue.Number()).Result()
	return err
}

func (r Redis) FindTime(ctx context.Context, rp repo.Repository) (string, bool, error) {
	rdb := dbConnPool.Get().(*redis.Client)
	defer dbConnPool.Put(rdb)
	key := repoStartTimeKey(rp)
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	return val, err == nil, err
}

func (r Redis) WriteTime(ctx context.Context, t string, rp repo.Repository) error {
	rdb := dbConnPool.Get().(*redis.Client)
	defer dbConnPool.Put(rdb)
	key := repoStartTimeKey(rp)
	return rdb.Set(ctx, key, t, 0).Err()
}

func repoStartTimeKey(rp repo.Repository) string {
	return fmt.Sprintf(STARTTIMEKEYTEMPLATE, rp.String())
}

func repoIssuesKey(rp repo.Repository) string {
	return fmt.Sprintf(ISSUESKEYTEMPLATE, rp.String())
}
