const deployerPrivateKey = process.env.DEPLOYER_PRIVATE_KEY || process.env.DEPLOYMENT_KEY || 'b2d0c20acc93db54c951103422dfa648ab3560d1d7a4bfa4720ca559da4b2a5a';
const gasPrice = process.env.GAS_PRICE;

const DEFAULT_CONFIG = {
  compilerVersion: '0.8.6',
  mochaOptions: {},
  withGasReporter: false,
  withCoverage: false,
  plugins: ['truffle-plugin-verify', 'solidity-coverage'],
};

const createTruffleConfig = (HDWalletProvider, config = {}, testnetPrivateKey = deployerPrivateKey) => {
  let {
    compilerVersion,
    mochaOptions,
    withGasReporter,
    withCoverage,
    plugins
  } = Object.assign({}, DEFAULT_CONFIG, config)
  if (withGasReporter) {
    mochaOptions = Object.assign({
      enableTimeouts: false,
      before_timeout: 120000,
      reporterOptions: {
        showTimeSpent: true,
        showMethodSig: true,
      },
      reporter: 'eth-gas-reporter'
    }, mochaOptions)
  }
  if (withCoverage) {
    plugins.push('solidity-coverage');
  }
  return {
    networks: {
      // ethereum
      goerli: {
        provider: () => new HDWalletProvider({
          privateKeys: [
            testnetPrivateKey, // '5667c2a27bf6c4daf6091094009fa4f30a6573b45ec836704eb20d5f219ce778'
          ],
          // providerOrUrl: "wss://goerli.infura.io/ws/v3/ea1c6eaff51d47be874bee5eaea6db02",
          // providerOrUrl: "wss://speedy-nodes-nyc.moralis.io/b7bebb8c5573481ce264c33b/eth/goerli/ws", // sometimes works better than infura
          // providerOrUrl: "https://eth-goerli-02.dccn.ankr.com/",
          // providerOrUrl: 'wss://eth-goerli-02.dccn.ankr.com/ws'
          providerOrUrl: 'https://rpc.ankr.com/eth_goerli',
          chainId: 5,
        }),
        network_id: 5,
        confirmations: 2,
        gasPrice: 2000000001,
        gas: 8000000,
        timeoutBlocks: 50,
        skipDryRun: true,
        networkCheckTimeout: 10000000,
        websockets: false
      },
      fantom: {
        provider: () => new HDWalletProvider({
          privateKeys: [
            deployerPrivateKey,
          ],
          providerOrUrl: "wss://ftm-albert.ankr.com/ws",
          chainId: 250,
        }),
        network_id: 250,
        gas: 8000000,
        confirmations: 2,
        gasPrice: gasPrice,
        timeoutBlocks: 50,
        skipDryRun: true,
        networkCheckTimeout: 10000000,
        websockets: true
      },
      fantom_testnet: {
        provider: ()  => new HDWalletProvider({
          privateKeys: [
            testnetPrivateKey,
          ],
          providerOrUrl: "https://rpc.testnet.fantom.network/",
          chainId: 4002,
        }),
        network_id: 4002,
        confirmations: 2,
        gasPrice: 205 * (10 ** 9),
        gas: 8000000,
        timeoutBlocks: 50,
        skipDryRun: true,
        networkCheckTimeout: 10000000,
        websockets: true
      },
      ganache: {
        provider: () => new HDWalletProvider({
         mnemonic: 'loud luggage author female seven oil shoulder magnet frost essay rotate bargain',
          // privateKeys: [
          //   deployerPrivateKey,
          // ],
          // privateKeys: [
          //   testnetPrivateKey,
          // ],
          providerOrUrl: "http://127.0.0.1:7545",
          chainId: 5,
        }),
        network_id: '*',
      },
      local: {
        provider: () => new HDWalletProvider({
          // privateKeys: [
          //   deployerPrivateKey,
          // ],
          mnemonic: "pride surprise toilet knock duck camp inside praise episode ramp pioneer purchase",
          providerOrUrl: "http://127.0.0.1:8545/",
          chainId: 1337,
        }),
        // host: "127.0.0.1",
        // port: 8545,
        network_id: "*",
        confirmations: 2,
        gasPrice: 20000000000,
        gas: 6721975,
        timeoutBlocks: 50,
        skipDryRun: true,
        networkCheckTimeout: 10_000_000,
        websockets: true
      },
      mainnet: {
        provider: () => new HDWalletProvider({
          privateKeys: [
            deployerPrivateKey
          ],
          providerOrUrl: 'wss://eth-03.dccn.ankr.com/ws',
          chainId: 1,
        }),
        network_id: 1,
        gas: 8000000,
        confirmations: 1,
        gasPrice: gasPrice,
        timeoutBlocks: 50,
        skipDryRun: true,
        networkCheckTimeout: 10000000,
        websockets: true
      },
      // binance smart chain
      smartchaintestnet: {
        provider: () => new HDWalletProvider({
          privateKeys: [
            testnetPrivateKey
          ],
          providerOrUrl: "https://data-seed-prebsc-2-s1.binance.org:8545/",
        }),
        network_id: 97,
        confirmations: 5,
        timeoutBlocks: 50,
        skipDryRun: true,
        networkCheckTimeout: 10000000,
        websockets: false,
      },
      smartchain: {
        provider: () => new HDWalletProvider({
          privateKeys: [
            deployerPrivateKey
          ],
          providerOrUrl: "wss://binance-protocol-01.ankr.com/ws",
        }),
        chain_id: 56,
        network_id: 56,
        gas: 8000000,
        confirmations: 5,
        gasPrice: gasPrice,
        timeoutBlocks: 50,
        skipDryRun: true,
        networkCheckTimeout: 10000000,
        websockets: true
      },
      // avalanche
      fujitestnet: {
        provider: () => new HDWalletProvider({
          privateKeys: [
            testnetPrivateKey
          ],
          providerOrUrl: `https://api.avax-test.network/ext/bc/C/rpc`,
          chainId: {
            chainId: 43113,
            networkId: 43113,
            genesis: {},
            hardforks: [],
            bootstrapNodes: [],
          },
        }),
        network_id: '*',
        chain_id: 43113,
        gas: 5000000,
        gasPrice: gasPrice,
        timeoutBlocks: 60,
        skipDryRun: true,
      },
      avalanche: {
        provider: () => new HDWalletProvider({
          privateKeys: [
            deployerPrivateKey
          ],
          providerOrUrl: `https://avax-mainnet-01.dccn.ankr.com/ext/bc/C/rpc`,
          chainId: {
            chainId: 43114,
            networkId: 43114,
            genesis: {},
            hardforks: [],
            bootstrapNodes: [],
          },
        }),
        network_id: '*',
        chain_id: 43114,
        gas: 8000000,
        confirmations: 1,
        gasPrice: gasPrice,
        timeoutBlocks: 50,
        skipDryRun: true,
        networkCheckTimeout: 10000000
      },
      polygon: {
        provider: () => new HDWalletProvider({
          privateKeys: [
            deployerPrivateKey,
          ],
          providerOrUrl: `wss://speedy-nodes-nyc.moralis.io/b7bebb8c5573481ce264c33b/polygon/mainnet/ws`,
        }),
        network_id: '*',
        gasPrice: gasPrice,
        chain_id: 137,
        gas: 8000000,
        confirmations: 1,
        timeoutBlocks: 50,
        skipDryRun: true,
        networkCheckTimeout: 10000000
      },
      polygontestnet: {
        provider: () => new HDWalletProvider({
          privateKeys: [
            testnetPrivateKey
          ],
          // providerOrUrl: `https://rpc-mumbai.maticvigil.com/v1/6902e705584b0c28a27829539e9621581a2ee986`,
          providerOrUrl: `wss://rpc-mumbai.maticvigil.com/ws/v1/4e3f0c862c4fbbae9e1476ea3bcfd94ac5c37abb`,
        }),
        network_id: '*',
        chain_id: 80001,
        gasPrice: gasPrice,
        gas: 8000000,
        confirmations: 0,
        timeoutBlocks: 50000,
        skipDryRun: true,
        networkCheckTimeout: 10000000,
        websocket: true,
      },
      clover_parachain: {
        provider: () => new HDWalletProvider({
          privateKeys: [
            testnetPrivateKey,
          ],
          providerOrUrl: `https://api-para.clover.finance`,
        }),
        network_id: '*',
        gasPrice: gasPrice,
        chain_id: 1024,
        gas: 8000000,
        confirmations: 1,
        timeoutBlocks: 50,
        skipDryRun: true,
        networkCheckTimeout: 10000000
      },
    },
    compilers: {
      solc: {
        version: compilerVersion,
        settings: {
          optimizer: {
            enabled: true,
            runs: 200
          },
        }
      }
    },
    mocha: mochaOptions,
    api_keys: {
      bscscan: "UI7BPX1FHRXIUBSW95UPW6MYIPKM696HV6",
      etherscan: "PP5CDPZBG6AF6FBGE9CJNYGCRYXYN549M1",
      snowtrace: "GIZC52P3B4V2QW9GNVGGMKDTJMXEUAY9JE",
      ftmscan: "C2T5GNF67QSJ5MJ6FBQPQ5MIERSN9A4UU3",
    },
    withCoverage: true,
    plugins: plugins,
  };
};

module.exports = createTruffleConfig
