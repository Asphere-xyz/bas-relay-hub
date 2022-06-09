// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.0;

import "./ParliaBlockVerifier.sol";

contract EmbeddedParliaBlockVerifier is ParliaBlockVerifier {

    address constant internal VERIFY_PARLIA_BLOCK_PRECOMPILE = address(0x0000000000000000000000000000004241530001);

    constructor(uint32 confirmationBlocks, uint32 epochInterval) ParliaBlockVerifier(confirmationBlocks, epochInterval) {
    }

    function _extractParliaSigningData(bytes calldata blockProof, uint256 chainId) internal view override returns (VerifiedParliaBlockResult memory) {
        bytes memory input = abi.encode(chainId, blockProof, _epochInterval);
        bytes memory output = new bytes(blockProof.length);
        assembly {
            let status := staticcall(0, 0x0000000000000000000000000000004241530001, add(input, 0x20), mload(input), add(output, 0x20), mload(output))
            switch status
            case 0 {
                revert(add(output, 0x20), returndatasize())
            }
        }
        return abi.decode(output, (VerifiedParliaBlockResult));
    }
}