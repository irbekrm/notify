package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v32/github"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Repository struct {
	Name  string
	Owner string
}

type RepositoriesList struct {
	Repositories []Repository
}

func main() {
	var configPath *string = flag.String("configpath", "", "path to directory with config.yaml")
	flag.Parse()
	viper.SetConfigName("config")
	viper.AddConfigPath(*configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("viper failed to read config: %v", err)
	}

	r := RepositoriesList{}
	if err := viper.Unmarshal(&r); err != nil {
		log.Fatalf("viper failed to unmarshal config: %v", err)
	}
	fmt.Printf("%#v", r)
	ctx := context.TODO()
	client := github.NewClient(nil)
	issues, _, err := client.Issues.ListByRepo(ctx, "irbekrm", "cidrgo", &github.IssueListByRepoOptions{})
	if err != nil {
		log.Fatalf("could not list issues: %v", err)
	}
	for _, i := range issues {
		fmt.Printf("Issue: %s ", *i.Title)
	}
}
