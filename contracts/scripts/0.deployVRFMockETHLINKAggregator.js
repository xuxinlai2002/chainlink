const {
    deployContract,
    writeConfig,
} = require('./utils/helper')

const main = async () => {
 
   
    let accounts = await ethers.getSigners();
    let deployer = accounts[0];

    console.log("xxl 0 deployer : ",deployer.address);

    contract = await deployContract(deployer, "VRFMockETHLINKAggregator","1000000000000000000");
    await writeConfig("0config","0config","VRFMockETHLINKAggregator",contract.address);

}

main();

