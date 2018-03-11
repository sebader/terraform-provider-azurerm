package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSchedulerJobCollection_basic(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_scheduler_job_collection.test"
	config := testAccAzureRMSchedulerJobCollection_basic(ri, testLocation(), "")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSchedulerJobCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMSchedulerJobCollection_complete(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_scheduler_job_collection.test"
	config := testAccAzureRMSchedulerJobCollection_complete(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSchedulerJobCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSchedulerJobCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func testCheckAzureRMSchedulerJobCollectionDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_scheduler_job_collection" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).applicationSecurityGroupsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Scheduler Job Collection still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMSchedulerJobCollectionExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %q", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Scheduler Job Collection: %q", name)
		}

		client := testAccProvider.Meta().(*ArmClient).schedulerJobCollectionsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Scheduler Job Collection %q (resource group: %q) was not found: %+v", name, resourceGroup, err)
			}

			return fmt.Errorf("Bad: Get on schedulerJobCollectionsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMSchedulerJobCollection_basic(rInt int, location string, additional string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_scheduler_job_collection" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
%s
}
`, rInt, location, rInt, additional)
}

func testAccAzureRMSchedulerJobCollection_complete(rInt int, location string) string {
	return testAccAzureRMSchedulerJobCollection_basic(rInt, location, `
  state = "disabled"
  quota {
    max_recurrence_frequency = "hour"
    max_recurrence_interval  = 10
    max_job_count            = 10
  }
`)
}
