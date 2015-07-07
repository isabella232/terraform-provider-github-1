package githubprovider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccGithubFork_Basic(t *testing.T) {
	// var githubCli github.Client

	// userKey := os.Getenv("GITHUB_USERKEY")

	var providers []*schema.Provider
	providerFactories := map[string]terraform.ResourceProviderFactory{
		"github": func() (terraform.ResourceProvider, error) {
			p := Provider()
			providers = append(providers, p.(*schema.Provider))
			return p, nil
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		Providers:         testAccProviders,
		CheckDestroy:      testAccCheckGithubForkDestroy,
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckGithubForkConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGithubForkExists("github_fork.foobar"),
					resource.TestCheckResourceAttr(
						"github_fork.foobar", "username", "mehmetalisavas"),
					resource.TestCheckResourceAttr(
						"github_fork.foobar", "organization", "organizasyon"),
				),
			},
		},
	})
}

func testAccCheckGithubForkDestroy(s *terraform.State) error {
	// client := testAccProvider.Meta().(*Clients).UserClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "github_fork" {
			continue
		}
		// TODO check if forked repo still exists
	}

	return nil
}

func testAccCheckGithubForkExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		fmt.Println("RS data -->", rs)

		if !ok {
			return fmt.Errorf("Not found: %s, res: %#v", n, s.RootModule())
		}

		return nil
	}
}

const testAccCheckGithubForkConfig_basic = `
resource "github_fork" "foobar" {
    username = "mehmetalisavas"
    repos = ["organization"]
    organization = "organizasyon"
}
`

//////////////////////////\
//
// BELOW CODE IS HARDCODED !!!!
//
//////////////////////////\

// import (
// 	"fmt"
// 	"testing"

// 	"golang.org/x/oauth2"

// 	"github.com/google/go-github/github"
// 	. "github.com/smartystreets/goconvey/convey"
// )

// func TestGithubForkRepo(t *testing.T) {
// 	Convey("while preparing requirements", t, func() {

// 		Convey("while forking repos of organization", func() {
// 			repos := "organization"
// 			org := "organizasyon"

// 			ts := oauth2.StaticTokenSource(
// 				&oauth2.Token{
// 					// fork token MAS.
// 					AccessToken: "b591b5dba904aee86c19b39361c680da2f9863b0",
// 				},
// 			)

// 			tc := oauth2.NewClient(oauth2.NoContext, ts)
// 			client := github.NewClient(tc)

// 			_, resp, err := client.Repositories.CreateFork(org, repos, nil)
// 			if err != nil {
// 				fmt.Println("response is=>", resp)
// 				fmt.Println("Err-->", err)
// 				return
// 			}

// 		})

// 	})
// }
