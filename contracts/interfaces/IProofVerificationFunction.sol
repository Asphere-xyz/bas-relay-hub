// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.0;

interface IProofVerificationFunction {

    function verifyCheckpointBlock(bytes calldata genesisBlock, uint256 chainId, bytes32 checkpointHash) external view returns (address[] memory initialValidatorSet);

    function verifyGenesisBlock(bytes calldata genesisBlock, uint256 chainId) external view returns (address[] memory initialValidatorSet);

    function verifyValidatorTransition(bytes[] calldata blockProofs, uint256 chainId, address[] calldata existingValidatorSet) external view returns (address[] memory newValidatorSet, uint64 epochNumber);

    struct VerifiedBlock {
        bytes32 blockHash;
        bytes32 parentHash;
        uint64 blockNumber;
        bytes32 stateRoot;
        bytes32 txRoot;
        bytes32 receiptRoot;
    }

    function verifyBlock(bytes[] calldata blockProofs, uint256 chainId, address[] calldata existingValidatorSet) external view returns (VerifiedBlock memory result);
}