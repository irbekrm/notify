package store

import (
	"context"
	"sync"

	"github.com/irbekrm/notify/internal/github"
)

type WriterFinder interface {
	Issuer
	Timer
}

type Issuer interface {
	FindIssue(context.Context, github.Issue, string) (bool, error)
	WriteIssue(context.Context, github.Issue, string) error
}

type Timer interface {
	FindTime(context.Context, string) (string, bool, error)
	WriteTime(context.Context, string, string) error
}

var dbConnPool sync.Pool
