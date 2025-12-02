// SPDX-License-Identifier: MIT
pragma solidity >=0.6.2 <0.9.0;

import "./Vm.sol";

contract Test {
    Vm public constant vm = Vm(address(uint160(uint256(keccak256("hevm cheat code")))));

    string private failedMsg;
    bool private testFailed;

    modifier mayRevert() { _; }
    modifier showsMessage(string memory msg) { _; }

    function fail(string memory err) internal virtual {
        failedMsg = err;
        testFailed = true;
    }
}
