package main

import (
	"log"
	"sync"

	"github.com/irbekrm/notify/internal/receiver"
	"github.com/irbekrm/notify/internal/repo"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	var configPath *string = flag.String("configpath", "", "path to directory with config.yaml")
	var webhookUrl *string = flag.String("webhookurl", "", "incoming webhook url for Slack notifications backend")
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

	wg := &sync.WaitGroup{}
	for _, r := range rl.Repositories {
		client := repo.NewClient(r, rec)
		wg.Add(1)
		go client.WatchAndTell(wg)
	}
	wg.Wait()
}
