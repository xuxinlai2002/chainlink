const {
    deployContract,
    writeConfig,
} = require('./utils/helper')

const main = async () => {
 
   
    let accounts = await ethers.getSigners();
    let deployer = accounts[0];

    console.log("xxl 2 deployer : ",deployer.address);

    contract = await deployContract(deployer, "src/v0.8/vrf/dev/BlockhashStore.sol:BlockhashStore");
    await writeConfig("0config","0config","BlockhashStore",contract.address);

}

main();