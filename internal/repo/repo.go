package repo

import (
	"fmt"
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
