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
    let VRFV2PlusLoadTestWithMetrics = await readConfig("0config","VRFV2PlusLoadTestWithMetrics");

    let VRFCoordinatorV2_5Contract = await attachContract(
        "VRFCoordinatorV2_5",
        VRFCoordinatorV2_5Address,
        deployer
    );

    let pubKeys = ["0xe0a57be0970f68f1d612a050bed5f8799d509de60e7367768fa1dc2b51e7ac5c","0x540352695411384e8f187c651539a793ad7b93244f7de903d9867c18406f370f"];
    let result = await isContractTransferSuccess(

        await VRFCoordinatorV2_5Contract.registerProvingKey(
            pubKeys,
            "10000000000",
        )
    )

    console.log("xxl addConsumer result ",result);
   


}

main();

// 82355e22-a200-4ccb-8ed4-6e2ace8264fb