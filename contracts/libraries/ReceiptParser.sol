// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.6;

import "@openzeppelin/contracts/utils/Strings.sol";

import "../interfaces/ICrossChainBridge.sol";

import "./RLP.sol";

library ReceiptParser {

    bytes32 constant TOPIC_PEG_IN_LOCKED = keccak256("DepositLocked(uint256,address,address,address,address,uint256,(bytes32,bytes32,uint256,address))");
    bytes32 constant TOPIC_PEG_IN_BURNED = keccak256("DepositBurned(uint256,address,address,address,address,uint256,(bytes32,bytes32,uint256,address),address)");

    enum PegInType {
        None,
        Lock,
        Burn
    }

    struct State {
        bytes32 receiptHash;
        address contractAddress;
        uint256 chainId;
        address fromAddress;
        address payable toAddress;
        address fromToken;
        address toToken;
        uint256 totalAmount;
        // metadata fields (we can't use Metadata struct here because of Solidity struct memory layout)
        bytes32 symbol;
        bytes32 name;
        uint256 originChain;
        address originAddress;
        address originToken;
    }

    function getMetadata(State memory state) internal pure returns (ICrossChainBridge.Metadata memory) {
        ICrossChainBridge.Metadata memory metadata;
        assembly {
            metadata := add(state, 0x100)
        }
        return metadata;
    }

    function parseTransactionReceipt(bytes calldata rawReceipt) internal view returns (State memory state, PegInType pegInType) {
        uint256 receiptOffset = RLP.openRlp(rawReceipt);
        // parse peg-in data from logs
        uint256 iter = RLP.beginIteration(receiptOffset);
        {
            // postStateOrStatus - we must ensure that tx is not reverted
            require(RLP.toUint256(iter, 1) == 0x01, "tx is reverted");
            iter = RLP.next(iter);
        }
        // skip cumulativeGasUsed
        iter = RLP.next(iter);
        iter = RLP.next(iter);
        // logs - we need to find our logs
        uint256 logs = iter;
        iter = RLP.next(iter);
        uint256 logsIter = RLP.beginIteration(logs);
        for (; logsIter < iter;) {
            uint256 log = logsIter;
            logsIter = RLP.next(logsIter);
            // make sure there is only one peg-in event in logs
            PegInType logType = _decodeReceiptLogs(state, log);
            if (logType != PegInType.None) {
                require(pegInType == PegInType.None, "multiple logs");
                pegInType = logType;
            }
        }
        // don't allow to process if peg-in type is unknown
        require(pegInType != PegInType.None, "missing logs");
        state.receiptHash = keccak256(rawReceipt);
        return (state, pegInType);
    }

    function _decodeReceiptLogs(
        State memory state,
        uint256 log
    ) internal view returns (PegInType pegInType) {
        uint256 logIter = RLP.beginIteration(log);
        address contractAddress;
        {
            // parse smart contract address
            uint256 addressOffset = logIter;
            logIter = RLP.next(logIter);
            contractAddress = RLP.toAddress(addressOffset);
        }
        // topics
        bytes32 mainTopic;
        address fromAddress;
        address toAddress;
        {
            uint256 topicsIter = logIter;
            logIter = RLP.next(logIter);
            // Must be 3 topics RLP encoded: event signature, fromAddress, toAddress
            // Each topic RLP encoded is 33 bytes (0xa0[32 bytes data])
            // Total payload: 99 bytes. Since it's list with total size bigger than 55 bytes we need 2 bytes prefix (0xf863)
            // So total size of RLP encoded topics array must be 101
            if (RLP.itemLength(topicsIter) != 101) {
                return PegInType.None;
            }
            topicsIter = RLP.beginIteration(topicsIter);
            mainTopic = bytes32(RLP.toUintStrict(topicsIter));
            topicsIter = RLP.next(topicsIter);
            fromAddress = address(bytes20(uint160(RLP.toUintStrict(topicsIter))));
            topicsIter = RLP.next(topicsIter);
            toAddress = address(bytes20(uint160(RLP.toUintStrict(topicsIter))));
            topicsIter = RLP.next(topicsIter);
            require(topicsIter == logIter);
            // safety check that iteration is finished
        }

        uint256 ptr = RLP.rawDataPtr(logIter);
        logIter = RLP.next(logIter);
        uint256 len = logIter - ptr;
        {
            // parse logs based on topic type and check that event data has correct length
            uint256 expectedLen;
            if (mainTopic == TOPIC_PEG_IN_LOCKED) {
                expectedLen = 0x100;
                pegInType = PegInType.Lock;
            } else if (mainTopic == TOPIC_PEG_IN_BURNED) {
                expectedLen = 0x120;
                pegInType = PegInType.Burn;
            } else {
                return PegInType.None;
            }
            if (len != expectedLen) {
                return PegInType.None;
            }
        }
        {
            // read chain id separately and verify that contract that emitted event is relevant
            uint256 chainId;
            assembly {
                chainId := calldataload(ptr)
            }
            if (chainId != block.chainid) {
                return PegInType.None;
            }
            // All checks are passed after this point, no errors allowed and we can modify state
            state.chainId = chainId;
            ptr += 0x20;
            len -= 0x20;
        }

        {
            uint256 structOffset;
            assembly {
            // skip 5 fields: receiptHash, contractAddress, chainId, fromAddress, toAddress
                structOffset := add(state, 0xa0)
                calldatacopy(structOffset, ptr, len)
            }
        }
        state.contractAddress = contractAddress;
        state.fromAddress = fromAddress;
        state.toAddress = payable(toAddress);
        return pegInType;
    }
}
