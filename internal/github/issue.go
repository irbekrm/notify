package github

//go:generate mockgen -source=issue.go -destination=../../mocks/mock_issue.go -package=mocks
import (
	"fmt"
	"strings"
)

type Issue interface {
	Number() int
	Description() string
}

type issue struct {
	number int
	repo   string
	url    string
	title  string
	labels []string
}

func (i issue) Number() int {
	return i.number
}

func (i issue) Description() string {
	return fmt.Sprintf("Issue #%d %q in repo %s!\nlabels: %s\nurl: %s", i.number, i.title, i.repo, strings.Join(i.labels, ", "), i.url)
}
