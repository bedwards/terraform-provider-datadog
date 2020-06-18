package datadog

import (
	"context"
	"fmt"
	"testing"

	datadogV2 "github.com/DataDog/datadog-api-client-go/api/v2/datadog"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

const archiveAzureConfigForCreation = `
resource "datadog_logs_archive" "my_azure_archive" {
	name = "my first azure archive"
	query = "service:toto"
	azure = {
		container 		= "my-container"
		client_id 		= "aaaaaaaa-1a1a-1a1a-1a1a-aaaaaaaaaaab"
		tenant_id       = "aaaaaaaa-1a1a-1a1a-1a1a-aaaaaaaaaaaa"
		storage_account = "storageAccount"
		region          = "my-region"
		path            = "/path/blou"
	}
}
`

var archiveAzure = datadogV2.LogsArchiveCreateRequest{
	Data: &datadogV2.LogsArchiveCreateRequestDefinition{
		Attributes: &datadogV2.LogsArchiveCreateRequestAttributes{
			Destination: datadogV2.LogsArchiveCreateRequestDestination{
				LogsArchiveDestinationAzure: &datadogV2.LogsArchiveDestinationAzure{
					Container: "my-container",
					Integration: datadogV2.LogsArchiveIntegrationAzure{
						ClientId: "aaaaaaaa-1a1a-1a1a-1a1a-aaaaaaaaaaab",
						TenantId: "aaaaaaaa-1a1a-1a1a-1a1a-aaaaaaaaaaaa",
					},
					Path:           datadogV2.PtrString("/path/blou"),
					Region:         datadogV2.PtrString("my-region"),
					StorageAccount: "storageAccount",
					Type:           "azure",
				},
			},
			Name:  "my first azure archive",
			Query: "service:toto",
		},
		Type: "archives",
	},
}

const archiveGCSConfigForCreation = `
resource "datadog_logs_archive" "my_gcs_archive" {
	name = "my first gcs archive"
	query = "service:tata"
	gcs = {
        bucket 		 = "dd-logs-test-datadog-api-client-go"
        path 	     = "/path/blah"
        client_email = "email@email.com"
        project_id   = "aaaaaaaa-1a1a-1a1a-1a1a-aaaaaaaaaaaa"
	}
}
`

var archiveGCS = datadogV2.LogsArchiveCreateRequest{
	Data: &datadogV2.LogsArchiveCreateRequestDefinition{
		Attributes: &datadogV2.LogsArchiveCreateRequestAttributes{
			Destination: datadogV2.LogsArchiveCreateRequestDestination{
				LogsArchiveDestinationGCS: &datadogV2.LogsArchiveDestinationGCS{
					Integration: datadogV2.LogsArchiveIntegrationGCS{
						ClientEmail: "email@email.com",
						ProjectId: "aaaaaaaa-1a1a-1a1a-1a1a-aaaaaaaaaaaa",
					},
					Path:           datadogV2.PtrString("/path/blah"),
					Bucket:         "dd-logs-test-datadog-api-client-go",
					Type:           "gcs",
				},
			},
			Name:  "my first gcs archive",
			Query: "service:tata",
		},
		Type: "archives",
	},
}

const archiveS3ConfigForCreation = `
resource "datadog_logs_archive" "my_s3_archive" {
	name = "my first azure archive"
	query = "service:toto"
	s3 = {
        bucket 		 = "bucket"
        path 		 = "/path/hello"
        client_email = "clientEmail"
        project_id   = "projectId"
        account_id   = "accountId"
        role_name    = "roleName"
	}
}
`

//Test
// create: OK azure
func TestAccDatadogLogsArchiveAzure_basic(t *testing.T) {
	defer gock.Disable()
	archiveType := "azure"
	expectedOut := readFixture(t, fmt.Sprintf("fixtures/logs/archives/%s/create.json", archiveType))
	gock.New("https://api.datadoghq.com").Post("/api/v2/logs/config/archives").MatchType("json").JSON(archiveAzure).Reply(200).Type("json").BodyString(expectedOut)
	id := "FooBar"
	byIdURL := fmt.Sprintf("/api/v2/logs/config/archives/%s", id)
	gock.New("https://api.datadoghq.com").Get(byIdURL).Reply(200).Type("json").BodyString(expectedOut)
	gock.New("https://api.datadoghq.com").Get(byIdURL).Reply(200).Type("json").BodyString(expectedOut)
	gock.New("https://api.datadoghq.com").Get(byIdURL).Reply(200).Type("json").BodyString(expectedOut)
	gock.New("https://api.datadoghq.com").Get(byIdURL).Reply(404).Type("json").BodyString(expectedOut)
	accProviders := testAccProvidersWithHttpClient(t, http.DefaultClient)
	accProvider := testAccProvider(t, accProviders)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    accProviders,
		CheckDestroy: testAccCheckArchiveDestroy(accProvider),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
				},
				Config: archiveAzureConfigForCreation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"datadog_logs_archive.my_azure_archive", "name", "my first azure archive"),
					resource.TestCheckResourceAttr(
						"datadog_logs_archive.my_azure_archive", "query", "service:toto"),
					resource.TestCheckResourceAttr(
						"datadog_logs_archive.my_azure_archive", "azure.container", "my-container"),
					resource.TestCheckResourceAttr(
						"datadog_logs_archive.my_azure_archive", "azure.client_id", "aaaaaaaa-1a1a-1a1a-1a1a-aaaaaaaaaaab"),
					resource.TestCheckResourceAttr(
						"datadog_logs_archive.my_azure_archive", "azure.tenant_id", "aaaaaaaa-1a1a-1a1a-1a1a-aaaaaaaaaaaa"),
					resource.TestCheckResourceAttr(
						"datadog_logs_archive.my_azure_archive", "azure.storage_account", "storageAccount"),
					resource.TestCheckResourceAttr(
						"datadog_logs_archive.my_azure_archive", "azure.path", "/path/blou"),
					resource.TestCheckResourceAttr(
						"datadog_logs_archive.my_azure_archive", "azure.region", "my-region"),
					resource.TestCheckResourceAttr(
						"datadog_logs_archive.my_azure_archive", "id", id),
				),
			},
		},
	})
}

// create: Ok gcs
func TestAccDatadogLogsArchiveGCS_basic(t *testing.T) {
	defer gock.Disable()
	archiveType := "gcs"
	expectedOut := readFixture(t, fmt.Sprintf("fixtures/logs/archives/%s/create.json", archiveType))
	gock.New("https://api.datadoghq.com").Post("/api/v2/logs/config/archives").MatchType("json").JSON(archiveGCS).Reply(200).Type("json").BodyString(expectedOut)
	id := "FooBar"
	byIdURL := fmt.Sprintf("/api/v2/logs/config/archives/%s", id)
	gock.New("https://api.datadoghq.com").Get(byIdURL).Reply(200).Type("json").BodyString(expectedOut)
	gock.New("https://api.datadoghq.com").Get(byIdURL).Reply(200).Type("json").BodyString(expectedOut)
	gock.New("https://api.datadoghq.com").Get(byIdURL).Reply(200).Type("json").BodyString(expectedOut)
	gock.New("https://api.datadoghq.com").Get(byIdURL).Reply(404).Type("json").BodyString(expectedOut)
	accProviders := testAccProvidersWithHttpClient(t, http.DefaultClient)
	accProvider := testAccProvider(t, accProviders)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    accProviders,
		CheckDestroy: testAccCheckArchiveDestroy(accProvider),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
				},
				Config: archiveGCSConfigForCreation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"datadog_logs_archive.my_gcs_archive", "name", "my first gcs archive"),
					resource.TestCheckResourceAttr(
						"datadog_logs_archive.my_gcs_archive", "query", "service:tata"),
					resource.TestCheckResourceAttr(
						"datadog_logs_archive.my_gcs_archive", "gcs.bucket", "dd-logs-test-datadog-api-client-go"),
					resource.TestCheckResourceAttr(
						"datadog_logs_archive.my_gcs_archive", "gcs.client_email", "email@email.com"),
					resource.TestCheckResourceAttr(
						"datadog_logs_archive.my_gcs_archive", "gcs.project_id", "aaaaaaaa-1a1a-1a1a-1a1a-aaaaaaaaaaaa"),
					resource.TestCheckResourceAttr(
						"datadog_logs_archive.my_gcs_archive", "gcs.path", "/path/blah"),
					resource.TestCheckResourceAttr(
						"datadog_logs_archive.my_gcs_archive", "id", id),
				),
			},
		},
	})
}

// create: Ok s3
// create: type azure + azure, s3 defined => Fail
// create: type azure + gcs defined => Fail
// create: type unknown => Fail
// update: OK
// update: does not exist
// delete: OK
// delete: does not exist

func testAccCheckArchiveExists(accProvider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		providerConf := accProvider.Meta().(*ProviderConfiguration)
		datadogClientV2 := providerConf.DatadogClientV2
		authV2 := providerConf.AuthV2

		if err := archiveExistsChecker(s, authV2, datadogClientV2); err != nil {
			return err
		}
		return nil
	}
}

func archiveExistsChecker(s *terraform.State, authV2 context.Context, datadogClientV2 *datadogV2.APIClient) error {
	for _, r := range s.RootModule().Resources {
		if r.Type == "datadog_logs_archive" {
			id := r.Primary.ID
			if _, _, err := datadogClientV2.LogsArchivesApi.GetLogsArchive(authV2, id).Execute(); err != nil {
				return fmt.Errorf("received an error when retrieving archive, (%s)", err)
			}
		}
	}
	return nil
}

func testAccCheckArchiveDestroy(accProvider *schema.Provider) func(*terraform.State) error {
	return func(s *terraform.State) error {
		providerConf := accProvider.Meta().(*ProviderConfiguration)
		datadogClientV2 := providerConf.DatadogClientV2
		authV2 := providerConf.AuthV2

		if err := archiveDestroyHelper(s, authV2, datadogClientV2); err != nil {
			return err
		}
		return nil
	}
}

func archiveDestroyHelper(s *terraform.State, authV2 context.Context, datadogClientV2 *datadogV2.APIClient) error {
	for _, r := range s.RootModule().Resources {
		if r.Type == "datadog_logs_archive" {
			id := r.Primary.ID
			archive, httpresp, err := datadogClientV2.LogsArchivesApi.GetLogsArchive(authV2, id).Execute()
			if err != nil {
				if httpresp != nil && httpresp.StatusCode == 404 {
					continue
				}
				return fmt.Errorf("received an error when retrieving pipeline, (%s)", err)
			}
			if &archive != nil {
				return fmt.Errorf("archive still exists")
			}
		}

	}
	return nil
}

// readFixture opens the file at path and returns the contents as a string
func readFixture(t *testing.T, path string) string {
	t.Helper()
	fixturePath, err := filepath.Abs(path)
	if err != nil {
		t.Fatalf("failed to get fixture file path: %v", err)
	}
	data, err := ioutil.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("failed to open fixture file: %v", err)
	}
	return string(data)
}
