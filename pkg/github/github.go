package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

// repositoryNamespace contains the owner and repo properties of a Github repository
type repositoryNamespace struct {
	owner string
	repo  string
}

// Github wraps go-github to provide convenient Github API access
type Github struct {
	client *github.Client
	ctx    context.Context
	repo   repositoryNamespace
}

func (g *Github) String() string {
	return fmt.Sprintf("Github client for %s/%s", g.repo.owner, g.repo.repo)
}

func (g *Github) GetBranch(name string) (*github.Branch, error) {
	branch, _, err := g.client.Repositories.GetBranch(g.ctx, g.repo.owner, g.repo.repo, name)
	return branch, err
}
