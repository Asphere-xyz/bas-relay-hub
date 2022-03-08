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

    function extractSigningData(bytes calldata blockProof) external {
        uint256 it = RLP.openRlp(blockProof);
        uint256 blockProofOffset = it;
        // chain id
        {
            uint256 chainId = RLP.toUint(it, 1);
            emit DebugUint256("chainId", bytes32(chainId));
        }
        it = RLP.next(it);
        // parent hash, uncle hash
        it = RLP.next(it);
        it = RLP.next(it);
        // coinbase
        {
            address coinbase = RLP.toAddress(it);
            emit DebugUint256("coinbase", bytes32(bytes20(coinbase)));
        }
        it = RLP.next(it);
        // state root, transactions root
        {
            uint256 chainId = RLP.toUint(it, 33);
            emit DebugUint256("stateRoot", bytes32(chainId));
        }
        it = RLP.next(it);
        {
            uint256 chainId = RLP.toUint(it, 33);
            emit DebugUint256("txRoot", bytes32(chainId));
        }
        it = RLP.next(it);
        // receipts root
        {
            uint256 receiptsRoot = RLP.toUint(it, 33);
            emit DebugUint256("receiptsRoot", bytes32(receiptsRoot));
        }
        it = RLP.next(it);
        // bloom, difficulty
        it = RLP.next(it);
        it = RLP.next(it);
        // block number
        {
            uint256 blockNumber = RLP.toUint(it, RLP.itemLength(it));
            emit DebugUint256("blockNumber", bytes32(blockNumber));
        }
        it = RLP.next(it);
        // gas limit, gas used, time
        it = RLP.next(it);
        it = RLP.next(it);
        it = RLP.next(it);

        // extra data (vanity[32] + validators[20*N] + forkHash[4] + signature)
        {
            uint256 extraDataOffset = it;
            uint256 extraDataLength = RLP.payloadLen(extraDataOffset, it - extraDataOffset);
            require(extraDataLength >= 65);
            uint256 signatureOffset = extraDataOffset + extraDataLength - 65;
            bytes32 signature;
            assembly {
                signature := calldataload(signatureOffset)
            }
            // payload is blockProofOffset->extraDatOffset + (extraDataOffset + extraDataLength - EXTRA_SEAL) + rest
            bytes memory signingData;
            assembly {
                calldatacopy(add(signingData, 0x20), blockProofOffset, sub(extraDataOffset, blockProofOffset))
                calldatacopy(add(signingData, 0x20), sub(extraDataOffset, blockProofOffset), sub(extraDataLength, 65))
                mstore(signingData, add(sub(extraDataOffset, blockProofOffset), sub(extraDataLength, 65)))
            }
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