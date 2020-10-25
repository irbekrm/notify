package repo

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/go-github/v32/github"
	"github.com/irbekrm/notify/internal/receiver"
	wait "github.com/jbeda/go-wait"
)

const defaultInterval time.Duration = time.Minute * 5

type Client struct {
	startTime time.Time
	seen      map[int64]bool
	interval  time.Duration
	repo      string
	owner     string
	reciever  receiver.Notifier
	ghClient  *github.Client
}

func NewClient(repo string, owner string, reciever receiver.Notifier) *Client {
	return &Client{
		startTime: time.Now(),
		interval:  defaultInterval,
		repo:      repo,
		owner:     owner,
		reciever:  reciever,
		seen:      make(map[int64]bool),
		ghClient:  github.NewClient(nil),
	}
}

func (c *Client) WatchAndTell(wg *sync.WaitGroup) {
	defer wg.Done()
	wait.Forever(c.watchAndTellFunc, c.interval)
}

func (c *Client) watchAndTellFunc() {
	ctx := context.TODO()
	issues, _, err := c.ghClient.Issues.ListByRepo(ctx, c.owner, c.repo, &github.IssueListByRepoOptions{Since: c.startTime})
	if err != nil {
		log.Fatalf("could not list issues: %v", err)
	}
	for _, i := range issues {
		id := i.GetID()
		if !c.seen[id] {
			issue := issue{id: id, repo: c.repo, owner: c.owner, url: i.GetHTMLURL(), title: i.GetTitle()}
			c.reciever.Notify(fmt.Sprintf("%s", issue))
			c.seen[id] = true
		}
	}
}

type issue struct {
	id    int64
	repo  string
	owner string
	url   string
	title string
}

func (i issue) String() string {
	return fmt.Sprintf("New issue in repo %s/%s: %s\n Id: %d, link %s", i.owner, i.repo, i.title, i.id, i.url)
}
