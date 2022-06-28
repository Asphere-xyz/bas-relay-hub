const HDWalletProvider = require('@truffle/hdwallet-provider');

const DEFAULT_PRIVATE_KEY = 'b2d0c20acc93db54c951103422dfa648ab3560d1d7a4bfa4720ca559da4b2a5a'
const ENABLE_GAS_REPORTER = false;

const deployerPrivateKey = process.env.DEPLOYER_PRIVATE_KEY || DEFAULT_PRIVATE_KEY,
  gasPrice = process.env.GAS_PRICE;

let mochaOptions = {enableTimeouts: false}
if (Boolean(ENABLE_GAS_REPORTER)) {
  Object.assign(mochaOptions, {
    reporterOptions: {
      showTimeSpent: true,
      showMethodSig: true
    },
    reporter: 'eth-gas-reporter'
  })
}

const DEFAULT_OPTS = {
  confirmations: 2,
  timeoutBlocks: 50,
  skipDryRun: true,
  networkCheckTimeout: 10000000,
  websockets: false,
};

module.exports = {
  compilers: {
    solc: {
      version: "0.8.14",
      settings: {
        optimizer: {
          enabled: true,
          runs: 20
        }
      }
    }
  },
  networks: {
    // unit test
    develop: {
      host: "localhost",
      port: 8545,
      network_id: "*"
    },
    ganache: {
      host: "localhost",
      port: 7545,
      network_id: "*",
      gas: 100_000_000
    },
    // BSC
    'chapel': {
      provider: () => new HDWalletProvider({
        privateKeys: [
          deployerPrivateKey,
        ],
        providerOrUrl: "https://data-seed-prebsc-1-s1.binance.org:8545/",
        chainId: 97,
      }),
      network_id: "*",
      ...DEFAULT_OPTS
    },
    // BAS devnets
    'bas-devnet-1': {
      provider: () => new HDWalletProvider({
        privateKeys: [
          deployerPrivateKey,
        ],
        providerOrUrl: "http://rpc.dev-01.bas.ankr.com:8545",
        chainId: 14000,
      }),
      network_id: "*",
      ...DEFAULT_OPTS
    },
    'bas-devnet-2': {
      provider: () => new HDWalletProvider({
        privateKeys: [
          deployerPrivateKey,
        ],
        providerOrUrl: "http://rpc.dev-02.bas.ankr.com:8545",
        chainId: 14001,
      }),
      network_id: "*",
      ...DEFAULT_OPTS
    },
  },
  mocha: mochaOptions,
  plugins: [
    "solidity-coverage"
  ]
};
