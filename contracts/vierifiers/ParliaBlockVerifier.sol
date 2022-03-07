// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.0;

import "../interfaces/IProofVerificationFunction.sol";

import "../libraries/CallDataRLPReader.sol";
import "../libraries/RLPReader.sol";
import "../libraries/ParliaParser.sol";

contract ParliaBlockVerifier is IProofVerificationFunction {

    uint256 internal constant EXTRA_VANITY = 32;
    uint256 internal constant ADDRESS_LENGTH = 20;
    uint256 internal constant EXTRA_SEAL = 65;
    uint256 internal constant NEXT_FORK_HASHES = 4;

    event ChainId(uint256);

//    struct ParliaBlockHeader {
//        uint256 chainId;
//        bytes32 parentHash;
//        bytes32 uncleHash;
//        bytes32 stateRoot;
//        bytes32 transactionsRoot;
//        bytes32 receiptsRoot;
//        bytes32 bloomFilter;
//        uint256 number;
//        uint256 gasLimit;
//        uint64 gasUsed;
//        uint64 time;
//        bytes extraData;
//        bytes32 mixDigest;
//        bytes32 nonce;
//    }

    function extractSigningData(bytes calldata blockProof) external {
        uint256 it = CallDataRLPReader.openRlp(blockProof);

//        RLPReader.RLPItem memory init = RLPReader.toRlpItem(blockProof);
//        uint256 fromOffset = init.memPtr;
//        RLPReader.Iterator memory it = RLPReader.iterator(init);
//        // chain id
//        {
//            uint256 chainId = RLPReader.toUint(RLPReader.next(it));
//        }
//        // parent hash
//        {
//            uint256 parentHash = RLPReader.toUint(RLPReader.next(it));
//        }
//        // uncle hash
//        RLPReader.next(it);
//        // coinbase
//        {
//            address coinbase = RLPReader.toAddress(RLPReader.next(it));
//        }
//        // state root, transactions root
//        RLPReader.next(it);
//        RLPReader.next(it);
//        // receipts root
//        {
//            uint256 receiptsRoot = RLPReader.toUint(RLPReader.next(it));
//        }
//        // bloom filter, difficulty
//        RLPReader.next(it);
//        RLPReader.next(it);
//        // block number
//        {
//            uint256 blockNumber = RLPReader.toUint(RLPReader.next(it));
//        }
//        // gas limit, gas used, time
//        RLPReader.next(it);
//        RLPReader.next(it);
//        uint256 bytesBetween = RLPReader.next(it).memPtr - fromOffset;
//
//        // extra data (vanity[32] + validators[20*N] + forkHash[4] + signature)
//
//        {
//            bytes memory extraData = RLPReader.toBytes(RLPReader.next(it));
//            // payload is blockProofOffset->extraDatOffset + (extraDataOffset + extraDataLength - EXTRA_SEAL) + rest
////            extraData[:extraData.length - 65];
//        }


        {
//            uint256 extraDataOffset = it;
//            RLPReader.next(it);
//            uint256 extraDataLength = RLPReader.payloadLen(extraDataOffset, it - extraDataOffset);
            //            require(extraDataLength >= 65);
            //            uint256 signatureOffset = extraDataOffset + extraDataLength - 65;
            //            bytes32 signature;
            //            assembly {
            //                signature := calldataload(signatureOffset)
            //            }
            //            // payload is blockProofOffset->extraDatOffset + (extraDataOffset + extraDataLength - EXTRA_SEAL) + rest
            //            bytes memory signingData;
            //            assembly {
            //                calldatacopy(add(signingData, 0x20), blockProofOffset, sub(extraDataOffset, blockProofOffset))
            //                calldatacopy(add(signingData, 0x20), sub(extraDataOffset, blockProofOffset), sub(extraDataLength, 65))
            //                mstore(signingData, add(sub(extraDataOffset, blockProofOffset), sub(extraDataLength, 65)))
            //            }
            //            return signingData;
        }

        // mix digest, nonce
//        RLPReader.next(it);
//        RLPReader.next(it);
    }

    function verifyProof(uint256 blockProofOffset, bytes32[] calldata existingValidatorSet) external pure returns (bytes32[] memory newValidatorSet) {
//        uint256 it = CallDataRLPReader.beginIteration(blockProofOffset);
//        // chain id
//        uint256 chainId = CallDataRLPReader.toUintStrict(it);
//        it = CallDataRLPReader.next(it);
//        // parent hash
//        bytes32 parentHash = CallDataRLPReader.toBytes32(it);
//        it = CallDataRLPReader.next(it);
//        // uncle hash
//        it = CallDataRLPReader.next(it);
//        // coinbase
//        address coinbase = CallDataRLPReader.toAddress(it);
//        it = CallDataRLPReader.next(it);
//        // state root, transactions root
//        it = CallDataRLPReader.next(it);
//        it = CallDataRLPReader.next(it);
//        // receipts root
//        bytes32 receiptsRoot = CallDataRLPReader.toBytes32(it);
//        it = CallDataRLPReader.next(it);
//        // bloom filter, difficulty
//        it = CallDataRLPReader.next(it);
//        it = CallDataRLPReader.next(it);
//        // block number
//        uint256 blockNumber = CallDataRLPReader.toUintStrict(it);
//        it = CallDataRLPReader.next(it);
//        // gas limit, gas used, time
//        it = CallDataRLPReader.next(it);
//        it = CallDataRLPReader.next(it);
//        it = CallDataRLPReader.next(it);
//        // extra data (vanity[32] + validators[20*N] + forkHash[4] + signature)
//        {
//            uint256 extraDataOffset = it;
//            it = CallDataRLPReader.next(it);
//            uint256 extraDataLength = CallDataRLPReader.payloadLen(extraDataOffset, it - extraDataOffset);
//            require(extraDataLength >= EXTRA_SEAL);
//            uint256 signatureOffset = extraDataOffset + extraDataLength - EXTRA_SEAL;
//            bytes32 signature;
//            assembly {
//                signature := calldataload(signatureOffset)
//            }
//            // payload is blockProofOffset->extraDatOffset + (extraDataOffset + extraDataLength - EXTRA_SEAL) + rest
//            bytes memory signingData;
//            assembly {
//                calldatacopy(add(signingData, 0x20), blockProofOffset, sub(extraDataOffset, blockProofOffset))
//                calldatacopy(add(signingData, 0x20), sub(extraDataOffset, blockProofOffset), sub(extraDataLength, 65))
//                mstore(signingData, add(sub(extraDataOffset, blockProofOffset), sub(extraDataLength, 65)))
//            }
//        }
//
//        // mix digest, nonce
//        it = CallDataRLPReader.next(it);
//        it = CallDataRLPReader.next(it);

        return existingValidatorSet;
    }
}