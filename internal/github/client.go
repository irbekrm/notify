package github

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type client struct {
	r         Repository
	authToken string
}

func NewGithubClient(r Repository, opts ...option) Finder {
	ghr := &client{r: r}
	ghr.applyOptions(opts...)
	return ghr
}

type Options []option

type option func(*client)

func (ghr *client) applyOptions(opts ...option) {
	for _, o := range opts {
		o(ghr)
	}
}

func AuthToken(s string) option {
	return func(ghr *client) {
		ghr.authToken = s
	}
}

func (ghr *client) Find(ctx context.Context, startTime time.Time) ([]Issue, error) {
	ghc := ghr.githubClient(ctx)
	result, _, err := ghc.Issues.ListByRepo(ctx, ghr.r.Owner, ghr.r.Name, &github.IssueListByRepoOptions{Since: startTime, Labels: ghr.r.Labels})
	if err != nil {
		return nil, fmt.Errorf("could not list issues: %v", err)
	}
	issues := []Issue{}
	for _, i := range result {
		createdAt := i.GetClosedAt()
		isInterestingUpdate := true
		number := i.GetNumber()
		if createdAt.Before(startTime) {
			isInterestingUpdate = false
			issueEvents, _, err := ghc.Issues.ListIssueEvents(ctx, ghr.r.Owner, ghr.r.Name, number, &github.ListOptions{})
			if err != nil {
				return nil, fmt.Errorf("could not list issue events: %v", err)
			}
			for _, ie := range issueEvents {
				if ie.CreatedAt.Before(startTime) {
					// old event
					continue
				}
				if *ie.Event == "labeled" {
					isInterestingUpdate = true
				}
			}
		}
		if !isInterestingUpdate {
			// the event is neither creation of a new issue nor adding of a new label
			continue
		}
		labels := []string{}
		for _, l := range i.Labels {
			labels = append(labels, *l.Name)
		}
		issue := Issue{
			number: number,
			repo:   fmt.Sprintf("%s", ghr.r),
			url:    i.GetHTMLURL(),
			title:  i.GetTitle(),
			labels: labels,
		}
		issues = append(issues, issue)
	}
	return issues, nil
}

func (ghr *client) RepoName() string {
	return ghr.r.String()
}

func (ghr *client) githubClient(ctx context.Context) *github.Client {
	var httpClient *http.Client
	if ghr.authToken != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ghr.authToken})
		httpClient = oauth2.NewClient(ctx, ts)
	}
	return github.NewClient(httpClient)
}
