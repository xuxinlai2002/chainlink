const {
    readConfig,
    writeConfig,
    deployContract,
} = require('./utils/helper')

const main = async () => {
 
   
    let accounts = await ethers.getSigners();
    let deployer = accounts[0];

    let VRFCoordinatorV2PlusAddress = await readConfig("0config","VRFCoordinatorV2_5");

    console.log("xxl 5 deployer : ",deployer.address);

    contract = await deployContract(deployer, "BatchVRFCoordinatorV2Plus",VRFCoordinatorV2PlusAddress);
    await writeConfig("0config","0config","BatchVRFCoordinatorV2Plus",contract.address);

}

main();