package github_adduser

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GITHUB_USERNAME", nil),
				Description: "A registered Github username.",
			},

			"userKey": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GITHUB_USERKEY", nil),
				Description: "The token key for user operations.",
			},

			"organizationKey": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GITHUB_ORGANIZATIONKEY", nil),
				Description: "The token key for organization operations.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"github_adduser_record": resourceGithubAddUserRecord(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Username:        d.Get("username").(string),
		UserKey:         d.Get("userKey").(string),
		OrganizationKey: d.Get("organizationKey").(string),
	}

	return config.Client()
}
