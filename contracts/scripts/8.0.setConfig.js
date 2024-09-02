const {
    readConfig,
    writeConfig,
    attachContract,
    isContractTransferSuccess
} = require('./utils/helper')

const gasPrice = "10000000000"
const gasLimit = "20000000"

const main = async () => {
 
   
    let accounts = await ethers.getSigners();
    let deployer = accounts[0];
    console.log("xxl deployer :",deployer.address);

    let VRFCoordinatorV2_5Address = await readConfig("0config","VRFCoordinatorV2_5");

    let VRFCoordinatorV2_5Contract = await attachContract(
        "VRFCoordinatorV2_5",
        VRFCoordinatorV2_5Address,
        deployer
    );

    // uint16 minimumRequestConfirmations,
    // uint32 maxGasLimit,
    // uint32 stalenessSeconds,
    // uint32 gasAfterPaymentCalculation,
    // int256 fallbackWeiPerUnitLink,
    // uint32 fulfillmentFlatFeeNativePPM,
    // uint32 fulfillmentFlatFeeLinkDiscountPPM,
    // uint8 nativePremiumPercentage,
    // uint8 linkPremiumPercentage

    let result = await isContractTransferSuccess(
        await VRFCoordinatorV2_5Contract.setConfig(
            0,
            "16000000",
            0,
            0,
            10,
            0,
            0,
            0,
            0,{
                gasPrice, gasLimit,
            }
        )
    )
    console.log("xxl txRep setConfig ",result);

}

main();

// 82355e22-a200-4ccb-8ed4-6e2ace8264fb