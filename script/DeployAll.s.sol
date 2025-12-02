// script/DeployAll.s.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "forge-std/Script.sol";
import "forge-std/console.sol";
import "../src/MetaNode.sol";
import "../src/MetaStake.sol";
import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract DeployAll is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);
        console.log("Deployer address:", deployer);

        vm.startBroadcast(deployerPrivateKey);

        MetaStake implementation = new MetaStake();
        console.log("MetaStake implementation deployed to:", address(implementation));

        ERC1967Proxy proxy = new ERC1967Proxy(address(implementation), "");
        console.log("ERC1967Proxy deployed to:", address(proxy));

        MetaNode metaNode = new MetaNode(address(proxy));
        console.log("MetaNode (Reward Token) deployed to:", address(metaNode));

        MetaStake stake = MetaStake(address(proxy));
        stake.initialize(address(metaNode));
        console.log("MetaStake initialized");

        metaNode.transferOwnership(address(stake));
        console.log("MetaNode ownership transferred to MetaStake");

        stake.transferOwnership(deployer);
        console.log("MetaStake ownership transferred to deployer:", deployer);

        stake.addPool(address(0), 100, 0.01 ether, 10080);
        console.log("ETH Pool (pid: 0) added successfully!");

        console.log("=== DEPLOYMENT SUCCESSFUL ===");
        console.log("MetaNode:", address(metaNode));
        console.log("MetaStake (Proxy):", address(stake));
        console.log("View on Sepolia:");
        console.log("https://sepolia.etherscan.io/address/", address(stake));

        vm.stopBroadcast();
    }
}