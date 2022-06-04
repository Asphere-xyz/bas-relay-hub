// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.0;

import "./IValidatorSet.sol";

interface IProofVerificationFunction is IValidatorSet {

    function verifyCheckpointBlock(bytes calldata genesisBlock, uint256 chainId, bytes32 checkpointHash) external view returns (address[] memory initialValidatorSet);

    function verifyGenesisBlock(bytes calldata genesisBlock, uint256 chainId) external view returns (address[] memory initialValidatorSet);

    function verifyValidatorTransition(bytes[] calldata blockProofs, uint256 chainId, address[] calldata existingValidatorSet) external view returns (address[] memory newValidatorSet, uint64 epochNumber);
}