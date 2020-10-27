package store

import (
	"sync"

	"github.com/irbekrm/notify/internal/repo"
)

type RWIssuerTimer interface {
	Issuer
	Timer
}

type Issuer interface {
	ReadIssues() ([]repo.Issue, error)
	WriteIssue(repo.Issue) error
}

type Timer interface {
	ReadTime(string) (string, error)
	WriteTime(string, string) error
}

var dbConnPool sync.Pool
