package github

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
)

// NewGithub creates a new Github instance using a given repository URL and token
func NewGithub(remoteURL, token string) (*Github, error) {
	components, err := newFromRemoteURL(remoteURL)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	tokenClient := newStaticTokenClient(ctx, token)
	apiClient, err := newAPIClient(ctx, tokenClient, components.hostname)
	if err != nil {
		return nil, err
	}

	repo := repositoryNamespace{owner: components.owner, repo: components.repo}
	return &Github{
		client: apiClient,
		ctx:    ctx,
		repo:   repo,
	}, nil
}

func newStaticTokenClient(ctx context.Context, token string) *http.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return tc
}

func newAPIClient(ctx context.Context, tokenClient *http.Client, hostname string) (*github.Client, error) {
	if hostname == "github.com" || hostname == "www.github.com" {
		return github.NewClient(tokenClient), nil
	}
	baseURL := fmt.Sprintf("https://%s/api/v3", hostname)
	client, err := github.NewEnterpriseClient(baseURL, baseURL, tokenClient)
	if err != nil {
		return nil, err
	}
	if zen, _, err := client.Zen(ctx); err != nil || zen == "" {
		return nil, err
	}
	return client, nil
}

type remoteURLComponents struct {
	hostname string
	owner    string
	repo     string
}

func newFromRemoteURL(remoteURL string) (*remoteURLComponents, error) {
	if strings.HasPrefix(remoteURL, "git@") {
		remoteURL = strings.Replace(remoteURL, ":", "/", 1)
		remoteURL = strings.Replace(remoteURL, "git@", "git://", 1)
	}
	parsedURL, err := url.Parse(remoteURL)
	if err != nil {
		return nil, fmt.Errorf("could not extrapolate service due to unparsable URL")
	}

	hostname := parsedURL.Hostname()
	path := parsedURL.Path

	endingDotGitPattern, err := regexp.Compile(`.git$`)
	if err != nil {
		return nil, err
	}
	path = strings.Replace(path, "/", "", 1)
	path = endingDotGitPattern.ReplaceAllString(path, "")

	components := strings.Split(path, "/")
	if len(components) != 2 {
		return nil, fmt.Errorf("could not extrapolate URL values from path")
	}
	owner := components[0]
	repo := components[1]

	return &remoteURLComponents{hostname, owner, repo}, nil
}
