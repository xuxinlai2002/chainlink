package soak

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
)

func TestOCRSoak(t *testing.T) {
	// Use this variable to pass in any custom EVM specific TOML values to your Chainlink nodes
	customNetworkTOML := ``
	// Uncomment below for debugging TOML issues on the node
	// network := networks.MustGetSelectedNetworksFromEnv()[0]
	// fmt.Println("Using Chainlink TOML\n---------------------")
	// fmt.Println(networks.AddNetworkDetailedConfig(config.BaseOCR1Config, customNetworkTOML, network))
	// fmt.Println("---------------------")
	config, err := tc.GetConfig("Soak", tc.OCR)
	require.NoError(t, err, "Error getting config")
	runOCRSoakTest(t, config, customNetworkTOML)
}

func TestOCRSoak_GethReorgBelowFinality_FinalityTagDisabled(t *testing.T) {
	config, err := tc.GetConfig(t.Name(), tc.OCR)
	require.NoError(t, err, "Error getting config")
	runOCRSoakTest(t, config, "")
}

func TestOCRSoak_GethReorgBelowFinality_FinalityTagEnabled(t *testing.T) {
	config, err := tc.GetConfig(t.Name(), tc.OCR)
	require.NoError(t, err, "Error getting config")
	runOCRSoakTest(t, config, "")
}

func runOCRSoakTest(t *testing.T, config tc.TestConfig, customNetworkTOML string) {
	l := logging.GetTestLogger(t)

	l.Info().Str("test", t.Name()).Msg("Starting OCR soak test")

	ocrSoakTest, err := testsetups.NewOCRSoakTest(t, &config, false)
	require.NoError(t, err, "Error creating soak test")
	if !ocrSoakTest.Interrupted() {
		ocrSoakTest.DeployEnvironment(customNetworkTOML, &config)
	}
	if ocrSoakTest.Environment().WillUseRemoteRunner() {
		return
	}
	t.Cleanup(func() {
		if err := actions_seth.TeardownRemoteSuite(ocrSoakTest.TearDownVals(t)); err != nil {
			l.Error().Err(err).Msg("Error tearing down environment")
		}
	})
	if ocrSoakTest.Interrupted() {
		err = ocrSoakTest.LoadState()
		require.NoError(t, err, "Error loading state")
		ocrSoakTest.Resume()
	} else {
		ocrSoakTest.Setup(&config)
		ocrSoakTest.Run()
	}
}
