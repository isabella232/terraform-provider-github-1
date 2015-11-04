package githubprovider

import (
	"os"
	"testing"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func TestGithubAddAndRemoveUser(t *testing.T) {
	orgClient := github.NewClient(
		oauth2.NewClient(
			oauth2.NoContext,
			oauth2.StaticTokenSource(
				&oauth2.Token{
					AccessToken: os.Getenv("GITHUB_ORG_ACCESS_TOKEN"),
				},
			),
		),
	)
	userClient := github.NewClient(
		oauth2.NewClient(
			oauth2.NoContext,
			oauth2.StaticTokenSource(
				&oauth2.Token{
					AccessToken: os.Getenv("GITHUB_USER_ACCESS_TOKEN"),
				},
			),
		),
	)

	org := "organizasyon"
	user := "mehmetalikoding"
	teamNames := []string{"team"}
	role := "member"
	active := "active"

	teamIDs, err := GetTeamIDs(orgClient, org, teamNames)

	optAddOrgMembership := &github.OrganizationAddTeamMembershipOptions{
		Role: role,
	}

	for _, teamID := range teamIDs {
		_, _, err := orgClient.Organizations.AddTeamMembership(teamID, user, optAddOrgMembership)
		if err != nil {
			t.Fail()
			// return
		}
	}

	// When you try to add a user to a team twice,
	// it doesn't return any error
	for _, teamID := range teamIDs {
		_, _, err := orgClient.Organizations.AddTeamMembership(teamID, user, optAddOrgMembership)
		if err != nil {
			t.Fail()
			return
		}
	}

	membership := &github.Membership{
		// state should be active to add the user into organization
		State: &active,

		// Role is the required for the membership
		Role: &role,
	}

	_, _, err = userClient.Organizations.EditOrgMembership("", org, membership)
	if err != nil {
		t.Fail()
		return
	}

	// When you try to add a user to an organization twice,
	// it doesn't return any error
	_, _, err = userClient.Organizations.EditOrgMembership("", org, membership)
	if err != nil {
		t.Fail()
		return
	}

	// You can remove a member from organization
	_, err = orgClient.Organizations.RemoveMember(org, user)
	if err != nil {
		t.Fail()
		return
	}

}
