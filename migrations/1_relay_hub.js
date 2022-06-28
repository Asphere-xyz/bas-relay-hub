const ParliaBlockVerifier = artifacts.require("ParliaBlockVerifier");
const RelayHub = artifacts.require("RelayHub");
const SimpleTokenFactory = artifacts.require("SimpleTokenFactory");
const BridgeRouter = artifacts.require("BridgeRouter");
const CrossChainBridge = artifacts.require("CrossChainBridge");

const UNIT_TEST_CONFIG = {
  epochLength: 200,
  nativeTokenSymbol: 'TEST',
  nativeTokenName: 'TEST',
};

const SMART_CHAIN_CONFIG = {
  epochLength: 200,
  nativeTokenSymbol: 'BNB',
  nativeTokenName: 'BNB',
};

const BAS_DEVNET_CONFIG = {
  epochLength: 1200,
  nativeTokenSymbol: 'BAS',
  nativeTokenName: 'BAS',
};

const BLOCKCHAIN_CONFIG = {
  // testnet configs
  'test': UNIT_TEST_CONFIG,
  'soliditycoverage': UNIT_TEST_CONFIG,
  'ganache': UNIT_TEST_CONFIG,
  // BSC config
  'chapel': SMART_CHAIN_CONFIG,
  'bsc': SMART_CHAIN_CONFIG,
  // BAS devnets
  'bas-devnet-1': BAS_DEVNET_CONFIG,
  'bas-devnet-2': BAS_DEVNET_CONFIG,
};

const deployOnlyOnce = async (deployer, contractType, ...constructor) => {
  if (contractType.isDeployed()) {
    return contractType.deployed();
  }
  await deployer.deploy(contractType, ...constructor)
  return contractType.deployed();
}

module.exports = async (deployer) => {
  if (['test', 'soliditycoverage', 'ganache'].includes(deployer.network)) {
    return
  }

  const config = BLOCKCHAIN_CONFIG[deployer.network];
  console.log(`Deploying relay hub to the network: ${deployer.network}`);
  if (!config) throw new Error(`There is no config for network: ${deployer.network}`)
  const {epochLength, nativeTokenSymbol, nativeTokenName} = config;
  console.log(`Network ${deployer.network} config: ${JSON.stringify(config, null, 2)}`);
  // deploy parlia block verifier as default verifier, relay hub and other contracts
  const parliaBlockVerifier = await deployOnlyOnce(deployer, ParliaBlockVerifier, epochLength)
  const epochInterval = await parliaBlockVerifier.getEpochInterval();
  if (`${epochInterval}` !== `${epochLength}`) {
    throw new Error(`Detected bad epoch length: ${epochLength} != ${epochInterval}`)
  }
  const relayHub = await deployOnlyOnce(deployer, RelayHub, parliaBlockVerifier.address)
  const simpleTokenFactory = await deployOnlyOnce(deployer, SimpleTokenFactory)
  const bridgeRouter = await deployOnlyOnce(deployer, BridgeRouter)
  const crossChainBridge = await deployOnlyOnce(deployer, CrossChainBridge)
  // initialize cross bridge
  try {
    await crossChainBridge.initialize(relayHub.address, relayHub.address, simpleTokenFactory.address, bridgeRouter.address, nativeTokenSymbol, nativeTokenName);
  } catch (e) {
    const [errorMessage] = e.message.split('\n')
    if (!errorMessage.includes('Initializable: contract is already initialized')) {
      throw e;
    }
  }
};
