package github

//go:generate mockgen -source=finder.go -destination=../../mocks/mock_finder.go -package=mocks
import (
	"context"
	"time"
)

type Finder interface {
	Find(context.Context, time.Time) ([]Issue, error)
	RepoName() string
}
