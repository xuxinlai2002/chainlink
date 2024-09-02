const {
    readConfig,
    writeConfig,
    deployContract,
} = require('./utils/helper')

const main = async () => {
 
   
    let accounts = await ethers.getSigners();
    let deployer = accounts[0];

    let blockhashStoreAddress = await readConfig("0config","BlockhashStore");

    console.log("xxl 4 deployer : ",deployer.address);

    contract = await deployContract(deployer, "VRFCoordinatorV2_5",blockhashStoreAddress);
    await writeConfig("0config","0config","VRFCoordinatorV2_5",contract.address);

}

main();