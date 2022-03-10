// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.0;

interface IProofVerificationFunction {

    function verifyProof(bytes calldata proof, bytes32[] calldata existingValidatorSet) external pure returns (bytes32[] memory newValidatorSet);
}