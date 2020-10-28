package watch

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/irbekrm/notify/internal/receiver"
	"github.com/irbekrm/notify/internal/repo"
	"github.com/irbekrm/notify/internal/store"
	"github.com/jbeda/go-wait"
)

const (
	DEFAULTINTERVAL time.Duration = time.Minute * 5
	TIMEFORMAT                    = time.RFC3339
	TIMEROUND                     = time.Minute
)

type Client struct {
	startTime StartTime
	interval  time.Duration
	repo      repo.Repository
	reciever  receiver.Notifier
	db        store.WriterFinder
}

func NewClient(ctx context.Context, rp repo.Repository, reciever receiver.Notifier, db store.WriterFinder, opts ...option) (*Client, error) {
	var st StartTime
	timeString, exists, err := db.FindTime(ctx, rp)
	// attempt to write start time to db even if failed reading it before
	if !exists || err != nil {
		st = StartTime{t: time.Now()}
		err := db.WriteTime(ctx, fmt.Sprintf("%s", st), rp)
		if err != nil {
			// If we cannot connect to database we probably dont' want to continue
			return nil, fmt.Errorf("failed writing start time to database: %v", err)
		}
	} else { // start time found in the database
		st, err = parseTime(timeString)
		if err != nil {
			return nil, fmt.Errorf("failed parsing start time: %v", err)
		}
	}

	c := &Client{
		startTime: st,
		interval:  DEFAULTINTERVAL,
		repo:      rp,
		reciever:  reciever,
		db:        db,
	}
	c.applyOptions(opts...)
	return c, nil
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

func (c *Client) PollRepo(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	f := c.pollRepoWithContextFunc(ctx)
	wait.Forever(f, c.interval)
}

func (c *Client) pollRepoWithContextFunc(ctx context.Context) func() {
	return func() {
		log.Printf("checking repo %s for new issues since %s...", c.repo, c.startTime)
		issues, err := c.repo.IssuesSince(ctx, c.startTime.t)
		if err != nil {
			log.Fatalf("could not retrieve issues for repo %s: %v", c.repo, err)
		}
		for _, i := range issues {
			issueExists, err := c.db.FindIssue(ctx, i, c.repo)
			if err != nil {
				log.Printf("could not check if issue exists in database: %v", err)
			}
			// notify about new issue even in case of db error
			if !issueExists || err != nil {
				log.Printf("New issue: %s", i.Description())
				c.reciever.Notify(fmt.Sprintf("New issue: %s", i.Description()))
				err := c.db.WriteIssue(ctx, i, c.repo)
				if err != nil {
					log.Printf("could not write issue to database: %v", err)
				}
			}
		}
	}
}

type StartTime struct {
	t time.Time
}

func (s StartTime) String() string {
	t := s.t.Round(TIMEROUND).Format(TIMEFORMAT)
	return fmt.Sprintf("%s", t)
}

func parseTime(s string) (StartTime, error) {
	t, err := time.Parse(TIMEFORMAT, s)
	if err != nil {
		return StartTime{}, fmt.Errorf("failed parsing %s as time: %v", s, err)
	}
	return StartTime{t: t}, nil
}
