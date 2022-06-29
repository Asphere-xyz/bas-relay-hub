// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.0;

import "../interfaces/IProofVerificationFunction.sol";

contract FakeBlockVerifier is IProofVerificationFunction {

    function verifyCheckpointBlock(bytes calldata /*genesisBlock*/, uint256 /*chainId*/, bytes32 /*checkpointHash*/, uint32 /*epochLength*/) external pure override returns (address[] memory initialValidatorSet) {
        return initialValidatorSet;
    }

    function verifyGenesisBlock(bytes calldata /*genesisBlock*/, uint256 /*chainId*/, uint32 /*epochLength*/) external pure override returns (address[] memory initialValidatorSet) {
        return initialValidatorSet;
    }

    function verifyValidatorTransition(bytes[] calldata /*blockProofs*/, uint256 /*chainId*/, address[] calldata existingValidatorSet, uint32 /*epochLength*/) external pure override returns (address[] memory newValidatorSet, uint64 epochNumber) {
        return (existingValidatorSet, 0);
    }

    function verifyBlock(bytes[] calldata /*blockProofs*/, uint256 /*chainId*/, address[] calldata /*existingValidatorSet*/, uint32 /*epochLength*/) external pure override returns (VerifiedBlock memory result) {
        return result;
    }
}