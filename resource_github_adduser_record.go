package github_adduser

import (
	"fmt"

	"github.com/google/go-github/github"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/koding/logging"
	"golang.org/x/oauth2"
)

// required field are here for adding a user to the organization
func resourceGithubAddUserRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceGithubAddUserRecordCreate,
		Read:   resourceGithubAddUserRecordRead,
		Update: resourceGithubAddUserRecordUpdate,
		Delete: resourceGithubAddUserRecordDelete,

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
			// role must be selected as member or admin in Membership struct
			"role": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"source_repos": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			// organization is the name of the organization
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

// resourceGithubAddUserRecordCreate adds the user to the organization & the teams
func resourceGithubAddUserRecordCreate(d *schema.ResourceData, meta interface{}) error {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: d.Get("organizationKey").(string),
		},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	user := d.Get("username").(string)
	userKey := d.Get("userKey").(string)
	org := d.Get("organization").(string)
	teamNames := d.Get("teams").([]string)

	teamIDs, err := GetTeamIDs(client, org, teamNames)

	for _, teamID := range teamIDs {
		fmt.Println("teamID-->", teamID)

		_, _, err := client.Organizations.AddTeamMembership(teamID, user)
		if err != nil {
			return err
		}
	}

	active := "active"
	member := "member"
	membership := &github.Membership{
		// state should be active to add the user into organization
		State: &active,
		Role:  &member,
	}

	ts = oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: userKey,
		},
	)
	tc = oauth2.NewClient(oauth2.NoContext, ts)
	client = github.NewClient(tc)

	_, _, err = client.Organizations.EditOrgMembership(org, membership)
	if err != nil {
		return err
	}

	teamIDs, err = GetTeamIDs(client, org, teamNames)
	for _, teamID := range teamIDs {
		_, _, err := client.Organizations.AddTeamMembership(teamID, user)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceGithubAddUserRecordRead(d *schema.ResourceData, meta interface{}) error {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: d.Get("organizationKey").(string),
		},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	// org := d.Get("organization").(string)

	active := "active"
	opt := &github.ListOrgMembershipsOptions{
		State: active,
	}

	_, _, err := client.Organizations.ListOrgMemberships(opt)
	if err != nil {
		return err
	}

	return nil
}

// Do we need this update function?
// Create func is doing same work what update does.
func resourceGithubAddUserRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	// client := meta.(*github.Client)

	return nil
}

// resourceGithubAddUserRecordCreate removes the user from the organization & the teams
func resourceGithubAddUserRecordDelete(d *schema.ResourceData, meta interface{}) error {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: d.Get("organizationKey").(string),
		},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	user := d.Get("username").(string)
	org := d.Get("organization").(string)

	// Removing a user from this list will remove them from all teams and
	// they will no longer have any access to the organization’s repositories.
	_, err = client.Organizations.RemoveMember(org, user)
	//
	if err != nil {
		logging.Error("err while removing github user from the organization: %s", err.Error())
		// we return nil here. when user removed before
		// we are gonna get err, but we dont want err
		// because of removing is done successfully.
		return nil
	}

	return nil
}
