package nexus

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccRepositoryAptHosted(t *testing.T) {
	t.Parallel()
	repoName := fmt.Sprintf("test-repo-%s", acctest.RandString(10))
	repoAptDistribution := "bionic"
	repoAptSigningKeypair := acctest.RandString(10)
	repoAptSigningPassphrase := acctest.RandString(10)
	repoCleanupPolicyNames := []string{"weekly-cleanup"}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRepositoryAptHosted(repoName, repoAptDistribution, repoAptSigningKeypair, repoAptSigningPassphrase, repoCleanupPolicyNames),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nexus_repository_apt_hosted.apt_hosted", "name", repoName),
				),
			},
			{
				ResourceName:      "nexus_repository_apt_hosted.apt_hosted",
				ImportStateId:     repoName,
				ImportState:       true,
				ImportStateVerify: true,
				// apt_signing not returned by API
				ImportStateVerifyIgnore: []string{"apt_signing"},
			},
		},
	})
}
func testAccRepositoryAptHosted(name string, aptDistribution string, aptSigningKEypair string, aptSigningPassphrase string, cleanupPolicyNames []string) string {
	return fmt.Sprintf(`
resource "nexus_repository_apt_hosted" "apt_hosted" {
	name   = "%s"
    apt {
		distribution = "%s"
	}
    apt_signing {
		keypair    = "%s"
		passphrase = "%s"
	}
	storage {
        write_policy = "ALLOW"
	}
}
`, name, aptDistribution, aptSigningKEypair, aptSigningPassphrase)
}
