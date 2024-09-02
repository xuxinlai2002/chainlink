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

    let subId = await readConfig("0config","SubId");

    let txObj = await VRFV2PlusLoadTestWithMetricsContract.requestRandomWords(
        subId,
        3,
        "0x60b510b4e6c29abdf5ff00f492c2196749b3fdc29a662b9679faff16feae52d1",
        "700000",
        false,
        3,
        1,{
            gasPrice, gasLimit,
        }
    )    
    let txRep = await txObj.wait();

//     console.log("xxl requestRandomWords result ",txRep.blockNumber);

    let blockhashStoreddress = await readConfig("0config","BlockhashStore");
    let blockhashStoreContract = await attachContract(
        "src/v0.8/vrf/dev/BlockhashStore.sol:BlockhashStore",
        blockhashStoreddress,
        deployer
    );

    let result = await isContractTransferSuccess(

        await blockhashStoreContract.store(
            txRep.blockNumber,{
                gasPrice, gasLimit,
            }
        )
    )

    console.log("xxl store result ",result);

}

main();

// 82355e22-a200-4ccb-8ed4-6e2ace8264fb