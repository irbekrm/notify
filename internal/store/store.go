package store

import (
	"context"
	"sync"

	"github.com/irbekrm/notify/internal/repo"
)

type WriterFinder interface {
	Issuer
	Timer
}

type Issuer interface {
	FindIssue(context.Context, repo.Issue, repo.Repository) (bool, error)
	WriteIssue(context.Context, repo.Issue, repo.Repository) error
}

type Timer interface {
	FindTime(context.Context, repo.Repository) (string, bool, error)
	WriteTime(context.Context, string, repo.Repository) error
}

var dbConnPool sync.Pool
