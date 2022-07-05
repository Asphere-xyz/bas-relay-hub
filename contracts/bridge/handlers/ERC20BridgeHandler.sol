// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.6;

interface IBridgeHandler {
}

contract ERC20BridgeHandler is IBridgeHandler {

    function factoryPeggedToken(address bridgeAddress, address originToken) external returns (address) {
    }

    function handleDeposit(address originToken, uint256 originChain, address fromAddress, address toAddress) external {
    }

    function handleWithdraw() external {
    }
}