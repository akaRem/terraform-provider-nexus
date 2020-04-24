package nexus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRepositoryPypiGroup(t *testing.T) {
	repoName := fmt.Sprintf("test-repo-%s", acctest.RandString(10))
	repoNameProxy := fmt.Sprintf("test-repo-%s", acctest.RandString(10))
	repoNameHosted := fmt.Sprintf("test-repo-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{

		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createTfStmtForResourcePypiProxy(repoNameProxy) + createTfStmtForResourcePypiHosted(repoNameHosted) + createTfStmtForResourcePypiGroup(repoName),
				Check: resource.ComposeTestCheckFunc(
					// Base and common repo props
					// Identity fields
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.pypi_group", "id", repoName),
						resource.TestCheckResourceAttr("nexus_repository.pypi_group", "name", repoName),
						resource.TestCheckResourceAttr("nexus_repository.pypi_group", "format", "pypi"),
						resource.TestCheckResourceAttr("nexus_repository.pypi_group", "type", "group"),
					),
					// Common fields
					// Online
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.pypi_group", "online", "true"),
						// Storage
						resource.TestCheckResourceAttr("nexus_repository.pypi_group", "storage.#", "1"),
						resource.TestCheckResourceAttr("nexus_repository.pypi_group", "storage.0.blob_store_name", "default"),
						resource.TestCheckResourceAttr("nexus_repository.pypi_group", "storage.0.strict_content_type_validation", "true"),
						// FIXME: (BUG) can't set ALLOW
						// resource.TestCheckResourceAttr("nexus_repository.pypi_group", "storage.0.write_policy", "ALLOW"),
					),
					// No fields related to other repo types
					// Format
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.pypi_group", "maven.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.pypi_group", "apt.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.pypi_group", "apt_signing.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.pypi_group", "bower.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.pypi_group", "docker.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.pypi_group", "docker_proxy.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.pypi_group", "http_client.#", "0"),
					),
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.pypi_group", "negative_cache.#", "0"),
						resource.TestCheckResourceAttr("nexus_repository.pypi_group", "proxy.#", "0"),
					),
					// Fields related to this format and type
					// Format
					// - No special fields
					// Type
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("nexus_repository.pypi_group", "group.#", "1"),
						// FIXME: (BUG) Incorrect member_names state representation.
						// For some reasons, 1st ans 2nd elements in array are not stored as group.0.member_names.0, but instead they're stored
						// as group.0.member_names.2126137474 where 2126137474 is a "random" number.
						// This number changes from test run to test run.
						// It may be a pointer to int instead of int itself, but it's not clear and requires additional research.
						// resource.TestCheckResourceAttr("nexus_repository.pypi_group", "group.0.member_names.2126137474", memberRepoName),

					),
				),
			},
			{
				ResourceName:      "nexus_repository.pypi_group",
				ImportStateId:     repoName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func createTfStmtForResourcePypiGroup(name string) string {
	return fmt.Sprintf(`
resource "nexus_repository" "pypi_group" {
	name   = "%s"
	format = "pypi"
	type   = "group"

	group {
		member_names = [ nexus_repository.pypi_proxy.name, nexus_repository.pypi_hosted.name ]
	}

	storage {

	}
}`, name)
}
