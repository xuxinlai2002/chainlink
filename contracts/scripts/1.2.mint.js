const {
    readConfig,
    attachContract,
    isContractTransferSuccess,
} = require('./utils/helper')


const gasPrice = "2100000000"
const gasLimit = "4000000"

const main = async () => {
 
   
    let accounts = await ethers.getSigners();
    let deployer = accounts[0];

    let linkTokenAddress = await readConfig("0config","LinkToken");

    console.log("xxl 1.1 deployer : ",deployer.address);

    let contract = await attachContract("LinkToken",linkTokenAddress,deployer);

    // function mint(address account, uint256 amount) external override onlyMinter validAddress(account) {
    let result = await isContractTransferSuccess(
        await contract.mint(deployer.address,
            "100000000000000000000",{
                gasPrice, gasLimit,
            }
        )
    )

    console.log("xxl result is : ",result);

}

main();