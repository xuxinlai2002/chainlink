const {
    deployContract,
    writeConfig,
} = require('./utils/helper')

const main = async () => {
 
   
    let accounts = await ethers.getSigners();
    let deployer = accounts[0];

    console.log("xxl 1 deployer : ",deployer.address);

    contract = await deployContract(deployer, "LinkToken");
    await writeConfig("0config","0config","LinkToken",contract.address);

}

main();