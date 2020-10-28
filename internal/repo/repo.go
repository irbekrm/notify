package repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/github"
)

type Repository struct {
	Name   string
	Owner  string
	Labels []string `json:",omitempty"`
}

type RepositoriesList struct {
	Repositories []Repository
}

func (r Repository) String() string {
	return fmt.Sprintf("%s/%s", r.Owner, r.Name)
}

func (r *Repository) IssuesSince(ctx context.Context, startTime time.Time) ([]Issue, error) {
	gh := github.NewClient(nil)
	result, _, err := gh.Issues.ListByRepo(ctx, r.Owner, r.Name, &github.IssueListByRepoOptions{Since: startTime, Labels: r.Labels})
	if err != nil {
		return nil, fmt.Errorf("could not list issues: %v", err)
	}
	issues := []Issue{}
	for _, i := range result {
		labels := []string{}
		for _, l := range i.Labels {
			labels = append(labels, *l.Name)
		}
		issue := Issue{
			number: i.GetNumber(),
			repo:   fmt.Sprintf("%s", r),
			url:    i.GetHTMLURL(),
			title:  i.GetTitle(),
			labels: labels,
		}
		issues = append(issues, issue)
	}
	return issues, nil
}

type Issue struct {
	number int
	repo   string
	url    string
	title  string
	labels []string
}

func (i Issue) Number() int {
	return i.number
}

func (i Issue) Description() string {
	return fmt.Sprintf("Issue #%d %q in repo %s!\nlabels: %s\nurl: %s", i.number, i.title, i.repo, strings.Join(i.labels, ", "), i.url)
}
