// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.6;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

import "../interfaces/IERC20.sol";

contract TestToken is Context, IERC20, IERC20Mintable, IERC20Pegged {

    // pre-defined state
    bytes32 private _symbol; // 0
    bytes32 private _name; // 1
    address private _owner; // 2

    // internal state
    mapping(address => uint256) private _balances;
    mapping(address => mapping(address => uint256)) private _allowances;

    uint256 private _totalSupply;
    address public _originAddress;
    uint256 private _originChain;

    function name() public view returns (string memory) {
        return "Test Token";
    }

    function symbol() public view returns (string memory) {
        return "TeST";
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
        return 10;
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

    function mint(address account, uint256 amount) public onlyOwner override {
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

contract SimpleTokenFactory_Test {
    address private _template;
    constructor() {
        _template = TestTokenFactoryUtils.deployTestTokenTemplate(this);
    }

    function getImplementation() public view returns (address) {
        return _template;
    }
}

library TestTokenFactoryUtils {

    bytes32 constant internal TEST_TOKEN_TEMPLATE_SALT = keccak256("TestTokenTemplateV0.1");

    bytes constant internal TEST_TOKEN_TEMPLATE_BYTECODE = hex"608060405234801561001057600080fd5b506107f9806100206000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c806394bfed881161008c578063a9059cbb11610066578063a9059cbb14610202578063dd62ed3e14610215578063df1f29ee1461024e578063e99fe3b31461027157600080fd5b806394bfed88146101bc57806395d89b41146101cf5780639dc29fac146101ef57600080fd5b806323b872dd116100c857806323b872dd1461015c578063313ce5671461016f57806340c10f191461017e57806370a082311461019357600080fd5b806306fdde03146100ef578063095ea7b31461012757806318160ddd1461014a575b600080fd5b60408051808201909152600a8152692a32b9ba102a37b5b2b760b11b60208201525b60405161011e9190610729565b60405180910390f35b61013a6101353660046106c0565b61029c565b604051901515815260200161011e565b6005545b60405190815260200161011e565b61013a61016a366004610684565b6102b2565b604051600a815260200161011e565b61019161018c3660046106c0565b610304565b005b61014e6101a136600461062f565b6001600160a01b031660009081526003602052604090205490565b6101916101ca3660046106ea565b6103b8565b604080518082019091526004815263151954d560e21b6020820152610111565b6101916101fd3660046106c0565b610409565b61013a6102103660046106c0565b6104b7565b61014e610223366004610651565b6001600160a01b03918216600090815260046020908152604080832093909416825291909152205490565b600754600654604080519283526001600160a01b0390911660208301520161011e565b600654610284906001600160a01b031681565b6040516001600160a01b03909116815260200161011e565b60006102a93384846104c4565b50600192915050565b60006102bf84848461054c565b6001600160a01b0384166000908152600460209081526040808320338085529252909120546102fa9186916102f5908690610796565b6104c4565b5060019392505050565b6002546001600160a01b0316331461031b57600080fd5b6001600160a01b03821661032e57600080fd5b8060056000828254610340919061077e565b90915550506001600160a01b0382166000908152600360205260408120805483929061036d90849061077e565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef906020015b60405180910390a35050565b6002546001600160a01b0316156103ce57600080fd5b60028054336001600160a01b031991821617909155600094909455600192909255600755600680549092166001600160a01b03909116179055565b6002546001600160a01b0316331461042057600080fd5b6001600160a01b03821661043357600080fd5b6001600160a01b0382166000908152600360205260408120805483929061045b908490610796565b9250508190555080600560008282546104749190610796565b90915550506040518181526000906001600160a01b038416907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef906020016103ac565b60006102a933848461054c565b6001600160a01b0383166104d757600080fd5b6001600160a01b0382166104ea57600080fd5b6001600160a01b0383811660008181526004602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591015b60405180910390a3505050565b6001600160a01b03831661055f57600080fd5b6001600160a01b03821661057257600080fd5b6001600160a01b0383166000908152600360205260408120805483929061059a908490610796565b90915550506001600160a01b038216600090815260036020526040812080548392906105c790849061077e565b92505081905550816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405161053f91815260200190565b80356001600160a01b038116811461062a57600080fd5b919050565b60006020828403121561064157600080fd5b61064a82610613565b9392505050565b6000806040838503121561066457600080fd5b61066d83610613565b915061067b60208401610613565b90509250929050565b60008060006060848603121561069957600080fd5b6106a284610613565b92506106b060208501610613565b9150604084013590509250925092565b600080604083850312156106d357600080fd5b6106dc83610613565b946020939093013593505050565b6000806000806080858703121561070057600080fd5b84359350602085013592506040850135915061071e60608601610613565b905092959194509250565b600060208083528351808285015260005b818110156107565785810183015185820160400152820161073a565b81811115610768576000604083870101525b50601f01601f1916929092016040019392505050565b60008219821115610791576107916107ad565b500190565b6000828210156107a8576107a86107ad565b500390565b634e487b7160e01b600052601160045260246000fdfea26469706673582212203970b5e6c2ed2b1da7852d84955c57ef6726ae4a1cb6c6bd39fa7828c73312bd64736f6c63430008060033";

    bytes32 constant internal TEST_TOKEN_TEMPLATE_HASH = keccak256(TEST_TOKEN_TEMPLATE_BYTECODE);

    function deployTestTokenTemplate(SimpleTokenFactory_Test templateFactory) internal returns (address) {
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

    function simpleTokenTemplateAddress(SimpleTokenFactory_Test templateFactory) internal pure returns (address) {
        bytes32 hash = keccak256(abi.encodePacked(uint8(0xff), address(templateFactory), TEST_TOKEN_TEMPLATE_SALT, TEST_TOKEN_TEMPLATE_HASH));
        return address(bytes20(hash << 96));
    }
}
