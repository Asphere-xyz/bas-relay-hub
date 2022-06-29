// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.0;

import "../verifiers/ParliaBlockVerifier.sol";

contract VerifierGasMeasurer is ParliaBlockVerifier {

    function measureVerifyGas(bytes calldata blockProof, uint256 chainId) external view returns (uint64 gasUsed) {
        gasUsed = uint64(gasleft());
        _extractParliaSigningData(blockProof, chainId, 200);
        return gasUsed - uint64(gasleft());
    }
}