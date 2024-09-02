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


    let txObj = await VRFCoordinatorV2_5Contract.createSubscription()
    let txRep = await txObj.wait();

    console.log("xxl txRep result ",txRep.events[0].args.subId);

}

main();

// 82355e22-a200-4ccb-8ed4-6e2ace8264fb