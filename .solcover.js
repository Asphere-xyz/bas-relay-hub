module.exports = {
  providerOptions: {
    _chainId: 1337,
    _chainIdRpc: 1337,
  },
  client: require("ganache-cli"), // Will load the outermost ganache-cli in node_modules
  skipFiles: [
  ]
};
