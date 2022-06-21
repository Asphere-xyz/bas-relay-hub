// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.6;

import "../interfaces/ICrossChainBridge.sol";

contract SimpleTokenProxy {

    bytes32 private constant BEACON_SLOT = bytes32(uint256(keccak256('eip1967.proxy.beacon')) - 1);

    fallback() external {
        address bridge;
        bytes32 slot = BEACON_SLOT;
        assembly {
            bridge := sload(slot)
        }
        address impl = ICrossChainBridge(bridge).getTokenImplementation();
        assembly {
            calldatacopy(0, 0, calldatasize())
            let result := delegatecall(gas(), impl, 0, calldatasize(), 0, 0)
            returndatacopy(0, 0, returndatasize())
            switch result
            case 0 {revert(0, returndatasize())}
            default {return (0, returndatasize())}
        }
    }

    function setBeacon(address newBeacon) external {
        address beacon;
        bytes32 slot = BEACON_SLOT;
        assembly {
            beacon := sload(slot)
        }
        require(beacon == address(0x00));
        assembly {
            sstore(slot, newBeacon)
        }
    }
}

library SimpleTokenProxyUtils {

    bytes constant internal SIMPLE_TOKEN_PROXY_BYTECODE = hex"608060405234801561001057600080fd5b50610201806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063d42afb56146100f0575b60008061005960017fa3f0ad74e5423aebfd80d3ef4346578335a9a72aeaee59ff6cb3582b35133d5161014d565b60001b9050805491506000826001600160a01b031663709bc7f36040518163ffffffff1660e01b81526004016020604051808303816000875af11580156100a4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906100c8919061018a565b90503660008037600080366000845af43d6000803e8080156100e9573d6000f35b3d6000fd5b005b6100ee6100fe3660046101ae565b60008061012c60017fa3f0ad74e5423aebfd80d3ef4346578335a9a72aeaee59ff6cb3582b35133d5161014d565b8054925090506001600160a01b0382161561014657600080fd5b9190915550565b60008282101561016d57634e487b7160e01b600052601160045260246000fd5b500390565b6001600160a01b038116811461018757600080fd5b50565b60006020828403121561019c57600080fd5b81516101a781610172565b9392505050565b6000602082840312156101c057600080fd5b81356101a78161017256fea264697066735822122027eb4b9b1405c9e00d3c0cc70d7d3039075720c3fd2738f82dc17b7771cb7ef464736f6c634300080e0033";

    bytes32 constant internal SIMPLE_TOKEN_PROXY_HASH = keccak256(SIMPLE_TOKEN_PROXY_BYTECODE);

    bytes4 constant internal SET_META_DATA_SIG = bytes4(keccak256("initialize(bytes32,bytes32,uint256,address)"));
    bytes4 constant internal SET_BEACON_SIG = bytes4(keccak256("setBeacon(address)"));

    function deploySimpleTokenProxy(address bridge, bytes32 salt, ICrossChainBridge.Metadata memory metaData) internal returns (address) {
        // lets concat bytecode with constructor parameters
        bytes memory bytecode = SIMPLE_TOKEN_PROXY_BYTECODE;
        // deploy new contract and store contract address in result variable
        address result;
        assembly {
            result := create2(0, add(bytecode, 0x20), mload(bytecode), salt)
        }
        require(result != address(0x00), "deploy failed");
        // setup impl
        (bool success,) = result.call(abi.encodePacked(SET_BEACON_SIG, abi.encode(bridge)));
        require(success, "setBeacon failed");
        // setup meta data
        (success,) = result.call(abi.encodePacked(SET_META_DATA_SIG, abi.encode(metaData)));
        require(success, "set metadata failed");
        // return generated contract address
        return result;
    }

    function simpleTokenProxyAddress(address deployer, bytes32 salt) internal pure returns (address) {
        bytes32 bytecodeHash = keccak256(SIMPLE_TOKEN_PROXY_BYTECODE);
        bytes32 hash = keccak256(abi.encodePacked(uint8(0xff), address(deployer), salt, bytecodeHash));
        return address(bytes20(hash << 96));
    }
}
