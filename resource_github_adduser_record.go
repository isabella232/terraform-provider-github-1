package githubprovider

import (
	"github.com/google/go-github/github"
	"github.com/hashicorp/terraform/helper/schema"
)

// required field are here for adding a user to the organization
func resourceGithubAddUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceGithubAddUserCreate,
		Read:   resourceGithubAddUserRead,
		Update: resourceGithubAddUserCreate,
		Delete: resourceGithubAddUserDelete,

		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"role": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "member",
			},

			"organization": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"teams": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

// GetTeamIDs gets the teams id of the organization
func GetTeamIDs(client *github.Client, org string, teamNames []string) ([]int, error) {
	currentPage := 1

	var teamIDs []int

	for {
		options := &github.ListOptions{
			PerPage: 100,
			Page:    currentPage,
		}

		teams, resp, err := client.Organizations.ListTeams(org, options)
		if err != nil {
			return nil, err
		}

		if len(teams) == 0 {
			break
		}
		// Iterate over all teams and add current user to realted team(s)
		for i, team := range teams {
			for _, teamName := range teamNames {
				if *team.Name == teamName {
					teamIDs = append(teamIDs, *teams[i].ID)
				}
			}
		}

		if currentPage == resp.LastPage {
			break
		}

		currentPage = resp.NextPage
	}

	return teamIDs, nil
}

// resourceGithubAddUserCreate adds the user to the organization & the teams
func resourceGithubAddUserCreate(d *schema.ResourceData, meta interface{}) error {
	clientOrg := meta.(*Clients).OrgClient
	org := d.Get("organization").(string)
	teamNames := interfaceToStringSlice(d.Get("teams"))

	user := d.Get("username").(string)

	teamIDs, err := GetTeamIDs(clientOrg, org, teamNames)

	for _, teamID := range teamIDs {
		_, _, err := clientOrg.Organizations.AddTeamMembership(teamID, user)
		if err != nil {
			return err
		}
	}

	active := "active"
	role := d.Get("role").(string)

	membership := &github.Membership{
		// state should be active to add the user into organization
		State: &active,
		Role:  &role,
	}

	client := meta.(*Clients).UserClient

	_, _, err = client.Organizations.EditOrgMembership(org, membership)
	return err
}

func resourceGithubAddUserRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// resourceGithubAddUserCreate removes the user from the organization & the teams
func resourceGithubAddUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Clients).OrgClient

	user := d.Get("username").(string)
	org := d.Get("organization").(string)

	// Removing a user from this list will remove them from all teams and
	// they will no longer have any access to the organization’s repositories.
	_, err := client.Organizations.RemoveMember(org, user)
	return err
}
