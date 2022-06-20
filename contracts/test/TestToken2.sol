// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.6;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

import "../interfaces/IERC20.sol";

contract TestToken2 is Context, IERC20, IERC20Mintable, IERC20Pegged {

    // pre-defined state
    bytes32 private _symbol; // 0
    bytes32 private _name; // 1
    address private _owner; // 2

    // internal state
    mapping(address => uint256) private _balances;
    mapping(address => mapping(address => uint256)) private _allowances;

    uint256 private _totalSupply;
    uint256 private _originChain;
    address private _originAddress;

    function name() public view returns (string memory) {
        return bytes32ToString(_name);
    }

    function symbol() public view returns (string memory) {
        return bytes32ToString(_symbol);
    }

    function bytes32ToString(bytes32 _bytes32) internal pure returns (string memory) {
        uint8 i = 0;
        while (i < 32 && _bytes32[i] != 0) {
            i++;
        }
        bytes memory bytesArray = new bytes(i);
        for (i = 0; i < 32 && _bytes32[i] != 0; i++) {
            bytesArray[i] = _bytes32[i];
        }
        return string(bytesArray);
    }

    function decimals() public pure returns (uint8) {
        return 18;
    }

    function totalSupply() public view virtual override returns (uint256) {
        return _totalSupply;
    }

    function balanceOf(address account) public view virtual override returns (uint256) {
        return _balances[account];
    }

    function transfer(address recipient, uint256 amount) public virtual override returns (bool) {
        _transfer(_msgSender(), recipient, amount);
        return true;
    }

    function allowance(address owner, address spender) public view virtual override returns (uint256) {
        return _allowances[owner][spender];
    }

    function approve(address spender, uint256 amount) public virtual override returns (bool) {
        _approve(_msgSender(), spender, amount);
        return true;
    }

    function transferFrom(address sender, address recipient, uint256 amount) public virtual override returns (bool) {
        _transfer(sender, recipient, amount);
        _approve(sender, _msgSender(), _allowances[sender][_msgSender()] - amount);
        return true;
    }

    function _approve(address owner, address spender, uint256 amount) internal {
        require(owner != address(0));
        require(spender != address(0));
        _allowances[owner][spender] = amount;
        emit Approval(owner, spender, amount);
    }

    function _transfer(address sender, address recipient, uint256 amount) internal {
        require(sender != address(0));
        require(recipient != address(0));
        _balances[sender] -= amount;
        _balances[recipient] += amount;
        emit Transfer(sender, recipient, amount);
    }

    function mint(address account, uint256 amount) public override {
        require(account != address(0));
        _totalSupply += amount;
        _balances[account] += amount;
        emit Transfer(address(0), account, amount);
    }

    function burn(address account, uint256 amount) public onlyOwner override {
        require(account != address(0));
        _balances[account] -= amount;
        _totalSupply -= amount;
        emit Transfer(account, address(0), amount);
    }

    modifier emptyOwner() {
        require(_owner == address(0x00));
        _;
    }

    function initAndObtainOwnership(bytes32 symbol, bytes32 name, uint256 originChain, address originAddress) public emptyOwner {
        _owner = msg.sender;
        _symbol = symbol;
        _name = name;
        _originChain = originChain;
        _originAddress = originAddress;
    }

    modifier onlyOwner() {
        require(msg.sender == _owner);
        _;
    }

    function getOrigin() public view override returns (uint256, address) {
        return (_originChain, _originAddress);
    }
}

contract SimpleTokenFactory_Test2 {
    address private _template;
    constructor() {
        _template = TestTokenFactoryUtils2.deployTestTokenTemplate(this);
    }

    function getImplementation() public view returns (address) {
        return _template;
    }
}

library TestTokenFactoryUtils2 {

    bytes32 constant internal TEST_TOKEN_TEMPLATE_SALT = keccak256("TestTokenTemplateV0.1");

    bytes constant internal TEST_TOKEN_TEMPLATE_BYTECODE = hex"608060405234801561001057600080fd5b50610905806100206000396000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c806370a082311161008c5780639dc29fac116100665780639dc29fac146101a2578063a9059cbb146101b5578063dd62ed3e146101c8578063df1f29ee1461020157600080fd5b806370a082311461015e57806394bfed881461018757806395d89b411461019a57600080fd5b806306fdde03146100d4578063095ea7b3146100f257806318160ddd1461011557806323b872dd14610127578063313ce5671461013a57806340c10f1914610149575b600080fd5b6100dc610224565b6040516100e991906107e9565b60405180910390f35b610105610100366004610780565b610236565b60405190151581526020016100e9565b6005545b6040519081526020016100e9565b610105610135366004610744565b61024c565b604051601281526020016100e9565b61015c610157366004610780565b61029e565b005b61011961016c3660046106f6565b6001600160a01b031660009081526003602052604090205490565b61015c6101953660046107aa565b61033b565b6100dc61038c565b61015c6101b0366004610780565b610399565b6101056101c3366004610780565b610447565b6101196101d6366004610711565b6001600160a01b03918216600090815260046020908152604080832093909416825291909152205490565b600654600754604080519283526001600160a01b039091166020830152016100e9565b6060610231600154610454565b905090565b600061024333848461058b565b50600192915050565b6000610259848484610613565b6001600160a01b03841660009081526004602090815260408083203380855292529091205461029491869161028f908690610856565b61058b565b5060019392505050565b6001600160a01b0382166102b157600080fd5b80600560008282546102c3919061083e565b90915550506001600160a01b038216600090815260036020526040812080548392906102f090849061083e565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef906020015b60405180910390a35050565b6002546001600160a01b03161561035157600080fd5b60028054336001600160a01b031991821617909155600094909455600192909255600655600780549092166001600160a01b03909116179055565b6060610231600054610454565b6002546001600160a01b031633146103b057600080fd5b6001600160a01b0382166103c357600080fd5b6001600160a01b038216600090815260036020526040812080548392906103eb908490610856565b9250508190555080600560008282546104049190610856565b90915550506040518181526000906001600160a01b038416907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200161032f565b6000610243338484610613565b606060005b60208160ff1610801561048d5750828160ff166020811061047c5761047c6108a3565b1a60f81b6001600160f81b03191615155b156104a4578061049c8161086d565b915050610459565b60008160ff1667ffffffffffffffff8111156104c2576104c26108b9565b6040519080825280601f01601f1916602001820160405280156104ec576020820181803683370190505b509050600091505b60208260ff161080156105285750838260ff1660208110610517576105176108a3565b1a60f81b6001600160f81b03191615155b1561058457838260ff1660208110610542576105426108a3565b1a60f81b818360ff168151811061055b5761055b6108a3565b60200101906001600160f81b031916908160001a9053508161057c8161086d565b9250506104f4565b9392505050565b6001600160a01b03831661059e57600080fd5b6001600160a01b0382166105b157600080fd5b6001600160a01b0383811660008181526004602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591015b60405180910390a3505050565b6001600160a01b03831661062657600080fd5b6001600160a01b03821661063957600080fd5b6001600160a01b03831660009081526003602052604081208054839290610661908490610856565b90915550506001600160a01b0382166000908152600360205260408120805483929061068e90849061083e565b92505081905550816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405161060691815260200190565b80356001600160a01b03811681146106f157600080fd5b919050565b60006020828403121561070857600080fd5b610584826106da565b6000806040838503121561072457600080fd5b61072d836106da565b915061073b602084016106da565b90509250929050565b60008060006060848603121561075957600080fd5b610762846106da565b9250610770602085016106da565b9150604084013590509250925092565b6000806040838503121561079357600080fd5b61079c836106da565b946020939093013593505050565b600080600080608085870312156107c057600080fd5b8435935060208501359250604085013591506107de606086016106da565b905092959194509250565b600060208083528351808285015260005b81811015610816578581018301518582016040015282016107fa565b81811115610828576000604083870101525b50601f01601f1916929092016040019392505050565b600082198211156108515761085161088d565b500190565b6000828210156108685761086861088d565b500390565b600060ff821660ff8114156108845761088461088d565b60010192915050565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fdfea2646970667358221220a294ea60ab6992335d5795ef44997a0c2d06e0bc4028c7a7f73ac2253ab47d0f64736f6c63430008060033";

    bytes32 constant internal TEST_TOKEN_TEMPLATE_HASH = keccak256(TEST_TOKEN_TEMPLATE_BYTECODE);

    function deployTestTokenTemplate(SimpleTokenFactory_Test2 templateFactory) internal returns (address) {
        /* we can use any deterministic salt here, since we don't care about it */
        bytes32 salt = TEST_TOKEN_TEMPLATE_SALT;
        /* concat bytecode with constructor */
        bytes memory bytecode = TEST_TOKEN_TEMPLATE_BYTECODE;
        /* deploy contract and store result in result variable */
        address result;
        assembly {
            result := create2(0, add(bytecode, 0x20), mload(bytecode), salt)
        }
        require(result != address(0x00), "deploy failed");
        /* check that generated contract address is correct */
        require(result == simpleTokenTemplateAddress(templateFactory), "address mismatched");
        return result;
    }

    function simpleTokenTemplateAddress(SimpleTokenFactory_Test2 templateFactory) internal pure returns (address) {
        bytes32 hash = keccak256(abi.encodePacked(uint8(0xff), address(templateFactory), TEST_TOKEN_TEMPLATE_SALT, TEST_TOKEN_TEMPLATE_HASH));
        return address(bytes20(hash << 96));
    }
}
