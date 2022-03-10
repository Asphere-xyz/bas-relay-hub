// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.0;

import "../interfaces/IProofVerificationFunction.sol";

import "../libraries/RLP.sol";
import "../libraries/ParliaParser.sol";

contract ParliaBlockVerifier is IProofVerificationFunction {

    uint256 internal constant EXTRA_VANITY = 32;
    uint256 internal constant ADDRESS_LENGTH = 20;
    uint256 internal constant EXTRA_SEAL = 65;
    uint256 internal constant NEXT_FORK_HASHES = 4;

    event DebugUint256(string, uint256);

    struct ParliaBlockHeader {
        bytes32 parentHash;
        bytes32 uncleHash;
        address coinbase;
        bytes32 stateRoot;
        bytes32 txRoot;
        bytes32 receiptRoot;
        uint64 blockNumber;
        uint64 gasLimit;
        uint64 gasUsed;
        uint64 blockTime;
        bytes32 mixDigest;
        uint64 nonce;
        bytes32 blockHash;
    }

    function extractParliaSigningData(bytes calldata blockProof, uint256 chainId) external pure returns (bytes memory signingData, bytes memory signature) {
        // support of >64 kB headers might make code much more complicated
        require(blockProof.length <= 65535);
        // open RLP and calc block header length after the prefix (it should be block proof length -3)
        uint256 it = RLP.openRlp(blockProof);
        uint256 originalLength = RLP.itemLength(it);
        // skip body length
        it = RLP.beginIteration(it);
        // fast skip for the fixed fields:
        // + 33 bytes (x5): parent hash, uncle hash, state root, tx hash, receipt hash
        // + 259 bytes (x1): bloom filter
        // + 21 bytes (x1): coinbase
        // total skip: 33*5+259+21=445
        it += 445;
        // slow skip for variadic fields: difficulty, number, gas limit, gas used, time
        it = RLP.next(RLP.next(RLP.next(RLP.next(RLP.next(it)))));
        // calculate and remember offsets for extra data begin and end
        uint256 beforeExtraDataOffset = it;
        it = RLP.next(it);
        uint256 afterExtraDataOffset = it;
        // create chain id and extra data RLPs
        uint256 oldExtraDataPrefixLength = RLP.prefixLength(beforeExtraDataOffset);
        uint256 newExtraDataPrefixLength = 3;
        {
            uint256 newEstExtraDataLength = afterExtraDataOffset - beforeExtraDataOffset - oldExtraDataPrefixLength - 65;
            if (newEstExtraDataLength < 56) {
                newExtraDataPrefixLength = 1;
            } else {
                newExtraDataPrefixLength = 1 + RLP.uintRlpPrefixLength(newEstExtraDataLength);
            }
        }
        bytes memory chainRlp = RLP.uintToRlp(chainId);
        // form signing data from block proof
        bytes memory signingData = new bytes(chainRlp.length + originalLength - oldExtraDataPrefixLength + newExtraDataPrefixLength - 65);
        // init first 3 bytes of signing data with RLP prefix and encoded length
        {
            signingData[0] = 0xf9;
            uint256 bodyLength = signingData.length - 3;
            signingData[1] = bytes1(uint8(bodyLength >> 8));
            signingData[2] = bytes1(uint8(bodyLength >> 0));
        }
        // copy chain id rlp right after the prefix
        for (uint256 i = 0; i < chainRlp.length; i++) {
            signingData[3 + i] = chainRlp[i];
        }
        // copy block calldata to the signing data before extra data [0;extraData-65)
        assembly {
        // copy first bytes before extra data
            let dst := add(signingData, add(mload(chainRlp), 0x23)) // 0x20+3 (size of prefix for 64kB list)
            let src := add(blockProof.offset, 3)
            let len := sub(beforeExtraDataOffset, src)
            calldatacopy(dst, src, len)
        // copy extra data with new prefix
            dst := add(add(dst, len), newExtraDataPrefixLength)
            src := add(beforeExtraDataOffset, oldExtraDataPrefixLength)
            len := sub(sub(sub(afterExtraDataOffset, beforeExtraDataOffset), oldExtraDataPrefixLength), 65)
            calldatacopy(dst, src, len)
        // copy rest (mix digest, nonce)
            dst := add(dst, len)
            src := afterExtraDataOffset
            len := 42 // its always 42 bytes
            calldatacopy(dst, src, len)
        }
        // patch extra data length inside RLP signing data
        {
            uint256 newExtraDataLength;
            uint256 patchExtraDataAt;
            assembly {
                newExtraDataLength := sub(sub(sub(afterExtraDataOffset, beforeExtraDataOffset), oldExtraDataPrefixLength), 65)
                patchExtraDataAt := sub(mload(signingData), add(add(newExtraDataLength, newExtraDataPrefixLength), 42))
            }
            // we don't need to cover more than 3 cases because we revert if block header >64kB
            if (newExtraDataPrefixLength == 4) {
                signingData[patchExtraDataAt + 0] = bytes1(uint8(0xb7 + 3));
                signingData[patchExtraDataAt + 1] = bytes1(uint8(newExtraDataLength >> 16));
                signingData[patchExtraDataAt + 2] = bytes1(uint8(newExtraDataLength >> 8));
                signingData[patchExtraDataAt + 3] = bytes1(uint8(newExtraDataLength >> 0));
            } else if (newExtraDataPrefixLength == 3) {
                signingData[patchExtraDataAt + 0] = bytes1(uint8(0xb7 + 2));
                signingData[patchExtraDataAt + 1] = bytes1(uint8(newExtraDataLength >> 8));
                signingData[patchExtraDataAt + 2] = bytes1(uint8(newExtraDataLength >> 0));
            } else if (newExtraDataPrefixLength == 2) {
                signingData[patchExtraDataAt + 0] = bytes1(uint8(0xb7 + 1));
                signingData[patchExtraDataAt + 1] = bytes1(uint8(newExtraDataLength >> 8));
            } else if (newExtraDataLength < 56) {
                signingData[patchExtraDataAt + 0] = bytes1(uint8(0x80 + newExtraDataLength));
            }
            // else can't be here, its unreachable
        }
        // save signature
        signature = new bytes(65);
        assembly {
            calldatacopy(add(signature, 0x20), sub(afterExtraDataOffset, 65), 65)
        }
        return (signingData, signature);
    }

    function parseParliaBlockHeader(bytes calldata blockProof) external pure returns (ParliaBlockHeader memory pbh) {
        uint256 it = RLP.beginRlp(blockProof);
        // parent hash, uncle hash
        pbh.parentHash = RLP.toBytes32(it);
        it = RLP.next(it);
        pbh.uncleHash = RLP.toBytes32(it);
        it = RLP.next(it);
        // coinbase
        pbh.coinbase = RLP.toAddress(it);
        it = RLP.next(it);
        // state root, transactions root, receipts root
        pbh.stateRoot = RLP.toBytes32(it);
        it = RLP.next(it);
        pbh.txRoot = RLP.toBytes32(it);
        it = RLP.next(it);
        pbh.receiptRoot = RLP.toBytes32(it);
        it = RLP.next(it);
        // bloom, difficulty
        it = RLP.next(it);
        it = RLP.next(it);
        // block number, gas limit, gas used, time
        pbh.blockNumber = uint64(RLP.toUint256(it, RLP.itemLength(it)));
        it = RLP.next(it);
        pbh.gasLimit = uint64(RLP.toUint256(it, RLP.itemLength(it)));
        it = RLP.next(it);
        pbh.gasUsed = uint64(RLP.toUint256(it, RLP.itemLength(it)));
        it = RLP.next(it);
        pbh.blockTime = uint64(RLP.toUint256(it, RLP.itemLength(it)));
        it = RLP.next(it);
        // extra data
        it = RLP.next(it);
        // mix digest, nonce
        pbh.mixDigest = RLP.toBytes32(it);
        it = RLP.next(it);
        pbh.nonce = uint64(RLP.toUint256(it, RLP.itemLength(it)));
        it = RLP.next(it);
        // calc block hash
        pbh.blockHash = keccak256(blockProof);
        return pbh;
    }

    function verifyProof(bytes calldata proof, bytes32[] calldata existingValidatorSet) external pure returns (bytes32[] memory newValidatorSet) {
        return existingValidatorSet;
    }
}