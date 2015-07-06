package githubprovider

import (
	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
)

type Config struct {
	UserKey         string
	OrganizationKey string
}

type Clients struct {
	OrgClient  *github.Client
	UserClient *github.Client
}

// Client returns  clients for accessing github.
func (c *Config) Clients() (*Clients, error) {
	orgClient := github.NewClient(
		oauth2.NewClient(
			oauth2.NoContext,
			oauth2.StaticTokenSource(
				&oauth2.Token{
					AccessToken: c.OrganizationKey,
				},
			),
		),
	)
	userClient := github.NewClient(
		oauth2.NewClient(
			oauth2.NoContext,
			oauth2.StaticTokenSource(
				&oauth2.Token{
					AccessToken: c.UserKey,
				},
			),
		),
	)
	return &Clients{
		OrgClient:  orgClient,
		UserClient: userClient,
	}, nil

}
