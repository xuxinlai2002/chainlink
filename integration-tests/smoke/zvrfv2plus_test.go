package smoke

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
	vrfcommon "github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/common"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrf/vrfv2plus"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestVRFv2PlusM(t *testing.T) {

	// t.Parallel()
	var (
		env                          *test_env.CLClusterTestEnv
		vrfContracts                 *vrfcommon.VRFContracts
		subIDsForCancellingAfterTest []*big.Int
		vrfKey                       *vrfcommon.VRFKeyData
		nodeTypeToNodeMap            map[vrfcommon.VRFNodeType]*vrfcommon.VRFNode
	)

	fmt.Printf("xxl 0000 env %v \n", env)
	fmt.Printf("xxl 0001 vrfContracts %v \n", vrfContracts)
	fmt.Printf("xxl 0002 subIDsForCancellingAfterTest %v \n", subIDsForCancellingAfterTest)
	fmt.Printf("xxl 0003 vrfKey %v \n", vrfKey)
	fmt.Printf("xxl 0004 nodeTypeToNodeMap %v \n", nodeTypeToNodeMap)

	//
	l := logging.GetTestLogger(t)
	//
	config, err := tc.GetConfig("Smoke", tc.VRFv2Plus)
	fmt.Printf("xxl 0005 config : %v - vrfv2plus : %v \n", config, tc.VRFv2Plus)
	require.NoError(t, err, "Error getting config")

	vrfv2PlusConfig := config.VRFv2Plus
	fmt.Printf("xxl 0006 vrfv2PlusConfig %v \n", vrfv2PlusConfig)

	chainID := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0].ChainID
	fmt.Printf("xxl 0007 chainID %v \n", chainID)

	cleanupFn := func() {

		fmt.Println("开始睡眠...")
		// 睡眠10分钟
		time.Sleep(10 * time.Minute)
		fmt.Println("睡眠结束！")

		fmt.Printf("xxl 0000 cleanupFn \n")
		sethClient, err := env.GetSethClient(chainID)
		require.NoError(t, err, "Getting Seth client shouldn't fail")

		if sethClient.Cfg.IsSimulatedNetwork() {
			fmt.Printf("xxl 0001 cleanupFn \n")
			l.Info().
				Str("Network Name", sethClient.Cfg.Network.Name).
				Msg("Network is a simulated network. Skipping fund return for Coordinator Subscriptions.")
		} else {

			fmt.Printf("xxl 0002 cleanupFn \n")
			if *vrfv2PlusConfig.General.CancelSubsAfterTestRun {
				//cancel subs and return funds to sub owner
				vrfv2plus.CancelSubsAndReturnFunds(testcontext.Get(t), vrfContracts, sethClient.MustGetRootKeyAddress().Hex(), subIDsForCancellingAfterTest, l)
			}
		}
		if !*vrfv2PlusConfig.General.UseExistingEnv {

			fmt.Printf("xxl 0003 cleanupFn \n")
			if err := env.Cleanup(test_env.CleanupOpts{TestName: t.Name()}); err != nil {
				l.Error().Err(err).Msg("Error cleaning up test environment")
			}
		}
	}

	vrfEnvConfig := vrfcommon.VRFEnvConfig{
		TestConfig: config,
		ChainID:    chainID,
		CleanupFn:  cleanupFn,
	}
	fmt.Printf("xxl 0008 vrfEnvConfig %v \n", vrfEnvConfig)

	newEnvConfig := vrfcommon.NewEnvConfig{
		NodesToCreate:                   []vrfcommon.VRFNodeType{vrfcommon.VRF},
		NumberOfTxKeysToCreate:          0,
		UseVRFOwner:                     false,
		UseTestCoordinator:              false,
		ChainlinkNodeLogScannerSettings: test_env.DefaultChainlinkNodeLogScannerSettings,
	}

	fmt.Printf("xxl 0009 vrfEnvConfig %v \n", newEnvConfig)

	env, vrfContracts, vrfKey, nodeTypeToNodeMap, err = vrfv2plus.SetupVRFV2PlusUniverse(testcontext.Get(t), t, vrfEnvConfig, newEnvConfig, l)
	fmt.Printf("xxl 0010 vrfContracts %v \n", vrfContracts)

	require.NoError(t, err, "Error setting up VRFv2Plus universe")

	fmt.Printf("xxl 0011 vrfEnvConfig %v \n", vrfKey)
	sethClient, err := env.GetSethClient(chainID)

	fmt.Printf("xxl 0012 sethClient %v \n", sethClient)

	require.NoError(t, err, "Getting Seth client shouldn't fail")
	fmt.Printf("xxl 0013 vrfEnvConfig %v \n", vrfKey)
	
	t.Run("Link Billing", func(t *testing.T) {
		configCopy := config.MustCopy().(tc.TestConfig)
		var isNativeBilling = false
		consumers, subIDsForRequestRandomness, err := vrfv2plus.SetupNewConsumersAndSubs(
			testcontext.Get(t),
			env,
			chainID,
			vrfContracts.CoordinatorV2Plus,
			configCopy,
			vrfContracts.LinkToken,
			1,
			1,
			l,
		)
		require.NoError(t, err, "error setting up new consumers and subs")
		subIDForRequestRandomness := subIDsForRequestRandomness[0]
		subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subIDForRequestRandomness)
		require.NoError(t, err, "error getting subscription information")
		vrfcommon.LogSubDetails(l, subscription, subIDForRequestRandomness.String(), vrfContracts.CoordinatorV2Plus)
		subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForRequestRandomness...)

		subBalanceBeforeRequest := subscription.Balance

		// test and assert
		_, randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
			consumers[0],
			vrfContracts.CoordinatorV2Plus,
			vrfKey,
			subIDForRequestRandomness,
			isNativeBilling,
			configCopy.VRFv2Plus.General,
			l,
			0,
		)
		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")

		require.False(t, randomWordsFulfilledEvent.OnlyPremium, "RandomWordsFulfilled Event's `OnlyPremium` field should be false")
		require.Equal(t, isNativeBilling, randomWordsFulfilledEvent.NativePayment, "RandomWordsFulfilled Event's `NativePayment` field should be false")
		require.True(t, randomWordsFulfilledEvent.Success, "RandomWordsFulfilled Event's `Success` field should be true")

		expectedSubBalanceJuels := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
		subscription, err = vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subIDForRequestRandomness)
		require.NoError(t, err, "error getting subscription information")
		subBalanceAfterRequest := subscription.Balance
		require.Equal(t, expectedSubBalanceJuels, subBalanceAfterRequest)

		status, err := consumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
		require.NoError(t, err, "error getting rand request status")
		require.True(t, status.Fulfilled)
		l.Info().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")

		require.Equal(t, *configCopy.VRFv2Plus.General.NumberOfWords, uint32(len(status.RandomWords)))
		for _, w := range status.RandomWords {
			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
		}

		fmt.Printf("xxl 0000 status.RandomWords %v \n", status.RandomWords)
	})

	//t.Run("xxl test 0000 ", func(t *testing.T) {
	//	fmt.Printf("xxl abc ... \n")
	//})

	//t.Run("Native Billing", func(t *testing.T) {
	//	configCopy := config.MustCopy().(tc.TestConfig)
	//	testConfig := configCopy.VRFv2Plus.General
	//	var isNativeBilling = true
	//
	//	consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
	//		testcontext.Get(t),
	//		env,
	//
	//		chainID,
	//		vrfContracts.CoordinatorV2Plus,
	//		configCopy,
	//		vrfContracts.LinkToken,
	//		1,
	//		1,
	//		l,
	//	)
	//	require.NoError(t, err, "error setting up new consumers and subs")
	//	subID := subIDs[0]
	//	subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
	//	require.NoError(t, err, "error getting subscription information")
	//	vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)
	//	subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)
	//
	//	subNativeTokenBalanceBeforeRequest := subscription.NativeBalance
	//
	//	// test and assert
	//	_, randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
	//		consumers[0],
	//		vrfContracts.CoordinatorV2Plus,
	//		vrfKey,
	//		subID,
	//		isNativeBilling,
	//		configCopy.VRFv2Plus.General,
	//		l,
	//		0,
	//	)
	//	require.NoError(t, err, "error requesting randomness and waiting for fulfilment")
	//	require.False(t, randomWordsFulfilledEvent.OnlyPremium)
	//	require.Equal(t, isNativeBilling, randomWordsFulfilledEvent.NativePayment)
	//	require.True(t, randomWordsFulfilledEvent.Success)
	//	expectedSubBalanceWei := new(big.Int).Sub(subNativeTokenBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
	//	subscription, err = vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
	//	require.NoError(t, err)
	//	subBalanceAfterRequest := subscription.NativeBalance
	//	require.Equal(t, expectedSubBalanceWei, subBalanceAfterRequest)
	//
	//	status, err := consumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
	//	require.NoError(t, err, "error getting rand request status")
	//	require.True(t, status.Fulfilled)
	//	l.Info().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")
	//
	//	require.Equal(t, *testConfig.NumberOfWords, uint32(len(status.RandomWords)))
	//	for _, w := range status.RandomWords {
	//		l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
	//		require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
	//	}
	//})

	//t.Run("VRF Node waits block confirmation number specified by the consumer before sending fulfilment on-chain", func(t *testing.T) {
	//	configCopy := config.MustCopy().(tc.TestConfig)
	//	testConfig := configCopy.VRFv2Plus.General
	//	var isNativeBilling = true
	//
	//	consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
	//		testcontext.Get(t),
	//		env,
	//		chainID,
	//		vrfContracts.CoordinatorV2Plus,
	//		configCopy,
	//		vrfContracts.LinkToken,
	//		1,
	//		1,
	//		l,
	//	)
	//	require.NoError(t, err, "error setting up new consumers and subs")
	//	subID := subIDs[0]
	//	subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
	//	require.NoError(t, err, "error getting subscription information")
	//	vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)
	//	subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)
	//
	//	expectedBlockNumberWait := uint16(10)
	//	testConfig.MinimumConfirmations = ptr.Ptr[uint16](expectedBlockNumberWait)
	//	randomWordsRequestedEvent, randomWordsFulfilledEvent, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
	//		consumers[0],
	//		vrfContracts.CoordinatorV2Plus,
	//		vrfKey,
	//		subID,
	//		isNativeBilling,
	//		testConfig,
	//		l,
	//		0,
	//	)
	//	require.NoError(t, err, "error requesting randomness and waiting for fulfilment")
	//
	//	// check that VRF node waited at least the number of blocks specified by the consumer in the rand request min confs field
	//	blockNumberWait := randomWordsRequestedEvent.Raw.BlockNumber - randomWordsFulfilledEvent.Raw.BlockNumber
	//	require.GreaterOrEqual(t, blockNumberWait, uint64(expectedBlockNumberWait))
	//
	//	status, err := consumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
	//	require.NoError(t, err, "error getting rand request status")
	//	require.True(t, status.Fulfilled)
	//	l.Info().Bool("Fulfilment Status", status.Fulfilled).Msg("Random Words Request Fulfilment Status")
	//})
	//t.Run("CL Node VRF Job Runs", func(t *testing.T) {
	//	configCopy := config.MustCopy().(tc.TestConfig)
	//	var isNativeBilling = false
	//	consumers, subIDsForRequestRandomness, err := vrfv2plus.SetupNewConsumersAndSubs(
	//		testcontext.Get(t),
	//		env,
	//		chainID,
	//		vrfContracts.CoordinatorV2Plus,
	//		configCopy,
	//		vrfContracts.LinkToken,
	//		1,
	//		1,
	//		l,
	//	)
	//	require.NoError(t, err, "error setting up new consumers and subs")
	//	subIDForRequestRandomness := subIDsForRequestRandomness[0]
	//	subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subIDForRequestRandomness)
	//	require.NoError(t, err, "error getting subscription information")
	//	vrfcommon.LogSubDetails(l, subscription, subIDForRequestRandomness.String(), vrfContracts.CoordinatorV2Plus)
	//	subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDsForRequestRandomness...)
	//
	//	jobRunsBeforeTest, err := nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API.MustReadRunsByJob(nodeTypeToNodeMap[vrfcommon.VRF].Job.Data.ID)
	//	require.NoError(t, err, "error reading job runs")
	//
	//	// test and assert
	//	_, _, err = vrfv2plus.RequestRandomnessAndWaitForFulfillment(
	//		consumers[0],
	//		vrfContracts.CoordinatorV2Plus,
	//		vrfKey,
	//		subIDForRequestRandomness,
	//		isNativeBilling,
	//		configCopy.VRFv2Plus.General,
	//		l,
	//		0,
	//	)
	//	require.NoError(t, err, "error requesting randomness and waiting for fulfilment")
	//
	//	jobRuns, err := nodeTypeToNodeMap[vrfcommon.VRF].CLNode.API.MustReadRunsByJob(nodeTypeToNodeMap[vrfcommon.VRF].Job.Data.ID)
	//	require.NoError(t, err, "error reading job runs")
	//	require.Equal(t, len(jobRunsBeforeTest.Data)+1, len(jobRuns.Data))
	//})
	//t.Run("Direct Funding", func(t *testing.T) {
	//	configCopy := config.MustCopy().(tc.TestConfig)
	//	wrapperContracts, wrapperSubID, err := vrfv2plus.SetupVRFV2PlusWrapperEnvironment(
	//		testcontext.Get(t),
	//		l,
	//		env,
	//		chainID,
	//		&configCopy,
	//		vrfContracts.LinkToken,
	//		vrfContracts.MockETHLINKFeed,
	//		vrfContracts.CoordinatorV2Plus,
	//		vrfKey.KeyHash,
	//		1,
	//	)
	//	require.NoError(t, err)
	//
	//	t.Run("Link Billing", func(t *testing.T) {
	//		configCopy := config.MustCopy().(tc.TestConfig)
	//		testConfig := configCopy.VRFv2Plus.General
	//		var isNativeBilling = false
	//
	//		wrapperConsumerJuelsBalanceBeforeRequest, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), wrapperContracts.LoadTestConsumers[0].Address())
	//		require.NoError(t, err, "error getting wrapper consumer balance")
	//
	//		wrapperSubscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), wrapperSubID)
	//		require.NoError(t, err, "error getting subscription information")
	//		subBalanceBeforeRequest := wrapperSubscription.Balance
	//
	//		randomWordsFulfilledEvent, err := vrfv2plus.DirectFundingRequestRandomnessAndWaitForFulfillment(
	//			wrapperContracts.LoadTestConsumers[0],
	//			vrfContracts.CoordinatorV2Plus,
	//			vrfKey,
	//			wrapperSubID,
	//			isNativeBilling,
	//			configCopy.VRFv2Plus.General,
	//			l,
	//		)
	//		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")
	//
	//		expectedSubBalanceJuels := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
	//		wrapperSubscription, err = vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), wrapperSubID)
	//		require.NoError(t, err, "error getting subscription information")
	//		subBalanceAfterRequest := wrapperSubscription.Balance
	//		require.Equal(t, expectedSubBalanceJuels, subBalanceAfterRequest)
	//
	//		consumerStatus, err := wrapperContracts.LoadTestConsumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
	//		require.NoError(t, err, "error getting rand request status")
	//		require.True(t, consumerStatus.Fulfilled)
	//
	//		expectedWrapperConsumerJuelsBalance := new(big.Int).Sub(wrapperConsumerJuelsBalanceBeforeRequest, consumerStatus.Paid)
	//
	//		wrapperConsumerJuelsBalanceAfterRequest, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), wrapperContracts.LoadTestConsumers[0].Address())
	//		require.NoError(t, err, "error getting wrapper consumer balance")
	//		require.Equal(t, expectedWrapperConsumerJuelsBalance, wrapperConsumerJuelsBalanceAfterRequest)
	//
	//		//todo: uncomment when VRF-651 will be fixed
	//		//require.Equal(t, 1, consumerStatus.Paid.Cmp(randomWordsFulfilledEvent.Payment), "Expected Consumer contract pay more than the Coordinator Sub")
	//		vrfcommon.LogFulfillmentDetailsLinkBilling(l, wrapperConsumerJuelsBalanceBeforeRequest, wrapperConsumerJuelsBalanceAfterRequest, consumerStatus, randomWordsFulfilledEvent)
	//
	//		require.Equal(t, *testConfig.NumberOfWords, uint32(len(consumerStatus.RandomWords)))
	//		for _, w := range consumerStatus.RandomWords {
	//			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
	//			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
	//		}
	//	})
	//	t.Run("Native Billing", func(t *testing.T) {
	//		configCopy := config.MustCopy().(tc.TestConfig)
	//		testConfig := configCopy.VRFv2Plus.General
	//		var isNativeBilling = true
	//
	//		wrapperConsumerBalanceBeforeRequestWei, err := sethClient.Client.BalanceAt(testcontext.Get(t), common.HexToAddress(wrapperContracts.LoadTestConsumers[0].Address()), nil)
	//		require.NoError(t, err, "error getting wrapper consumer balance")
	//
	//		wrapperSubscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), wrapperSubID)
	//		require.NoError(t, err, "error getting subscription information")
	//		subBalanceBeforeRequest := wrapperSubscription.NativeBalance
	//
	//		randomWordsFulfilledEvent, err := vrfv2plus.DirectFundingRequestRandomnessAndWaitForFulfillment(
	//			wrapperContracts.LoadTestConsumers[0],
	//			vrfContracts.CoordinatorV2Plus,
	//			vrfKey,
	//			wrapperSubID,
	//			isNativeBilling,
	//			configCopy.VRFv2Plus.General,
	//			l,
	//		)
	//		require.NoError(t, err, "error requesting randomness and waiting for fulfilment")
	//
	//		expectedSubBalanceWei := new(big.Int).Sub(subBalanceBeforeRequest, randomWordsFulfilledEvent.Payment)
	//		wrapperSubscription, err = vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), wrapperSubID)
	//		require.NoError(t, err, "error getting subscription information")
	//		subBalanceAfterRequest := wrapperSubscription.NativeBalance
	//		require.Equal(t, expectedSubBalanceWei, subBalanceAfterRequest)
	//
	//		consumerStatus, err := wrapperContracts.LoadTestConsumers[0].GetRequestStatus(testcontext.Get(t), randomWordsFulfilledEvent.RequestId)
	//		require.NoError(t, err, "error getting rand request status")
	//		require.True(t, consumerStatus.Fulfilled)
	//
	//		expectedWrapperConsumerWeiBalance := new(big.Int).Sub(wrapperConsumerBalanceBeforeRequestWei, consumerStatus.Paid)
	//
	//		wrapperConsumerBalanceAfterRequestWei, err := sethClient.Client.BalanceAt(testcontext.Get(t), common.HexToAddress(wrapperContracts.LoadTestConsumers[0].Address()), nil)
	//		require.NoError(t, err, "error getting wrapper consumer balance")
	//		require.Equal(t, expectedWrapperConsumerWeiBalance, wrapperConsumerBalanceAfterRequestWei)
	//
	//		//todo: uncomment when VRF-651 will be fixed
	//		//require.Equal(t, 1, consumerStatus.Paid.Cmp(randomWordsFulfilledEvent.Payment), "Expected Consumer contract pay more than the Coordinator Sub")
	//		vrfcommon.LogFulfillmentDetailsNativeBilling(l, wrapperConsumerBalanceBeforeRequestWei, wrapperConsumerBalanceAfterRequestWei, consumerStatus, randomWordsFulfilledEvent)
	//
	//		require.Equal(t, *testConfig.NumberOfWords, uint32(len(consumerStatus.RandomWords)))
	//		for _, w := range consumerStatus.RandomWords {
	//			l.Info().Str("Output", w.String()).Msg("Randomness fulfilled")
	//			require.Equal(t, 1, w.Cmp(big.NewInt(0)), "Expected the VRF job give an answer bigger than 0")
	//		}
	//	})
	//})
	//t.Run("Canceling Sub And Returning Funds", func(t *testing.T) {
	//	configCopy := config.MustCopy().(tc.TestConfig)
	//	_, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
	//		testcontext.Get(t),
	//		env,
	//		chainID,
	//		vrfContracts.CoordinatorV2Plus,
	//		configCopy,
	//		vrfContracts.LinkToken,
	//		1,
	//		1,
	//		l,
	//	)
	//	require.NoError(t, err, "error setting up new consumers and subs")
	//	subID := subIDs[0]
	//	subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
	//	require.NoError(t, err, "error getting subscription information")
	//	vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)
	//	subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)
	//
	//	testWalletAddress, err := actions.GenerateWallet()
	//	require.NoError(t, err)
	//
	//	testWalletBalanceNativeBeforeSubCancelling, err := sethClient.Client.BalanceAt(testcontext.Get(t), testWalletAddress, nil)
	//	require.NoError(t, err)
	//
	//	testWalletBalanceLinkBeforeSubCancelling, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), testWalletAddress.String())
	//	require.NoError(t, err)
	//
	//	subscriptionForCancelling, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
	//	require.NoError(t, err, "error getting subscription information")
	//
	//	subBalanceLink := subscriptionForCancelling.Balance
	//	subBalanceNative := subscriptionForCancelling.NativeBalance
	//	l.Info().
	//		Str("Subscription Amount Native", subBalanceNative.String()).
	//		Str("Subscription Amount Link", subBalanceLink.String()).
	//		Str("Returning funds from SubID", subID.String()).
	//		Str("Returning funds to", testWalletAddress.String()).
	//		Msg("Canceling subscription and returning funds to subscription owner")
	//
	//	cancellationTx, cancellationEvent, err := vrfContracts.CoordinatorV2Plus.CancelSubscription(subID, testWalletAddress)
	//	require.NoError(t, err, "Error canceling subscription")
	//
	//	txGasUsed := new(big.Int).SetUint64(cancellationTx.Receipt.GasUsed)
	//	// we don't have that information for older Geth versions
	//	if cancellationTx.Receipt.EffectiveGasPrice == nil {
	//		cancellationTx.Receipt.EffectiveGasPrice = new(big.Int).SetUint64(0)
	//	}
	//	cancellationTxFeeWei := new(big.Int).Mul(txGasUsed, cancellationTx.Receipt.EffectiveGasPrice)
	//
	//	l.Info().
	//		Str("Cancellation Tx Fee Wei", cancellationTxFeeWei.String()).
	//		Str("Effective Gas Price", cancellationTx.Receipt.EffectiveGasPrice.String()).
	//		Uint64("Gas Used", cancellationTx.Receipt.GasUsed).
	//		Msg("Cancellation TX Receipt")
	//
	//	l.Info().
	//		Str("Returned Subscription Amount Native", cancellationEvent.AmountLink.String()).
	//		Str("Returned Subscription Amount Link", cancellationEvent.AmountLink.String()).
	//		Str("SubID", cancellationEvent.SubId.String()).
	//		Str("Returned to", cancellationEvent.To.String()).
	//		Msg("Subscription Canceled Event")
	//
	//	require.Equal(t, subBalanceNative, cancellationEvent.AmountNative, "SubscriptionCanceled event native amount is not equal to sub amount while canceling subscription")
	//	require.Equal(t, subBalanceLink, cancellationEvent.AmountLink, "SubscriptionCanceled event LINK amount is not equal to sub amount while canceling subscription")
	//
	//	testWalletBalanceNativeAfterSubCancelling, err := sethClient.Client.BalanceAt(testcontext.Get(t), testWalletAddress, nil)
	//	require.NoError(t, err)
	//
	//	testWalletBalanceLinkAfterSubCancelling, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), testWalletAddress.String())
	//	require.NoError(t, err)
	//
	//	//Verify that sub was deleted from Coordinator
	//	_, err = vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
	//	require.Error(t, err, "error not occurred when trying to get deleted subscription from old Coordinator after sub migration")
	//
	//	subFundsReturnedNativeActual := new(big.Int).Sub(testWalletBalanceNativeAfterSubCancelling, testWalletBalanceNativeBeforeSubCancelling)
	//	subFundsReturnedLinkActual := new(big.Int).Sub(testWalletBalanceLinkAfterSubCancelling, testWalletBalanceLinkBeforeSubCancelling)
	//
	//	subFundsReturnedNativeExpected := new(big.Int).Sub(subBalanceNative, cancellationTxFeeWei)
	//	deltaSpentOnCancellationTxFee := new(big.Int).Sub(subBalanceNative, subFundsReturnedNativeActual)
	//	l.Info().
	//		Str("Sub Balance - Native", subBalanceNative.String()).
	//		Str("Delta Spent On Cancellation Tx Fee - `NativeBalance - subFundsReturnedNativeActual`", deltaSpentOnCancellationTxFee.String()).
	//		Str("Cancellation Tx Fee Wei", cancellationTxFeeWei.String()).
	//		Str("Sub Funds Returned Actual - Native", subFundsReturnedNativeActual.String()).
	//		Str("Sub Funds Returned Expected - `NativeBalance - cancellationTxFeeWei`", subFundsReturnedNativeExpected.String()).
	//		Str("Sub Funds Returned Actual - Link", subFundsReturnedLinkActual.String()).
	//		Str("Sub Balance - Link", subBalanceLink.String()).
	//		Msg("Sub funds returned")
	//
	//	//todo - this fails on SIMULATED env as tx cost is calculated different as for testnets and it's not receipt.EffectiveGasPrice*receipt.GasUsed
	//	//require.Equal(t, subFundsReturnedNativeExpected, subFundsReturnedNativeActual, "Returned funds are not equal to sub balance that was cancelled")
	//	require.Equal(t, 1, testWalletBalanceNativeAfterSubCancelling.Cmp(testWalletBalanceNativeBeforeSubCancelling), "Native funds were not returned after sub cancellation")
	//	require.Equal(t, 0, subBalanceLink.Cmp(subFundsReturnedLinkActual), "Returned LINK funds are not equal to sub balance that was cancelled")
	//
	//})
	//t.Run("Owner Canceling Sub And Returning Funds While Having Pending Requests", func(t *testing.T) {
	//	configCopy := config.MustCopy().(tc.TestConfig)
	//	testConfig := configCopy.VRFv2Plus.General
	//
	//	//underfund subs in order rand fulfillments to fail
	//	testConfig.SubscriptionFundingAmountNative = ptr.Ptr(float64(0))
	//	testConfig.SubscriptionFundingAmountLink = ptr.Ptr(float64(0))
	//
	//	consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
	//		testcontext.Get(t),
	//		env,
	//		chainID,
	//		vrfContracts.CoordinatorV2Plus,
	//		configCopy,
	//		vrfContracts.LinkToken,
	//		1,
	//		1,
	//		l,
	//	)
	//	require.NoError(t, err, "error setting up new consumers and subs")
	//	subID := subIDs[0]
	//	subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
	//	require.NoError(t, err, "error getting subscription information")
	//	vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)
	//	subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)
	//	activeSubscriptionIdsBeforeSubCancellation, err := vrfContracts.CoordinatorV2Plus.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
	//	require.NoError(t, err)
	//
	//	require.True(t, it_utils.BigIntSliceContains(activeSubscriptionIdsBeforeSubCancellation, subID))
	//
	//	pendingRequestsExist, err := vrfContracts.CoordinatorV2Plus.PendingRequestsExist(testcontext.Get(t), subID)
	//	require.NoError(t, err)
	//	require.False(t, pendingRequestsExist, "Pending requests should not exist")
	//
	//	configCopy.VRFv2Plus.General.RandomWordsFulfilledEventTimeout = ptr.Ptr(blockchain.StrDuration{Duration: 5 * time.Second})
	//	_, _, err = vrfv2plus.RequestRandomnessAndWaitForFulfillment(
	//		consumers[0],
	//		vrfContracts.CoordinatorV2Plus,
	//		vrfKey,
	//		subID,
	//		false,
	//		configCopy.VRFv2Plus.General,
	//		l,
	//		0,
	//	)
	//
	//	require.Error(t, err, "error should occur for waiting for fulfilment due to low sub balance")
	//
	//	_, _, err = vrfv2plus.RequestRandomnessAndWaitForFulfillment(
	//		consumers[0],
	//		vrfContracts.CoordinatorV2Plus,
	//		vrfKey,
	//		subID,
	//		true,
	//		configCopy.VRFv2Plus.General,
	//		l,
	//		0,
	//	)
	//
	//	require.Error(t, err, "error should occur for waiting for fulfilment due to low sub balance")
	//
	//	pendingRequestsExist, err = vrfContracts.CoordinatorV2Plus.PendingRequestsExist(testcontext.Get(t), subID)
	//	require.NoError(t, err)
	//	require.True(t, pendingRequestsExist, "Pending requests should exist after unfulfilled rand requests due to low sub balance")
	//
	//	walletBalanceNativeBeforeSubCancelling, err := sethClient.Client.BalanceAt(testcontext.Get(t), common.HexToAddress(sethClient.MustGetRootKeyAddress().Hex()), nil)
	//	require.NoError(t, err)
	//
	//	walletBalanceLinkBeforeSubCancelling, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), sethClient.MustGetRootKeyAddress().Hex())
	//	require.NoError(t, err)
	//
	//	subscriptionForCancelling, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
	//	require.NoError(t, err, "error getting subscription information")
	//
	//	subBalanceLink := subscriptionForCancelling.Balance
	//	subBalanceNative := subscriptionForCancelling.NativeBalance
	//	l.Info().
	//		Str("Subscription Amount Native", subBalanceNative.String()).
	//		Str("Subscription Amount Link", subBalanceLink.String()).
	//		Str("Returning funds from SubID", subID.String()).
	//		Str("Returning funds to", sethClient.MustGetRootKeyAddress().Hex()).
	//		Msg("Canceling subscription and returning funds to subscription owner")
	//
	//	cancellationTx, cancellationEvent, err := vrfContracts.CoordinatorV2Plus.OwnerCancelSubscription(subID)
	//	require.NoError(t, err, "Error canceling subscription")
	//
	//	txGasUsed := new(big.Int).SetUint64(cancellationTx.Receipt.GasUsed)
	//	// we don't have that information for older Geth versions
	//	if cancellationTx.Receipt.EffectiveGasPrice == nil {
	//		cancellationTx.Receipt.EffectiveGasPrice = new(big.Int).SetUint64(0)
	//	}
	//	cancellationTxFeeWei := new(big.Int).Mul(txGasUsed, cancellationTx.Receipt.EffectiveGasPrice)
	//
	//	l.Info().
	//		Str("Cancellation Tx Fee Wei", cancellationTxFeeWei.String()).
	//		Str("Effective Gas Price", cancellationTx.Receipt.EffectiveGasPrice.String()).
	//		Uint64("Gas Used", cancellationTx.Receipt.GasUsed).
	//		Msg("Cancellation TX Receipt")
	//
	//	l.Info().
	//		Str("Returned Subscription Amount Native", cancellationEvent.AmountNative.String()).
	//		Str("Returned Subscription Amount Link", cancellationEvent.AmountLink.String()).
	//		Str("SubID", cancellationEvent.SubId.String()).
	//		Str("Returned to", cancellationEvent.To.String()).
	//		Msg("Subscription Canceled Event")
	//
	//	require.Equal(t, subBalanceNative, cancellationEvent.AmountNative, "SubscriptionCanceled event native amount is not equal to sub amount while canceling subscription")
	//	require.Equal(t, subBalanceLink, cancellationEvent.AmountLink, "SubscriptionCanceled event LINK amount is not equal to sub amount while canceling subscription")
	//
	//	walletBalanceNativeAfterSubCancelling, err := sethClient.Client.BalanceAt(testcontext.Get(t), common.HexToAddress(sethClient.MustGetRootKeyAddress().Hex()), nil)
	//	require.NoError(t, err)
	//
	//	walletBalanceLinkAfterSubCancelling, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), sethClient.MustGetRootKeyAddress().Hex())
	//	require.NoError(t, err)
	//
	//	//Verify that sub was deleted from Coordinator
	//	_, err = vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
	//	require.Error(t, err, "error not occurred when trying to get deleted subscription from old Coordinator after sub migration")
	//
	//	subFundsReturnedNativeActual := new(big.Int).Sub(walletBalanceNativeAfterSubCancelling, walletBalanceNativeBeforeSubCancelling)
	//	subFundsReturnedLinkActual := new(big.Int).Sub(walletBalanceLinkAfterSubCancelling, walletBalanceLinkBeforeSubCancelling)
	//
	//	subFundsReturnedNativeExpected := new(big.Int).Sub(subBalanceNative, cancellationTxFeeWei)
	//	deltaSpentOnCancellationTxFee := new(big.Int).Sub(subBalanceNative, subFundsReturnedNativeActual)
	//	l.Info().
	//		Str("Sub Balance - Native", subBalanceNative.String()).
	//		Str("Delta Spent On Cancellation Tx Fee - `NativeBalance - subFundsReturnedNativeActual`", deltaSpentOnCancellationTxFee.String()).
	//		Str("Cancellation Tx Fee Wei", cancellationTxFeeWei.String()).
	//		Str("Sub Funds Returned Actual - Native", subFundsReturnedNativeActual.String()).
	//		Str("Sub Funds Returned Expected - `NativeBalance - cancellationTxFeeWei`", subFundsReturnedNativeExpected.String()).
	//		Str("Sub Funds Returned Actual - Link", subFundsReturnedLinkActual.String()).
	//		Str("Sub Balance - Link", subBalanceLink.String()).
	//		Str("walletBalanceNativeBeforeSubCancelling", walletBalanceNativeBeforeSubCancelling.String()).
	//		Str("walletBalanceNativeAfterSubCancelling", walletBalanceNativeAfterSubCancelling.String()).
	//		Msg("Sub funds returned")
	//
	//	//todo - need to use different wallet for each test to verify exact amount of Native/LINK returned
	//	//todo - as defaultWallet is used in other tests in parallel which might affect the balance - TT-684
	//	//require.Equal(t, 1, walletBalanceNativeAfterSubCancelling.Cmp(walletBalanceNativeBeforeSubCancelling), "Native funds were not returned after sub cancellation")
	//
	//	//todo - this fails on SIMULATED env as tx cost is calculated different as for testnets and it's not receipt.EffectiveGasPrice*receipt.GasUsed
	//	//require.Equal(t, subFundsReturnedNativeExpected, subFundsReturnedNativeActual, "Returned funds are not equal to sub balance that was cancelled")
	//	require.Equal(t, 0, subBalanceLink.Cmp(subFundsReturnedLinkActual), "Returned LINK funds are not equal to sub balance that was cancelled")
	//
	//	activeSubscriptionIdsAfterSubCancellation, err := vrfContracts.CoordinatorV2Plus.GetActiveSubscriptionIds(testcontext.Get(t), big.NewInt(0), big.NewInt(0))
	//	require.NoError(t, err, "error getting active subscription ids")
	//
	//	require.False(
	//		t,
	//		it_utils.BigIntSliceContains(activeSubscriptionIdsAfterSubCancellation, subID),
	//		"Active subscription ids should not contain sub id after sub cancellation",
	//	)
	//})
	//t.Run("Owner Withdraw", func(t *testing.T) {
	//	configCopy := config.MustCopy().(tc.TestConfig)
	//	consumers, subIDs, err := vrfv2plus.SetupNewConsumersAndSubs(
	//		testcontext.Get(t),
	//		env,
	//		chainID,
	//		vrfContracts.CoordinatorV2Plus,
	//		configCopy,
	//		vrfContracts.LinkToken,
	//		1,
	//		1,
	//		l,
	//	)
	//	require.NoError(t, err, "error setting up new consumers and subs")
	//	subID := subIDs[0]
	//	subscription, err := vrfContracts.CoordinatorV2Plus.GetSubscription(testcontext.Get(t), subID)
	//	require.NoError(t, err, "error getting subscription information")
	//	vrfcommon.LogSubDetails(l, subscription, subID.String(), vrfContracts.CoordinatorV2Plus)
	//	subIDsForCancellingAfterTest = append(subIDsForCancellingAfterTest, subIDs...)
	//
	//	_, fulfilledEventLink, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
	//		consumers[0],
	//		vrfContracts.CoordinatorV2Plus,
	//		vrfKey,
	//		subID,
	//		false,
	//		configCopy.VRFv2Plus.General,
	//		l,
	//		0,
	//	)
	//	require.NoError(t, err)
	//
	//	_, fulfilledEventNative, err := vrfv2plus.RequestRandomnessAndWaitForFulfillment(
	//		consumers[0],
	//		vrfContracts.CoordinatorV2Plus,
	//		vrfKey,
	//		subID,
	//		true,
	//		configCopy.VRFv2Plus.General,
	//		l,
	//		0,
	//	)
	//	require.NoError(t, err)
	//	amountToWithdrawLink := fulfilledEventLink.Payment
	//
	//	defaultWalletBalanceNativeBeforeWithdraw, err := sethClient.Client.BalanceAt(testcontext.Get(t), common.HexToAddress(sethClient.MustGetRootKeyAddress().Hex()), nil)
	//	require.NoError(t, err)
	//
	//	defaultWalletBalanceLinkBeforeWithdraw, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), sethClient.MustGetRootKeyAddress().Hex())
	//	require.NoError(t, err)
	//
	//	l.Info().
	//		Str("Returning to", sethClient.MustGetRootKeyAddress().Hex()).
	//		Str("Amount", amountToWithdrawLink.String()).
	//		Msg("Invoking Oracle Withdraw for LINK")
	//
	//	err = vrfContracts.CoordinatorV2Plus.Withdraw(
	//		common.HexToAddress(sethClient.MustGetRootKeyAddress().Hex()),
	//	)
	//	require.NoError(t, err, "error withdrawing LINK from coordinator to default wallet")
	//	amountToWithdrawNative := fulfilledEventNative.Payment
	//
	//	l.Info().
	//		Str("Returning to", sethClient.MustGetRootKeyAddress().Hex()).
	//		Str("Amount", amountToWithdrawNative.String()).
	//		Msg("Invoking Oracle Withdraw for Native")
	//
	//	err = vrfContracts.CoordinatorV2Plus.WithdrawNative(
	//		common.HexToAddress(sethClient.MustGetRootKeyAddress().Hex()),
	//	)
	//	require.NoError(t, err, "error withdrawing Native tokens from coordinator to default wallet")
	//
	//	defaultWalletBalanceNativeAfterWithdraw, err := sethClient.Client.BalanceAt(testcontext.Get(t), common.HexToAddress(sethClient.MustGetRootKeyAddress().Hex()), nil)
	//	require.NoError(t, err)
	//
	//	defaultWalletBalanceLinkAfterWithdraw, err := vrfContracts.LinkToken.BalanceOf(testcontext.Get(t), sethClient.MustGetRootKeyAddress().Hex())
	//	require.NoError(t, err)
	//
	//	//not possible to verify exact amount of Native/LINK returned as defaultWallet is used in other tests in parallel which might affect the balance
	//	require.Equal(t, 1, defaultWalletBalanceNativeAfterWithdraw.Cmp(defaultWalletBalanceNativeBeforeWithdraw), "Native funds were not returned after oracle withdraw native")
	//	require.Equal(t, 1, defaultWalletBalanceLinkAfterWithdraw.Cmp(defaultWalletBalanceLinkBeforeWithdraw), "LINK funds were not returned after oracle withdraw")
	//})
}
