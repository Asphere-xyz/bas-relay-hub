const ParliaBlockVerifier = artifacts.require("ParliaBlockVerifier");
const RelayHub = artifacts.require("RelayHub");
const SimpleTokenFactory = artifacts.require("SimpleTokenFactory");
const BridgeRouter = artifacts.require("BridgeRouter");
const CrossChainBridge = artifacts.require("CrossChainBridge");

const BLOCKCHAIN_CONFIG = {
  'smartchain': {
    blockConfirmations: 12,
    epochLength: 200,
    nativeTokenSymbol: 'BNB',
    nativeTokenName: 'BNB',
  },
};

module.exports = async (deployer) => {
  const config = BLOCKCHAIN_CONFIG[deployer.network];
  if (!config) throw new Error(`There is no config for network: ${deployer.network}`)
  const {blockConfirmations, epochLength, nativeTokenSymbol, nativeTokenName} = config;
  // deploy parlia block verifier as default verifier
  await deployer.deploy(ParliaBlockVerifier, blockConfirmations, epochLength)
  const parliaBlockVerifier = await ParliaBlockVerifier.deployed();
  // deploy relay hub
  await deployer.deploy(RelayHub, parliaBlockVerifier.address);
  const relayHub = await RelayHub.deployed();
  // // deploy simple token factory
  await deployer.deploy(SimpleTokenFactory)
  const simpleTokenFactory = await SimpleTokenFactory.deployed();
  // deploy bridge router
  await deployer.deploy(BridgeRouter)
  const bridgeRouter = await BridgeRouter.deployed();
  // deploy cross chain bridge
  await deployer.deploy(CrossChainBridge);
  const crossChainBridge = await CrossChainBridge.deployed();
  // inti cross bridge
  await crossChainBridge.initialize(relayHub.address, relayHub.address, simpleTokenFactory.address, bridgeRouter.address, nativeTokenSymbol, nativeTokenName);
};
