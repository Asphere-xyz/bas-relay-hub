const HDWalletProvider = require('@truffle/hdwallet-provider');

module.exports = {
  compilers: {
    solc: {
      version: "0.8.11",
      settings: {
        optimizer: {
          enabled: true,
          runs: 200
        },
      }
    }
  },
  networks: {
    ganache: {
      provider: () => new HDWalletProvider({
        mnemonic: 'inherit daring trust actual engine pair swap cargo subject lawsuit length hurt',
        providerOrUrl: "http://127.0.0.1:7545",
        chainId: 5777,
      }),
      network_id: 5777,
      confirmations: 1,
      // gasPrice: 20_000000000,
      // gas: 30_000_000,
      timeoutBlocks: 50,
      skipDryRun: true,
      networkCheckTimeout: 10_000_000,
      websockets: true
    },

  },
  mocha: {
    enableTimeouts: false,
    // reporterOptions: {
    //   showTimeSpent: true,
    //   showMethodSig: true,
    // },
    // reporter: 'eth-gas-reporter'
  },
  plugins: [
    "solidity-coverage"
  ]
};
