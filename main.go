package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/irbekrm/notify/internal/github"
	"github.com/irbekrm/notify/internal/receiver"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	var configPath *string = flag.String("configpath", "", "path to directory with config.yaml")
	var webhookUrl *string = flag.String("webhookurl", "", "incoming webhook url for Slack notifications backend")
	var interval *time.Duration = flag.Duration("interval", 0, fmt.Sprintf(`Custom polling interval in format that would be accepted by time.ParseDuration (i.e 1m3s, 1h etc). Default: %v`, github.DEFAULTWATCHINTERVAL))
	flag.Parse()
	viper.SetConfigName("config")
	viper.AddConfigPath(*configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("viper failed to read config: %v", err)
	}

	rl := github.RepositoriesList{}
	if err := viper.Unmarshal(&rl); err != nil {
		log.Fatalf("viper failed to unmarshal config: %v", err)
	}

	rec, err := receiver.NewSlackReceiver(*webhookUrl)
	if err != nil {
		log.Fatalf("failed creating new receiver: %v", err)
	}

	var opts github.Options
	if f := flag.Lookup("interval"); f != nil && f.Changed {
		opts = append(opts, github.WatchInterval(*interval))
	}

	wg := &sync.WaitGroup{}
	for _, r := range rl.Repositories {
		client := github.NewClient(r, rec, opts...)
		wg.Add(1)
		go client.WatchAndTell(wg)
	}
	wg.Wait()
}
