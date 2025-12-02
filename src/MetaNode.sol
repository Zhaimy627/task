// src/MetaNode.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title MetaNode 奖励代币（v5 完全兼容版）
 * @dev 部署时直接把所有权给质押合约，100% 合规
 */
contract MetaNode is ERC20, Ownable {
    constructor(address stakeContract) 
        ERC20("MetaNode", "MTN") 
        Ownable(stakeContract)   // v5 必须传初始 owner
    {
        require(stakeContract != address(0), "Stake contract cannot be zero");
        // 10亿总供应，全给质押合约
        _mint(stakeContract, 1_000_000_000 * 1e18);
    }
}