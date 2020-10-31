package store

//go:generate mockgen -source=db.go -destination=../../mocks/mock_db.go -package=mocks
import (
	"context"
	"sync"

	"github.com/irbekrm/notify/internal/github"
)

type DB interface {
	FWIssue
	FWTime
}

type FWIssue interface {
	FindIssue(context.Context, github.Issue, string) (bool, error)
	WriteIssue(context.Context, github.Issue, string) error
}

type FWTime interface {
	FindTime(context.Context, string) (string, bool, error)
	WriteTime(context.Context, string, string) error
}

var dbConnPool sync.Pool
