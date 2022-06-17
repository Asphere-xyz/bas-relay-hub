// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

import "../interfaces/IProofVerificationFunction.sol";

import "../libraries/RLP.sol";

contract ParliaBlockVerifier is IProofVerificationFunction {

    using EnumerableSet for EnumerableSet.AddressSet;
    using BitMaps for BitMaps.BitMap;

    uint32 internal immutable _confirmationBlocks;
    uint32 internal immutable _epochInterval;

    constructor(uint32 confirmationBlocks, uint32 epochInterval) {
        _confirmationBlocks = confirmationBlocks;
        _epochInterval = epochInterval;
    }

    function extractParliaSigningData(bytes calldata blockProof, uint256 chainId) external view returns (VerifiedParliaBlockResult memory result) {
        return _extractParliaSigningData(blockProof, chainId);
    }

    struct VerifiedParliaBlockResult {
        bytes32 blockHash;
        uint64 blockNumber;
        address coinbase;
        address[] validators;
        bytes32 parentHash;
        bytes signingData;
        bytes signature;
    }

    function _extractParliaSigningData(bytes calldata blockProof, uint256 chainId) internal virtual view returns (VerifiedParliaBlockResult memory result) {
        // support of >64 kB headers might make code much more complicated and such blocks doesn't exist
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
        result.parentHash = RLP.toBytes32(it);
        it += 445;
        // slow skip for variadic fields: difficulty, number, gas limit, gas used, time
        it = RLP.next(it);
        result.blockNumber = uint64(RLP.toUint256(it, RLP.itemLength(it)));
        it = RLP.next(RLP.next(RLP.next(RLP.next(it))));
        // calculate and remember offsets for extra data begin and end
        uint256 beforeExtraDataOffset = it;
        it = RLP.next(it);
        uint256 afterExtraDataOffset = it;
        // create chain id and extra data RLPs
        uint256 oldExtraDataPrefixLength = RLP.prefixLength(beforeExtraDataOffset);
        uint256 newExtraDataPrefixLength;
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
            let dst := add(signingData, add(mload(chainRlp), 0x23)) // 0x20+3 (3 is a size of prefix for 64kB list)
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
        bytes memory signature = new bytes(65);
        assembly {
            calldatacopy(add(signature, 0x20), sub(afterExtraDataOffset, 65), 65)
        }
        // recover signer from signature (genesis block doesn't have signature)
        if (result.blockNumber != 0) {
            if (signature[64] == bytes1(uint8(1))) {
                signature[64] = bytes1(uint8(28));
            } else {
                signature[64] = bytes1(uint8(27));
            }
            result.signingData = signingData;
            result.signature = signature;
            result.coinbase = ECDSA.recover(keccak256(signingData), signature);
        }
        // parse validators for zero block epoch
        if (result.blockNumber % _epochInterval == 0) {
            uint256 totalValidators = (afterExtraDataOffset - beforeExtraDataOffset + oldExtraDataPrefixLength - 65 - 32) / 20;
            address[] memory newValidators = new address[](totalValidators);
            for (uint256 i = 0; i < totalValidators; i++) {
                uint256 validator;
                assembly {
                    validator := calldataload(add(add(add(beforeExtraDataOffset, oldExtraDataPrefixLength), mul(i, 20)), 32))
                }
                newValidators[i] = address(uint160(validator >> 96));
            }
            result.validators = newValidators;
        }
        // calc block hash
        result.blockHash = keccak256(blockProof);
        return result;
    }

    function verifyCheckpointBlock(bytes calldata genesisBlock, uint256 chainId, bytes32 checkpointHash) external view override returns (address[] memory initialValidatorSet) {
        VerifiedParliaBlockResult memory result = _extractParliaSigningData(genesisBlock, chainId);
        require(result.blockHash == checkpointHash, "not a checkpoint block");
        return result.validators;
    }

    function verifyGenesisBlock(bytes calldata genesisBlock, uint256 chainId) external view override returns (address[] memory initialValidatorSet) {
        VerifiedParliaBlockResult memory result = _extractParliaSigningData(genesisBlock, chainId);
        require(result.blockNumber == 0, "not a genesis block");
        return result.validators;
    }

    function verifyValidatorTransition(bytes[] calldata blockProofs, uint256 chainId, address[] calldata existingValidatorSet) external view returns (address[] memory newValidatorSet, uint64 epochNumber) {
        bytes32 parentHash;
        // copy to the stack to avoid SLOAD's
        (uint32 confirmationBlocks, uint32 epochInterval) = (_confirmationBlocks, _epochInterval);
        require(blockProofs.length >= confirmationBlocks, "not enough proofs");
        // we must store somehow set of active validators to check is quorum reached
        address[] memory uniqueValidators = new address[](blockProofs.length);
        uint64 uniqueValidatorsLength = 0;
        // check all blocks
        for (uint256 i = 0; i < confirmationBlocks; i++) {
            VerifiedParliaBlockResult memory result = _extractParliaSigningData(blockProofs[i], chainId);
            address signer = result.coinbase;
            // make sure signer exists (we should know validator order, it can be optimized)
            bool signerFound = false;
            for (uint256 j = 0; j < existingValidatorSet.length; j++) {
                if (existingValidatorSet[j] != signer) continue;
                signerFound = true;
                break;
            }
            require(signerFound, "unknown signer");
            bool uniqueFound = false;
            for (uint256 j = 0; j < uniqueValidatorsLength; j++) {
                if (uniqueValidators[j] != signer) continue;
                uniqueFound = true;
                break;
            }
            if (!uniqueFound) {
                uniqueValidators[uniqueValidatorsLength] = signer;
                uniqueValidatorsLength++;
            }
            // first block must be epoch block
            if (i == 0) {
                require(result.blockNumber % epochInterval == 0, "epoch block");
                epochNumber = result.blockNumber / epochInterval;
                newValidatorSet = result.validators;
                parentHash = result.blockHash;
            } else {
                require(result.parentHash == parentHash, "bad parent hash");
                parentHash = result.blockHash;
            }
        }
        require(uniqueValidatorsLength >= confirmationBlocks, "quorum not reached");
        return (newValidatorSet, epochNumber);
    }

    function verifyBlock(bytes[] calldata blockProofs, uint256 chainId, address[] calldata existingValidatorSet) external view returns (VerifiedBlock memory verifiedBlock) {
        require(blockProofs.length >= _confirmationBlocks, "make sure proofs are enough");
        // we must store somehow set of active validators to check is quorum reached
        address[] memory uniqueValidators = new address[](blockProofs.length);
        uint64 uniqueValidatorsLength = 0;
        // check all blocks
        bytes32 parentHash;
        for (uint256 i = 0; i < _confirmationBlocks; i++) {
            VerifiedParliaBlockResult memory result = _extractParliaSigningData(blockProofs[i], chainId);
            address signer = result.coinbase;
            // make sure signer exists (we should know validator order, it can be optimized)
            bool signerFound = false;
            for (uint256 j = 0; j < existingValidatorSet.length; j++) {
                if (existingValidatorSet[j] != signer) continue;
                signerFound = true;
                break;
            }
            require(signerFound, "unknown signer");
            bool uniqueFound = false;
            for (uint256 j = 0; j < uniqueValidatorsLength; j++) {
                if (uniqueValidators[j] != signer) continue;
                uniqueFound = true;
                break;
            }
            if (!uniqueFound) {
                uniqueValidators[uniqueValidatorsLength] = signer;
                uniqueValidatorsLength++;
            }
            // first block is block with proof
            if (i == 0) {
                ParliaBlockHeader memory pbh = parseParliaBlockHeader(blockProofs[i]);
                verifiedBlock.blockHash = result.blockHash;
                verifiedBlock.parentHash = result.parentHash;
                verifiedBlock.blockNumber = result.blockNumber;
                verifiedBlock.stateRoot = pbh.stateRoot;
                verifiedBlock.txRoot = pbh.txRoot;
                verifiedBlock.receiptRoot = pbh.receiptRoot;
                parentHash = result.blockHash;
            } else {
                require(result.parentHash == parentHash, "bad parent hash");
                parentHash = result.blockHash;
            }
        }
        require(uniqueValidatorsLength >= _confirmationBlocks, "quorum not reached");
        return verifiedBlock;
    }

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

    function parseParliaBlockHeader(bytes calldata blockProof) public pure returns (ParliaBlockHeader memory pbh) {
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
}