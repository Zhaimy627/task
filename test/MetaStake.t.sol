// test/MetaStake.t.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "forge-std/Test.sol";
import "../src/MetaStake.sol";
import "../src/MetaNode.sol";
import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract MockERC20 is ERC20 {
    constructor() ERC20("Mock Token", "MOCK") {
        _mint(msg.sender, 1000000 * 1e18);
    }
    function mint(address to, uint256 amount) public { _mint(to, amount); }
}

contract MetaStakeTest is Test {
    MetaStake stake;
    MetaNode rewardToken;
    MockERC20 mockToken;
    ERC1967Proxy proxy;
    MetaStake implementation;

    address deployer = address(0x1111);
    address user1 = address(0x2222);

    // 复制合约里的事件
    event UnstakeRequested(address indexed user, uint256 indexed pid, uint256 amount, uint256 unlockBlock);

   function setUp() public {
    vm.startPrank(deployer);

    // 1. 部署实现合约
    implementation = new MetaStake();

    // 2. 部署代理
    proxy = new ERC1967Proxy(address(implementation), "");

    // 3. 部署奖励代币（owner 直接是代理）
    rewardToken = new MetaNode(address(proxy));

    // 4. 初始化代理合约
    stake = MetaStake(address(proxy));
    stake.initialize(address(rewardToken));

    // 关键：到这里，代理合约才真正有 owner（就是 stake 自己）
    // 现在可以安全调用 addPool 了！

    stake.addPool(address(0), 100, 0.01 ether, 100);     // ETH 池
    mockToken = new MockERC20();
    stake.addPool(address(mockToken), 50, 10 * 1e18, 50); // ERC20 池

    vm.stopPrank();
}

    function test_ETH_Deposit() public {
        vm.deal(user1, 1 ether);
        vm.prank(user1);
        stake.deposit{value: 0.5 ether}(0);

        (uint256 amount,,) = stake.userInfo(0, user1);
        assertEq(amount, 0.5 ether);
    }

    function test_Withdraw_With_Lock() public {
        vm.deal(user1, 1 ether);
        vm.prank(user1);
        stake.deposit{value: 0.5 ether}(0);

        vm.expectEmit(true, true, false, true);
        emit UnstakeRequested(user1, 0, 0.2 ether, block.number + 100);

        vm.prank(user1);
        stake.withdraw(0, 0.2 ether);
    }

    function test_Claim_After_Unlock() public {
        test_Withdraw_With_Lock();

        vm.roll(block.number + 101);

        uint256 balBefore = user1.balance;
        vm.prank(user1);
        stake.claimUnstaked(0, 0);
        assertGt(user1.balance, balBefore);
    }

    function test_Harvest_Rewards() public {
        vm.deal(user1, 1 ether);
        vm.prank(user1);
        stake.deposit{value: 0.5 ether}(0);

        vm.roll(block.number + 1000);
        vm.prank(user1);
        stake.harvest(0);

        assertGt(rewardToken.balanceOf(user1), 0);
    }

    function test_Pause() public {
        vm.prank(deployer);
        stake.setPause(true, false);

        vm.deal(user1, 1 ether);
        vm.expectRevert("Stake paused");
        vm.prank(user1);
        stake.deposit{value: 0.1 ether}(0);
    }
}