const {deployments,upgrades, ethers} =require("hardhat")
const fs = require("fs");
const path = require("path");
// deploy/00_deploy_my_contract.js
// export a function that get passed the Hardhat runtime environment
module.exports = async ({ getNamedAccounts, deployments }) => {
  const { save } = deployments;
  const { deployer } = await getNamedAccounts();
  console.log("部署用户地址",deployer);
  const NftAuction = await ethers.getContractFactory("NftAuction");
  //通过代理合约部署
  const ntAuctionProxy =await upgrades.deployProxy(NftAuction, 
    [],{initializer: "initialize"}
);
   await ntAuctionProxy.waitForDeployment();
   console.log("代理合约地址", await ntAuctionProxy.getAddress());
   const implAddress =await upgrades.erc1967.getImplementationAddress(proxyAddress);
   console.log("实现合约地址",implAddress)

   const storePath = path.join(__dirname, "./.catche/proxyNftAuction.josn");
//    console.log("storePath",JSON.stringify
//     ({proxyAddress,
//         implAddress,
//         abi: ntAuctionProxy.interface.format("json"),
//        }),);
   fs.writeFileSync(storePath, JSON.stringify({proxyAddress,
    implAddress,
    implAddress,
    abi: NftAuction.interface.format("json"),
   }));
//    await save("NftAuctionProxy", {
//     proxyAddress,
//     implAddress,
//     implAddress: await upgrades.erc1967.getImplementationAddress(proxyAddress),
//   });
   await save("NftAuctionProxy", {
    abi: NftAuction.interface.format("json"),
    address: implAddress,
    // args: [],
    // log: true,
  }); 
//   await deploy("MyContract", {
//     from: deployer,
//     args: ["Hello"],
//     log: true,
//   });
};
// add tags and dependencies
module.exports.tags = ["deployNFTAuction"];
//deploy/01_deploy_nft_auction.js
// const fs = require("fs");
// const path = require("path");
// module.exports = async ({ getNamedAccounts, deployments, getChainId }) => {
//   const { deploy, save } = deployments;        // ← 这行必须有
//   const { deployer } = await getNamedAccounts(); // ← 部署者地址

//   console.log("部署用户地址", deployer);

//   // 如果是普通合约（不可升级）用这行
//   // const nftAuction = await deploy("NftAuction", {
//   //   from: deployer,
//   //   args: [], // 构造函数参数
//   //   log: true,
//   // });

//   // 如果是可升级合约（推荐）用这行
//   const nftAuction = await deploy("NftAuction", {
//     from: deployer,
//     proxy: {
//       proxyContract: "OpenZeppelinTransparentProxy",
//       execute: {
//         init: {
//           methodName: "initialize",   // 如果有 initialize 函数就写这个
//           args: [],                   // initialize 参数
//         },
//       },
//     },
//     log: true,
//   });

//   console.log("代理合约地址", nftAuction.address);

//   // 可选：手动保存一下 implementation 地址方便后面升级
//   if (nftAuction.newlyDeployed && nftAuction.implementation) {
//     await save("NftAuction_Implementation", {
//       abi: nftAuction.abi,
//       address: nftAuction.implementation,
//     });
//   }
// };
// const storePath = path.join(__dirname, "./.catche/proxyNftAuction.josn");



// module.exports.tags = ["NftAuction"];