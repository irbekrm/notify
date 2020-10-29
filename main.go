package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/irbekrm/notify/internal/receiver"
	"github.com/irbekrm/notify/internal/repo"
	"github.com/irbekrm/notify/internal/store"
	"github.com/irbekrm/notify/internal/watch"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	var configPath *string = flag.String("configpath", "", "path to directory with config.yaml")
	var webhookUrl *string = flag.String("webhookurl", "", "incoming webhook url for Slack notifications backend")
	var interval *time.Duration = flag.Duration("interval", 0, fmt.Sprintf("Custom polling interval in format that would be accepted by time.ParseDuration (i.e 1m3s, 1h etc). Default: %v", watch.DEFAULTINTERVAL))
	flag.Parse()
	viper.SetConfigName("config")
	viper.AddConfigPath(*configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("viper failed to read config: %v", err)
	}

	rl := repo.RepositoriesList{}
	if err := viper.Unmarshal(&rl); err != nil {
		log.Fatalf("viper failed to unmarshal config: %v", err)
	}

	rec, err := receiver.NewSlackReceiver(*webhookUrl)
	if err != nil {
		log.Fatalf("failed creating new receiver: %v", err)
	}

	ctx := context.Background()

	db, err := store.NewRedisStore(ctx, "localhost:6379", "")
	if err != nil {
		log.Fatalf("failed to establish db connection: %v", err)
	}
	var opts watch.Options
	if f := flag.Lookup("interval"); f != nil && f.Changed {
		opts = append(opts, watch.Interval(*interval))
	}

	wg := &sync.WaitGroup{}
	for _, r := range rl.Repositories {
		f := repo.NewFinder(r)
		watcher, err := watch.NewClient(ctx, f, rec, db, opts...)
		if err != nil {
			log.Printf("failed creating new client for %s: %v", r, err)
			// try to continue with the other repos
			continue
		}
		wg.Add(1)
		go watcher.PollRepo(ctx, wg)
	}
	wg.Wait()
}
