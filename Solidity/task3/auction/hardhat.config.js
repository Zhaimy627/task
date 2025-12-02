require("@nomicfoundation/hardhat-toolbox");
require('hardhat-deploy')
require("@nomicfoundation/hardhat-ethers");
require('@openzeppelin/hardhat-upgrades')
/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.28",
  namedAccounts: {
    deployer: 0,
    user: 1,
    user2: 2,
  },
};
