//go:build integration

package integration

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
)

const (
	testProjectName = "Test Project"
	testProjectId   = "a1Bxx0000001XYZAAA"
	testCampaignId  = "11111111-1111-1111-1111-111111111111"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
