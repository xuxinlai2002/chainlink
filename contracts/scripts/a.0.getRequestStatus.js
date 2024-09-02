const {
    readConfig,
    writeConfig,
    attachContract,
    isContractTransferSuccess
} = require('./utils/helper')

const gasPrice = "1000000000"
const gasLimit = "75000000"

const main = async () => {
 
   
    let accounts = await ethers.getSigners();
    let deployer = accounts[0];
    console.log("xxl deployer :",deployer.address);

    let VRFV2PlusLoadTestWithMetricsAddress = await readConfig("0config","VRFV2PlusLoadTestWithMetrics");

    let VRFV2PlusLoadTestWithMetricsContract = await attachContract(
        "VRFV2PlusLoadTestWithMetrics",
        VRFV2PlusLoadTestWithMetricsAddress,
        deployer
    );


    let result = await VRFV2PlusLoadTestWithMetricsContract.getRequestStatus(
        "73621413529974585259301230257926336279586817889598873452727497212140972465272"
    );
    console.log("xxl result is : ",result);

}

main();

// 82355e22-a200-4ccb-8ed4-6e2ace8264fb