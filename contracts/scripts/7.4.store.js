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

    let blockhashStoreddress = await readConfig("0config","BlockhashStore");
    
    let blockhashStoreContract = await attachContract(
        "src/v0.8/vrf/dev/BlockhashStore.sol:BlockhashStore",
        blockhashStoreddress,
        deployer
    );

    let blockNumber = 71970175
    let result = await isContractTransferSuccess(

        await blockhashStoreContract.store(
            blockNumber,{
                gasPrice, gasLimit,
            }
        )
    )

    console.log("xxl store result ",result);
   


}

main();

// 82355e22-a200-4ccb-8ed4-6e2ace8264fb