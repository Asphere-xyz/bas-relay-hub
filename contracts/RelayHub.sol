// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "@openzeppelin/contracts/utils/structs/BitMaps.sol";
import "@openzeppelin/contracts/utils/Multicall.sol";

import "./interfaces/IProofVerificationFunction.sol";
import "./interfaces/IRelayHub.sol";
import "./interfaces/IValidatorChecker.sol";
import "./interfaces/IBridgeRegistry.sol";

import "./libraries/BitUtils.sol";
import "./libraries/MerklePatriciaProof.sol";

import "./BlockVerifierFactory.sol";

contract RelayHub is Multicall, IRelayHub, IBridgeRegistry, IValidatorChecker {

    using EnumerableSet for EnumerableSet.AddressSet;
    using BitMaps for BitMaps.BitMap;

    // lets keep default verification function as zero to make it manageable by BAS relay hub itself
    IProofVerificationFunction internal constant DEFAULT_VERIFICATION_FUNCTION = IProofVerificationFunction(0x0000000000000000000000000000000000000000);

    bytes32 internal constant ZERO_BLOCK_HASH = bytes32(0x00);
    address internal constant ZERO_ADDRESS = address(0x00);

    event ChainRegistered(uint256 indexed chainId, address[] initValidatorSet);
    event ValidatorSetUpdated(uint256 indexed chainId, address[] newValidatorSet);

    struct ValidatorHistory {
        // set with all validators and their indices (never remove values)
        EnumerableSet.AddressSet allValidators;
        // mapping from epoch to the bitmap with active validators indices
        mapping(uint64 => BitMaps.BitMap) activeValidators;
        mapping(uint64 => uint64) validatorCount;
        // latest published epoch
        uint64 latestKnownEpoch;
    }

    enum ChainStatus {
        NotFound,
        Verifying,
        Active
    }

    struct BAS {
        ChainStatus chainStatus;
        IProofVerificationFunction verificationFunction;
        address bridgeAddress;
        uint32 epochLength;
    }

    // default verification function for certified chains
    IProofVerificationFunction internal _defaultVerificationFunction;
    address internal _checkpointOracle;
    // mapping with all registered chains
    mapping(uint256 => ValidatorHistory) _validatorHistories;
    mapping(uint256 => BAS) internal _registeredChains;

    constructor(IProofVerificationFunction defaultVerificationFunction) {
        _defaultVerificationFunction = defaultVerificationFunction;
    }

    function getBridgeAddress(uint256 chainId) external view returns (address) {
        return _registeredChains[chainId].bridgeAddress;
    }

    function registerCertifiedBAS(
        uint256 chainId,
        bytes calldata genesisBlock,
        address bridgeAddress,
        uint32 epochLength
    ) external {
        _registerChainWithVerificationFunction(chainId, DEFAULT_VERIFICATION_FUNCTION, genesisBlock, ZERO_BLOCK_HASH, ChainStatus.Verifying, bridgeAddress, epochLength);
    }

    function registerUsingCheckpoint(
        uint256 chainId,
        bytes calldata checkpointBlock,
        bytes32 checkpointHash,
        address bridgeAddress,
        uint32 epochLength
    ) external {
        _registerChainWithVerificationFunction(chainId, DEFAULT_VERIFICATION_FUNCTION, checkpointBlock, checkpointHash, ChainStatus.Verifying, bridgeAddress, epochLength);
    }

    function registerBAS(
        uint256 chainId,
        IProofVerificationFunction verificationFunction,
        bytes calldata genesisBlock,
        address bridgeAddress,
        uint32 epochLength
    ) external {
        _registerChainWithVerificationFunction(chainId, verificationFunction, genesisBlock, ZERO_BLOCK_HASH, ChainStatus.Verifying, bridgeAddress, epochLength);
    }

    function _registerChainWithVerificationFunction(
        uint256 chainId,
        IProofVerificationFunction verificationFunction,
        bytes calldata blockProof,
        bytes32 checkpointHash,
        ChainStatus defaultStatus,
        address bridgeAddress,
        uint32 epochLength
    ) internal {
        BAS memory bas = _registeredChains[chainId];
        require(bas.chainStatus == ChainStatus.NotFound || bas.chainStatus == ChainStatus.Verifying, "already registered");
        (
        bytes32 blockHash,
        address[] memory initialValidatorSet,
        uint64 blockNumber
        ) = _verificationFunction(verificationFunction).verifyBlockWithoutQuorum(chainId, blockProof, epochLength);
        if (checkpointHash != ZERO_BLOCK_HASH) {
            require(checkpointHash == blockHash, "bad checkpoint hash");
        }
        bas.chainStatus = defaultStatus;
        bas.verificationFunction = verificationFunction;
        bas.bridgeAddress = bridgeAddress;
        bas.epochLength = epochLength;
        ValidatorHistory storage validatorHistory = _validatorHistories[chainId];
        _updateActiveValidatorSet(validatorHistory, initialValidatorSet, blockNumber / epochLength);
        _registeredChains[chainId] = bas;
        emit ChainRegistered(chainId, initialValidatorSet);
    }

    function _updateActiveValidatorSet(ValidatorHistory storage validatorHistory, address[] memory newValidatorSet, uint64 epochNumber) internal {
        // make sure epochs updated one by one (don't do this check for the first transition)
        if (validatorHistory.latestKnownEpoch > 0 && epochNumber > 0) {
            require(epochNumber == validatorHistory.latestKnownEpoch + 1, "bad epoch");
        }
        uint256[] memory buckets = new uint256[]((validatorHistory.allValidators.length() >> 8) + 1);
        // build set of buckets with new bits
        for (uint256 i = 0; i < newValidatorSet.length; i++) {
            // add validator to the set of all validators
            address validator = newValidatorSet[i];
            validatorHistory.allValidators.add(validator);
            // get index of the validator in the set (-1 because 0 is not used)
            uint256 index = validatorHistory.allValidators._inner._indexes[bytes32(uint256(uint160(validator)))] - 1;
            buckets[index >> 8] |= 1 << (index & 0xff);
        }
        // copy buckets (its cheaper to keep buckets in memory)
        BitMaps.BitMap storage currentBitmap = validatorHistory.activeValidators[epochNumber];
        for (uint256 i = 0; i < buckets.length; i++) {
            currentBitmap._data[i] = buckets[i];
        }
        // remember total amount of validators and latest verified epoch
        validatorHistory.validatorCount[epochNumber] = uint64(newValidatorSet.length);
        validatorHistory.latestKnownEpoch = epochNumber;
    }

    function getActiveValidators(uint256 chainId) external view returns (address[] memory) {
        ValidatorHistory storage validatorHistory = _validatorHistories[chainId];
        return _extractActiveValidators(validatorHistory, validatorHistory.latestKnownEpoch);
    }

    function _extractActiveValidators(ValidatorHistory storage validatorHistory, uint64 atEpoch) internal view returns (address[] memory) {
        uint256 validatorsLength = validatorHistory.allValidators.length();
        uint256 totalBuckets = (validatorsLength >> 8) + 1;
        address[] memory activeValidators = new address[](validatorsLength);
        BitMaps.BitMap storage bitmap = validatorHistory.activeValidators[atEpoch];
        uint256 j = 0;
        for (uint256 i = 0; i < totalBuckets; i++) {
            uint256 bucket = bitmap._data[i];
            while (bucket != 0) {
                uint256 zeroes = BitUtils.ctz(bucket);
                bucket ^= (1 << zeroes);
                activeValidators[j] = address(uint160(uint256(bytes32(validatorHistory.allValidators._inner._values[(i << 8) + zeroes]))));
                j++;
            }
        }
        assembly {
            mstore(activeValidators, j)
        }
        return activeValidators;
    }

    function checkValidators(uint256 chainId, address[] memory validators, uint64 epoch) external view returns (uint64 uniqueValidators) {
        ValidatorHistory storage validatorHistory = _validatorHistories[chainId];
        BitMaps.BitMap storage activeValidators = validatorHistory.activeValidators[epoch];
        for (uint256 i = 0; i < validators.length; i++) {
            uint256 index = validatorHistory.allValidators._inner._indexes[bytes32(uint256(uint160(validators[i])))] - 1;
            require(activeValidators.get(index), "not a validator");
            uniqueValidators++;
        }
        return uniqueValidators;
    }

    function getLatestTransitionedEpoch(uint256 chainId) external view returns (uint64) {
        ValidatorHistory storage validatorHistory = _validatorHistories[chainId];
        return validatorHistory.latestKnownEpoch;
    }

    function updateValidatorSet(uint256 chainId, bytes[] calldata blockProofs) external {
        BAS memory bas = _registeredChains[chainId];
        require(bas.chainStatus == ChainStatus.Verifying || bas.chainStatus == ChainStatus.Active, "not active");
        ValidatorHistory storage validatorHistory = _validatorHistories[chainId];
        (address[] memory newValidatorSet, uint64 epochNumber) = _verificationFunction(bas.verificationFunction).verifyValidatorTransition(chainId, blockProofs, bas.epochLength, this);
        bas.chainStatus = ChainStatus.Active;
        _updateActiveValidatorSet(validatorHistory, newValidatorSet, epochNumber);
        _registeredChains[chainId] = bas;
        emit ValidatorSetUpdated(chainId, newValidatorSet);
    }

    function checkpointTransition(
        uint256 chainId,
        bytes calldata rawEpochBlock,
        bytes32 checkpointHash,
        bytes[] calldata signatures
    ) external {
        // make sure bas is registered and active
        BAS memory bas = _registeredChains[chainId];
        require(bas.chainStatus == ChainStatus.Verifying || bas.chainStatus == ChainStatus.Active, "not active");
        // verify next epoch block with new validator set
        (
        bytes32 blockHash,
        address[] memory newValidatorSet,
        uint64 blockNumber
        ) = _verificationFunction(bas.verificationFunction).verifyBlockWithoutQuorum(chainId, rawEpochBlock, bas.epochLength);
        uint64 newEpochNumber = blockNumber / bas.epochLength;
        // lets check signatures and make sure quorum is reached
        {
            address[] memory signers = new address[](signatures.length);
            bytes32 signingRoot = keccak256(abi.encode(blockHash, checkpointHash));
            for (uint256 i = 0; i < signatures.length; i++) {
                signers[i] = ECDSA.recover(signingRoot, signatures[i]);
            }
            require(checkValidatorsAndQuorumReached(chainId, signers, newEpochNumber - 1), "quorum not reached");
        }
        // update validator set
        {
            ValidatorHistory storage validatorHistory = _validatorHistories[chainId];
            _updateActiveValidatorSet(validatorHistory, newValidatorSet, newEpochNumber);
        }
        // remember bas status
        bas.chainStatus = ChainStatus.Active;
        _registeredChains[chainId] = bas;
    }

    function checkValidatorsAndQuorumReached(uint256 chainId, address[] memory validatorSet, uint64 epochNumber) public view returns (bool) {
        // find validator history for epoch and bitmap with active validators
        ValidatorHistory storage validatorHistory = _validatorHistories[chainId];
        BitMaps.BitMap storage bitMap = validatorHistory.activeValidators[epochNumber];
        // we must know total active validators and unique validators to check reachability of the quorum
        uint256 totalValidators = validatorHistory.validatorCount[epochNumber];
        uint256 uniqueValidators = 0;
        uint256[] memory markedValidators = new uint256[]((totalValidators + 0xff) >> 8);
        for (uint256 i = 0; i < validatorSet.length; i++) {
            // find validator's index and make sure it exists in the validator set
            uint256 index = validatorHistory.allValidators._inner._indexes[bytes32(uint256(uint160(validatorSet[i])))] - 1;
            require(bitMap.get(index), "bad validator");
            // mark used validators to be sure quorum is well-calculated
            uint256 usedMask = 1 << (index & 0xff);
            if (markedValidators[index >> 8] & usedMask == 0) {
                uniqueValidators++;
            }
            markedValidators[index >> 8] |= usedMask;
        }
        return uniqueValidators >= totalValidators * 2 / 3;
    }

    function checkReceiptProof(
        uint256 chainId,
        bytes[] calldata blockProofs,
        bytes calldata rawReceipt,
        bytes calldata proofSiblings,
        bytes calldata proofPath
    ) external view virtual override returns (bool) {
        // make sure bas chain is registered and active
        BAS memory bas = _registeredChains[chainId];
        require(bas.chainStatus == ChainStatus.Active, "not active");
        // verify block transition
        IProofVerificationFunction pvf = _verificationFunction(bas.verificationFunction);
        IProofVerificationFunction.BlockHeader memory blockHeader = pvf.verifyBlockAndReachedQuorum(chainId, blockProofs, bas.epochLength, this);
        // check receipt proof
        return pvf.checkReceiptProof(rawReceipt, blockHeader.receiptsRoot, proofSiblings, proofPath);
        return true;
    }

    function _verificationFunction(IProofVerificationFunction verificationFunction) internal view returns (IProofVerificationFunction) {
        if (verificationFunction == DEFAULT_VERIFICATION_FUNCTION) {
            return _defaultVerificationFunction;
        } else {
            return verificationFunction;
        }
    }
}