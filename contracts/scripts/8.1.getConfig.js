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


    let result = await VRFCoordinatorV2_5Contract.s_config()
    console.log("xxl txRep result ",result);

}

main();

// 82355e22-a200-4ccb-8ed4-6e2ace8264fb