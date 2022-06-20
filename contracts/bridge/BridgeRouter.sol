// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.6;
pragma abicoder v2;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

import "../libraries/ProofParser.sol";
import "../interfaces/ICrossChainBridge.sol";

import "./SimpleTokenProxy.sol";

contract BridgeRouter {

    function peggedTokenAddress(address bridge, address fromToken) public pure returns (address) {
        return SimpleTokenProxyUtils.simpleTokenProxyAddress(bridge, bytes32(bytes20(fromToken)));
    }

    function factoryPeggedToken(address fromToken, address toToken, ICrossChainBridge.Metadata memory metaData, address bridge) public returns (IERC20Mintable) {
        /* we must use delegate call because we need to deploy new contract from bridge contract to have valid address */
        address targetToken = SimpleTokenProxyUtils.deploySimpleTokenProxy(bridge, bytes32(bytes20(fromToken)), metaData);
        require(targetToken == toToken, "bad chain");
        /* to token is our new pegged token */
        return IERC20Mintable(toToken);
    }
}
