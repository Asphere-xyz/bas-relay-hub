// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.0;

import "../libraries/ReceiptParser.sol";

contract ReceiptParserHelper {

    function parseTransactionReceipt(bytes calldata rawReceipt) external view returns (ReceiptParser.State memory state, ReceiptParser.PegInType pegInType) {
        return ReceiptParser.parseTransactionReceipt(rawReceipt);
    }
}