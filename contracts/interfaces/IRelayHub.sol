// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.0;

interface IRelayHub {

    function checkReceiptProof(uint256 chainId, bytes[] calldata blockProofs, bytes memory rawReceipt, bytes memory path, bytes calldata siblings) external view returns (bool);
}