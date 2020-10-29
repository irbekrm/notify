package repo

import (
	"fmt"
	"strings"
)

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
