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
	seen      map[int]bool
	interval  time.Duration
	repo      string
	owner     string
	reciever  receiver.Notifier
	ghClient  *github.Client
}

func NewClient(repo string, owner string, reciever receiver.Notifier, opts ...option) *Client {
	c := &Client{
		startTime: time.Now(),
		interval:  defaultInterval,
		repo:      repo,
		owner:     owner,
		reciever:  reciever,
		seen:      make(map[int]bool),
		ghClient:  github.NewClient(nil),
	}
	c.applyOptions(opts...)
	return c
}

type option func(*Client)

func (c *Client) applyOptions(opts ...option) {
	for _, o := range opts {
		o(c)
	}
}

func WatchInterval(interval time.Duration) option {
	return func(c *Client) {
		c.interval = interval
	}
}

func (c *Client) WatchAndTell(wg *sync.WaitGroup) {
	defer wg.Done()
	wait.Forever(c.watchAndTellFunc, c.interval)
}

func (c *Client) watchAndTellFunc() {
	log.Printf("checking repo %s/%s for new issues...", c.owner, c.repo)
	ctx := context.TODO()
	issues, _, err := c.ghClient.Issues.ListByRepo(ctx, c.owner, c.repo, &github.IssueListByRepoOptions{Since: c.startTime})
	if err != nil {
		log.Fatalf("could not list issues: %v", err)
	}
	for _, i := range issues {
		number := i.GetNumber()
		if !c.seen[number] {
			log.Printf("found new issue (#%d) in %s/%s", number, c.owner, c.repo)
			issue := issue{number: number, repo: c.repo, owner: c.owner, url: i.GetHTMLURL(), title: i.GetTitle()}
			c.reciever.Notify(fmt.Sprintf("%s", issue))
			c.seen[number] = true
		}
	}
}

type issue struct {
	number int
	repo   string
	owner  string
	url    string
	title  string
}

func (i issue) String() string {
	return fmt.Sprintf("New issue in repo %s/%s!\nIssue #%d: %s\n %s", i.owner, i.repo, i.number, i.title, i.url)
}
