package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/github"
)

type Finder interface {
	Find(context.Context, time.Time) ([]Issue, error)
	RepoName() string
}

type ghRepo struct {
	r         Repository
	authToken string
}

func NewFinder(r Repository) Finder {
	return &ghRepo{r: r}
}

func (ghr *ghRepo) Find(ctx context.Context, startTime time.Time) ([]Issue, error) {
	ghc := github.NewClient(nil)
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

func (ghr *ghRepo) RepoName() string {
	return ghr.r.String()
}
