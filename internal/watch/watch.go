package watch

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/irbekrm/notify/internal/receiver"
	"github.com/irbekrm/notify/internal/repo"
	"github.com/jbeda/go-wait"
)

const DEFAULTINTERVAL time.Duration = time.Minute * 5

type Client struct {
	startTime time.Time
	seen      map[int]bool
	interval  time.Duration
	repo      repo.Repository
	reciever  receiver.Notifier
}

func NewClient(repo repo.Repository, reciever receiver.Notifier, opts ...option) *Client {
	c := &Client{
		startTime: time.Now(),
		interval:  DEFAULTINTERVAL,
		repo:      repo,
		reciever:  reciever,
		seen:      make(map[int]bool),
	}
	c.applyOptions(opts...)
	return c
}

type Options []option

type option func(*Client)

func (c *Client) applyOptions(opts ...option) {
	for _, o := range opts {
		o(c)
	}
}

func Interval(interval time.Duration) option {
	return func(c *Client) {
		c.interval = interval
	}
}

func (c *Client) PollRepo(wg *sync.WaitGroup) {
	defer wg.Done()
	wait.Forever(c.pollRepo, c.interval)
}

func (c *Client) pollRepo() {
	log.Printf("checking repo %s for new issues...", c.repo)
	ctx := context.TODO()
	issues, err := c.repo.IssuesSince(ctx, c.startTime)
	if err != nil {
		log.Fatalf("could not retrieve issues for repo %s: %v", c.repo, err)
	}
	for _, i := range issues {
		number := i.Number()
		if !c.seen[number] {
			log.Printf("New issue: %s", i.Description())
			c.reciever.Notify(fmt.Sprintf("New issue: %s", i.Description()))
			c.seen[number] = true
		}
	}
}
