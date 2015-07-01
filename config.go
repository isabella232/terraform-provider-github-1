package github_adduser

import (
	"log"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
)

type Config struct {
	Username        string
	UserKey         string
	OrganizationKey string
}

// Client() returns a new client for accessing cloudflare.
func (c *Config) Client() (*github.Client, error) {
	// client, err := github.NewClient(c.Username, c.UserKey)
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: c.OrganizationKey,
		},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	log.Printf("[INFO] Github Client configured for user: %s", c.Username)

	return client, nil
}
