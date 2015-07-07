package githubprovider

import "github.com/hashicorp/terraform/helper/schema"

// required field are here for adding a user to the organization
func resourceGithubFork() *schema.Resource {
	return &schema.Resource{
		Create: resourceGithubForkCreate,
		Read:   resourceGithubForkCreate,
		Update: resourceGithubForkCreate,
		Delete: resourceGithubForkCreate,

		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
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

// resourceGithubForkCreate forks the repos of the organization
func resourceGithubForkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Clients).UserClient
	// meta.(*github.Client)
	org := d.Get("organization").(string)

	for _, repo := range interfaceToStringSlice(d.Get("repos")) {
		// Creates a fork for the authenticated user.
		_, _, err := client.Repositories.CreateFork(org, repo, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
