// SPDX-License-Identifier: MIT
pragma solidity >=0.6.2 <0.9.0;

interface Vm {
    function envUint(string calldata name) external returns (uint256);
    function envInt(string calldata name) external returns (int256);
    function envAddress(string calldata name) external returns (address);
    function envBytes32(string calldata name) external returns (bytes32);
    function envBool(string calldata name) external returns (bool);
    function envString(string calldata name) external returns (string memory);
    function envBytes(string calldata name) external returns (bytes memory);

    function startBroadcast() external;
    function startBroadcast(uint256 privateKey) external;
    function startBroadcast(address from) external;
    function stopBroadcast() external;

    function addr(uint256 privateKey) external returns (address);
    function warp(uint256 newTimestamp) external;
    function roll(uint256 newBlockNumber) external;
    function fee(uint256 newBlockFee) external;
    function store(address target, bytes32 slot, bytes32 data) external;
    function load(address target, bytes32 slot) external returns (bytes32);
}
