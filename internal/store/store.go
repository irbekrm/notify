package store

import (
	"context"
	"sync"

	"github.com/irbekrm/notify/internal/repo"
)

type RWIssuerTimer interface {
	Issuer
	Timer
}

type Issuer interface {
	ReadIssues(context.Context, repo.Repository) ([]repo.Issue, error)
	WriteIssue(context.Context, repo.Issue) error
}

type Timer interface {
	ReadTime(context.Context, repo.Repository) (string, bool, error)
	WriteTime(context.Context, string, repo.Repository) error
}

var dbConnPool sync.Pool
