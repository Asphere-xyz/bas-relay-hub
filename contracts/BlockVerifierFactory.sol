// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.0;

import "./verifiers/EmbeddedParliaBlockVerifier.sol";
import "./verifiers/ParliaBlockVerifier.sol";

contract BlockVerifierFactory {

    function factoryParliaBlockVerifier() external returns (IProofVerificationFunction) {
        return new ParliaBlockVerifier();
    }

    function factoryEmbeddedParliaBlockVerifier() external returns (IProofVerificationFunction) {
        return new EmbeddedParliaBlockVerifier();
    }
}