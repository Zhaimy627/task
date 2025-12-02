// src/MetaStake.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/**
 * @title MetaStake Pro - 多池质押挖矿系统（2025 生产级）
 */
contract MetaStake is Initializable, UUPSUpgradeable, Ownable2StepUpgradeable {
    using SafeERC20 for IERC20;

    IERC20 public rewardToken;
    uint256 public rewardPerBlock = 10 * 1e18;
    uint256 public totalWeight;
    bool public stakePaused;
    bool public withdrawPaused;

    struct PoolInfo {
        IERC20 stToken;
        uint128 poolWeight;
        uint128 lastRewardBlock;
        uint256 accPerShare;
        uint256 totalStaked;
        uint256 minDeposit;
        uint64  lockBlocks;
    }

    struct UserInfo {
        uint256 amount;
        uint256 rewardDebt;
        uint256 pendingRewards;
    }

    struct UnstakeRequest {
        uint256 amount;
        uint256 unlockBlock;
    }

    PoolInfo[] public pools;
    mapping(uint256 => mapping(address => UserInfo)) public userInfo;
    // 关键：必须加 public 才能在测试里访问 .length
    mapping(uint256 => mapping(address => UnstakeRequest[])) public unstakeRequests;

    event Deposit(address indexed user, uint256 indexed pid, uint256 amount);
    event Withdraw(address indexed user, uint256 indexed pid, uint256 amount);
    event Harvest(address indexed user, uint256 indexed pid, uint256 amount);
    event PoolAdded(uint256 indexed pid, address stToken, uint256 weight);
    event UnstakeRequested(address indexed user, uint256 indexed pid, uint256 amount, uint256 unlockBlock);

    constructor() {
        _disableInitializers();
    }

    function initialize(address _rewardToken) public initializer {
        __Ownable2Step_init();
        // v5 自动包含 UUPS 初始化，无需手动调用
        rewardToken = IERC20(_rewardToken);
    }

    function addPool(
        address _stToken,
        uint256 _weight,
        uint256 _minDeposit,
        uint256 _lockBlocks
    ) external onlyOwner {
        require(_weight > 0, "Weight > 0");
        _updateAllPools();

        pools.push(PoolInfo({
            stToken: IERC20(_stToken),
            poolWeight: uint128(_weight),
            lastRewardBlock: uint128(block.number),
            accPerShare: 0,
            totalStaked: 0,
            minDeposit: _minDeposit,
            lockBlocks: uint64(_lockBlocks)
        }));
        totalWeight += _weight;
        emit PoolAdded(pools.length - 1, _stToken, _weight);
    }

    function setPause(bool _stake, bool _withdraw) external onlyOwner {
        stakePaused = _stake;
        withdrawPaused = _withdraw;
    }

    function deposit(uint256 _pid) external payable {
        require(!stakePaused, "Stake paused");
        PoolInfo storage pool = pools[_pid];
        require(address(pool.stToken) == address(0), "Only ETH pool");
        require(msg.value >= pool.minDeposit, "Below min");

        _updatePool(_pid);
        UserInfo storage user = userInfo[_pid][msg.sender];

        if (user.amount > 0) {
            uint256 pending = (user.amount * pool.accPerShare / 1e18) - user.rewardDebt;
            user.pendingRewards += pending;
        }

        user.amount += msg.value;
        pool.totalStaked += msg.value;
        user.rewardDebt = user.amount * pool.accPerShare / 1e18;

        emit Deposit(msg.sender, _pid, msg.value);
    }

    function withdraw(uint256 _pid, uint256 _amount) external {
        require(!withdrawPaused, "Withdraw paused");
        UserInfo storage user = userInfo[_pid][msg.sender];
        require(user.amount >= _amount, "Insufficient");

        _updatePool(_pid);
        uint256 pending = (user.amount * pools[_pid].accPerShare / 1e18) - user.rewardDebt;
        user.pendingRewards += pending;

        user.amount -= _amount;
        user.rewardDebt = user.amount * pools[_pid].accPerShare / 1e18;

        unstakeRequests[_pid][msg.sender].push(UnstakeRequest({
            amount: _amount,
            unlockBlock: block.number + pools[_pid].lockBlocks
        }));

        emit UnstakeRequested(msg.sender, _pid, _amount, block.number + pools[_pid].lockBlocks);
        emit Withdraw(msg.sender, _pid, _amount);
    }

    function harvest(uint256 _pid) external {
        _updatePool(_pid);
        UserInfo storage user = userInfo[_pid][msg.sender];
        uint256 pending = user.pendingRewards;
        if (pending > 0) {
            user.pendingRewards = 0;
            rewardToken.safeTransfer(msg.sender, pending);
            emit Harvest(msg.sender, _pid, pending);
        }
    }

    function claimUnstaked(uint256 _pid, uint256 _index) external {
        UnstakeRequest[] storage requests = unstakeRequests[_pid][msg.sender];
        require(_index < requests.length, "Invalid index");
        require(block.number >= requests[_index].unlockBlock, "Still locked");

        uint256 amount = requests[_index].amount;
        requests[_index] = requests[requests.length - 1];
        requests.pop();

        (bool sent, ) = payable(msg.sender).call{value: amount}("");
        require(sent, "ETH transfer failed");
    }

    function _updatePool(uint256 _pid) internal {
        PoolInfo storage pool = pools[_pid];
        if (block.number <= pool.lastRewardBlock || pool.totalStaked == 0) {
            pool.lastRewardBlock = uint128(block.number);
            return;
        }
        uint256 blocks = block.number - pool.lastRewardBlock;
        uint256 reward = blocks * rewardPerBlock * pool.poolWeight / totalWeight;
        pool.accPerShare += reward * 1e18 / pool.totalStaked;
        pool.lastRewardBlock = uint128(block.number);
    }

    function _updateAllPools() internal {
        for (uint256 i = 0; i < pools.length; i++) {
            _updatePool(i);
        }
    }

    function _authorizeUpgrade(address) internal override onlyOwner {}
}