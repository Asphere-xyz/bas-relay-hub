module.exports = {
  providerOptions: {
    _chainId: 1337,
    _chainIdRpc: 1337,
  },
  client: require("ganache-cli"), // Will load the outermost ganache-cli in node_modules
  skipFiles: [
    'test/TestBond.sol',
    'test/TestToken.sol',
    'test/TestToken2.sol',
    'test/TestNonRebasingBond.sol',
    'interfaces/ICrossChainBridge.sol',
    'interfaces/IERC20.sol',
    'interfaces/IInternetBondRatioFeed.sol',
  ]
};
