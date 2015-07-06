package github_adduser

import (
	"github.com/google/go-github/github"
	"github.com/hashicorp/terraform/helper/schema"
	"golang.org/x/oauth2"
)

// required field are here for adding a user to the organization
func resourceGithubForkRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceGithubForkRecordCreate,
		// Read:   resourceGithubAddUserRecordRead,
		// Update: resourceGithubAddUserRecordUpdate,
		// Delete: resourceGithubAddUserRecordDelete,

		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			// userKey is the token of the authenticated user
			"userKey": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			// organizationKey is the token of the authenticated
			// user that owner or admin of organization
			"organizationKey": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"repos": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			// organization is the name of the organization
			"organization": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// interfaceToStringSlice converts the interface to slice of string
func interfaceToStringSlice(s interface{}) []string {
	slice, ok := s.([]interface{})
	if !ok {
		return nil
	}

	sslice := make([]string, len(slice))
	for i := range slice {
		sslice[i] = slice[i].(string)
	}
	return sslice
}

// resourceGithubForkRecordCreate forks the repos of the organization
func resourceGithubForkRecordCreate(d *schema.ResourceData, meta interface{}) error {
	// user := d.Get("username").(string)
	userKey := d.Get("userKey").(string)
	org := d.Get("organization").(string)
	repos := interfaceToStringSlice(d.Get("repos"))

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: userKey,
		},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	for _, repo := range repos {
		// Creates a fork for the authenticated user.
		_, _, err := client.Repositories.CreateFork(org, repo, nil)
		if err != nil {
			return err
		}
	}
	return nil

}
