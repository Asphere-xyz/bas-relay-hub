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

    event DebugUint256(string, bytes32);

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

    function extractParliaSigningData(bytes calldata blockProof, uint256 chainId) external pure returns (bytes memory signingData) {
        // support of >=64 kB headers might make code much more complicated
        require(blockProof.length < 65536);
        uint256 it = RLP.openRlp(blockProof);
        uint256 originalLength = RLP.itemLength(it);
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
        // remember old and new extra data prefix lengths
        uint256 oldExtraDataPrefixLength = RLP.prefixLength(beforeExtraDataOffset);
        uint256 newExtraDataPrefixLength = RLP.estimatePrefixLength(afterExtraDataOffset - beforeExtraDataOffset - oldExtraDataPrefixLength - 65);
        // form signing data from block proof
        bytes memory chainRlp = RLP.uintToRlp(chainId);
        bytes memory signingData = new bytes(originalLength + newExtraDataPrefixLength - oldExtraDataPrefixLength + chainRlp.length - 65);
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
            // extraDataLength = afterExtraDataOffset - beforeExtraDataOffset - oldExtraDataPrefixLength + newExtraDataPrefixLength - 65
//            let extraDataLength := add(sub(sub(sub(afterExtraDataOffset, beforeExtraDataOffset), oldExtraDataPrefixLength), 65), newExtraDataPrefixLength)
            // copy first bytes before extra data
            let startAt := add(signingData, add(mload(chainRlp), 0x23))
            calldatacopy(startAt, add(blockProof.offset, 3), sub(beforeExtraDataOffset, add(blockProof.offset, 3)))
            // copy extra data but leave some bytes for the prefix
            startAt := add(add(startAt, sub(beforeExtraDataOffset, add(blockProof.offset, 3))))
            calldatacopy(startAt, startAt)
        }
        // copy extra data
//        assembly {
//            let startAt := add(signingData, add(mload(chainRlp), 0x23))
//            startAt := add(startAt)
//        }
        assembly {
            //let startOffset := add(signingData, add(mload(chainRlp), 0x23))
            //let lengthWithExtraData := sub(sub(afterExtraDataOffset, 65), add(blockProof.offset, newExtraDataPrefixLength))
            //calldatacopy(startOffset, add(blockProof.offset, 3), lengthWithExtraData)
            //calldatacopy(add(startOffset, lengthWithExtraData), afterExtraDataOffset, add(afterExtraDataOffset, 42))
        }
        // patch extra data length inside RLP signing data
        {
            uint256 extraDataLength = RLP.itemLength(beforeExtraDataOffset) - RLP.prefixLength(beforeExtraDataOffset) - 65;
//            uint256 patchExtraDataAt;
//            assembly {
//                patchExtraDataAt := add(sub(beforeExtraDataOffset, blockProof.offset), 2)
//            }
//            signingData[patchExtraDataAt + 0] = bytes1(uint8(extraDataLength >> 8));
//            signingData[patchExtraDataAt + 1] = bytes1(uint8(extraDataLength >> 0));
        }
        return signingData;
    }

    function verifyProof(bytes calldata proof, bytes32[] calldata existingValidatorSet) external pure returns (bytes32[] memory newValidatorSet) {
        return existingValidatorSet;
    }
}